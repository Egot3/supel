package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/Egot3/supel/backend/timetable/internal/database"
	"github.com/Egot3/supel/backend/timetable/internal/database/repositories"
	storage "github.com/Egot3/supel/backend/timetable/internal/s3"
	"github.com/Egot3/supel/backend/timetable/internal/server"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	database.InitDB()
	ctx := context.Background()
	if err := database.RunMigrations(ctx, database.DB); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("timetable", grpc_health_v1.HealthCheckResponse_SERVING)

	s3Config := storage.Config{
		Endpoint: fmt.Sprintf("http://%v:%v", os.Getenv("STORAGE_HOST"),
			os.Getenv("STORAGE_PORT")),
		AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		Bucket:    os.Getenv("STORAGE_TIMETABLE_BUCKET"),
	}

	client, err := storage.NewClient(ctx, s3Config)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	s3ConfigPresigned := storage.Config{
		Endpoint: fmt.Sprintf("http://%v:%v", os.Getenv("STORAGE_PUBLIC_HOST"),
			os.Getenv("STORAGE_PUBLIC_PORT")),
		AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		Bucket:    os.Getenv("STORAGE_TIMETABLE_BUCKET"),
	}

	presignedClientClient, err := storage.NewClient(ctx, s3ConfigPresigned)
	if err != nil {
		log.Fatalf("Failed to create signed storage client: %v", err)
	}
	presignedClient := s3.NewPresignClient(presignedClientClient)

	s3Service := storage.NewStorageService(client, presignedClient, s3Config.Bucket)
	if err := s3Service.EnsureBuckets(ctx, []string{os.Getenv("STORAGE_TIMETABLE_BUCKET")}); err != nil {
		log.Fatal(err)
	}

	abstractRepo := repositories.NewAbstractLessonRepository(database.DB)
	concreteRepo := repositories.NewConcreteLessonRepository(database.DB)
	homeworkAttachmentRepo := repositories.NewHomeworkAttachmentRepository(database.DB)

	TimetableServer := server.NewTimetableService(*s3Service, abstractRepo, concreteRepo, homeworkAttachmentRepo)
	ttpb.RegisterTimetableServiceServer(grpcServer, TimetableServer)

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("timetable Service gRPC server on %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
