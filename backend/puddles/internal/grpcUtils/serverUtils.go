package grpcutils

import (
	"context"

	"github.com/google/uuid"
)

type ServerUtils interface {
	UserFromContext(ctx context.Context) (userUUID uuid.UUID, ok bool)
}

type GRPCServerUtils struct{}
