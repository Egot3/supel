package main

import (
	"context"
	"log"
	"net"

	ttpb "github.com/Egot3/supel/backend/contracts/timetable"
	"github.com/egot3/supel/backend/timetable/internal/database"
	homeworkattachment "github.com/egot3/supel/backend/timetable/internal/database/repositories/homeworkAttachment"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/lesson"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/period"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/subject"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/timetable"
	"github.com/egot3/supel/backend/timetable/internal/database/repositories/timetableentry"
	grpcutils "github.com/egot3/supel/backend/timetable/internal/grpcUtils"
	storage "github.com/egot3/supel/backend/timetable/internal/s3"
	"github.com/egot3/supel/backend/timetable/internal/server"
	"github.com/egot3/supel/backend/timetable/internal/services"
	testutils "github.com/egot3/supel/backend/timetable/internal/testUtils"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	injector := do.New()
	do.Provide(injector, database.InitDB)

	db, err := do.Invoke[*bun.DB](injector)

	ctx := context.Background()
	if err := database.RunMigrations(ctx, db); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("timetable", grpc_health_v1.HealthCheckResponse_SERVING)

	do.Provide(injector, period.NewPeriodRepository)
	do.Provide(injector, lesson.NewLessonRepository)
	do.Provide(injector, subject.NewSubjectRepository)
	do.Provide(injector, homeworkattachment.NewHomeworkAttachmentRepository)
	do.Provide(injector, timetable.NewTimetableRepository)
	do.Provide(injector, timetableentry.NewTimetableEntryRepository)
	do.ProvideNamed(injector, "s3config.unsigned", storage.GenerateUnsignedS3Config)
	do.Provide(injector, storage.NewUnsignedClient)

	do.ProvideNamed(injector, "s3config.presigned", storage.GeneratePresignedS3Config)
	do.Provide(injector, storage.NewPresignedClient)

	do.Provide(injector, storage.NewStorageService)

	do.Provide(injector, testutils.AllowAllStub)
	do.Provide(injector, testutils.AllowAllUsers)

	do.Provide(injector, services.NewRBACClient)
	do.Provide(injector, services.NewUserService)
	do.Provide(injector, services.NewIdentityService)

	do.Provide(injector, services.NewGRPCClient)
	do.Provide(injector, func(injector do.Injector) (grpcutils.ServerUtils, error) {
		return grpcutils.GRPCServerUtils{}, nil
	})

	do.Provide(injector, server.NewTimetableService)

	svc, err := do.Invoke[*server.TimetableService](injector)
	if err != nil {
		log.Fatalf("Unable to run server: %v", err)
	}
	ttpb.RegisterTimetableServiceServer(grpcServer, svc)

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
