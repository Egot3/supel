package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Egot3/supel/backend/contracts"
	"github.com/Egot3/supel/backend/news/internal/database"
	"github.com/Egot3/supel/backend/news/internal/middleware"
	storage "github.com/Egot3/supel/backend/news/internal/s3"
	"github.com/Egot3/supel/backend/news/internal/server"

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

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.AuthInterceptor),
	)

	healthServer := grpc.NewServer()
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(healthServer, healthService)
	healthService.SetServingStatus("news", grpc_health_v1.HealthCheckResponse_SERVING)
	healthService.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	s3Config := storage.Config{
		Endpoint: fmt.Sprintf("http://%v:%v", os.Getenv("STORAGE_HOST"),
			os.Getenv("STORAGE_PORT")),
		AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		Bucket:    os.Getenv("STORAGE_NEWS_BUCKET"),
	}

	client, err := storage.NewClient(ctx, s3Config)
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	s3Service := storage.NewStorageService(client, s3Config.Bucket)

	if err := s3Service.EnsureBuckets(ctx, []string{os.Getenv("STORAGE_NEWS_BUCKET")}); err != nil {
		log.Fatal(err)
	}

	NewsServer := server.NewNewsService(*s3Service)
	pb.RegisterNewsServiceServer(grpcServer, NewsServer)

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Printf("news Service gRPC server on %s", port)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
