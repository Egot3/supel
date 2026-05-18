package interceptors

import (
	"context"
	"log/slog"
	"time"

	"github.com/egot3/supel/backend/rbac/internal/logctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LoggingUnaryInterceptor(baseLogger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		requestUUID := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get("request-uuid"); len(vals) != 0 {
				requestUUID = vals[0]
			}
		}

		logger := baseLogger.With(
			slog.String("method", info.FullMethod),
			slog.String("request_uuid", requestUUID),
		)

		ctx = logctx.WithLogger(ctx, logger)

		logger.InfoContext(ctx, "entered method")

		resp, err = handler(ctx, req)

		code := status.Code(err)
		elapsed := time.Since(start)
		if err != nil {
			logger.ErrorContext(ctx, "done",
				slog.Duration("elapsed", elapsed),
				slog.String("code", code.String()),
				slog.String("err", err.Error()),
			)
		} else {
			logger.InfoContext(ctx, "done",
				slog.Duration("elapsed", elapsed),
				slog.String("code", code.String()),
			)
		}
		return resp, err
	}
}
