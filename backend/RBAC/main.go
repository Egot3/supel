package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/rbac/internal/database"
	"github.com/egot3/supel/backend/rbac/internal/database/repositories"
	"github.com/egot3/supel/backend/rbac/internal/hooks"
	"github.com/egot3/supel/backend/rbac/internal/interceptors"
	"github.com/egot3/supel/backend/rbac/internal/models"
	"github.com/egot3/supel/backend/rbac/internal/server"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	injector := do.New()
	do.Provide(injector, database.InitDB)

	db := do.MustInvoke[*bun.DB](injector)

	ctx := context.Background()
	if err := database.RunMigrations(ctx, db); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}
	db.RegisterModel((*models.RolesActions)(nil))
	db.AddQueryHook(&hooks.SlogQueryHook{})

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.LoggingUnaryInterceptor(logger)),
	)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("RBAC", grpc_health_v1.HealthCheckResponse_SERVING)

	do.Provide(injector, repositories.NewActionRepository)
	do.Provide(injector, repositories.NewRoleRepository)
	do.Provide(injector, server.NewRBACService)

	RBACServer, err := do.Invoke[*server.RBACService](injector)
	if err != nil {
		log.Fatalf("Unable to run server: %v", err)
	}
	rbacpb.RegisterRBACServiceServer(grpcServer, RBACServer)

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
