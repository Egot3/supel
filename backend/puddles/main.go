package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/egot3/supel/backend/puddles/internal/database"
	"github.com/egot3/supel/backend/puddles/internal/hooks"
	"github.com/egot3/supel/backend/puddles/internal/interceptors"
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

	db, _ := do.Invoke[*bun.DB](injector)

	ctx := context.Background()
	if err := database.RunMigrations(ctx, db); err != nil {
		log.Fatalf("Fatal Migraton Fail(FMF): %s", err)
	}
	db.AddQueryHook(&hooks.SlogQueryHook{})

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.LoggingUnaryInterceptor(logger)),
	)
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("group", grpc_health_v1.HealthCheckResponse_SERVING)

}
