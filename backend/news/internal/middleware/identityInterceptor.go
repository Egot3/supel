package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUserRole contextKey = "user_role"
)

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if info.FullMethod == "/grpc.health.v1.Health/Check" ||
		info.FullMethod == "/grpc.health.v1.Health/Watch" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	userIDs := md.Get("x-user-id")
	if len(userIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing user identity")
	}

	roles := md.Get("x-user-role")
	role := ""
	if len(roles) > 0 {
		role = roles[0]
	}

	ctx = context.WithValue(ctx, ContextKeyUserID, userIDs[0])
	ctx = context.WithValue(ctx, ContextKeyUserRole, role)

	return handler(ctx, req)
}

func UserFromContext(ctx context.Context) (userID string, role string, ok bool) {
	userID, ok1 := ctx.Value(ContextKeyUserID).(string)
	role, ok2 := ctx.Value(ContextKeyUserRole).(string)
	return userID, role, ok1 && ok2
}
