package grpcutils

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func (GRPCServerUtils) UserFromContext(ctx context.Context) (userUUID uuid.UUID, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, false
	}
	rawUserUUID := md.Get("user-uuid")[0]
	userUUID, err := uuid.Parse(rawUserUUID)
	if err != nil {
		return uuid.Nil, false
	}

	return userUUID, len(userUUID) != 0
}
