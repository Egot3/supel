package main

import (
	"context"
	"log"

	"github.com/egot3/supel/backend/rbac/internal/database"
	"github.com/egot3/supel/backend/rbac/internal/models"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	injector := do.New()
	do.Provide(injector, database.InitDB)

	db, _ := do.Invoke[*bun.DB](injector)

	ctx := context.Background()
	if err := database.RunMigrations(ctx, db); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}
	db.RegisterModel((*models.RolesActions)(nil))

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("RBAC", grpc_health_v1.HealthCheckResponse_SERVING)

}
