package testutils

import (
	"context"

	grpcutils "github.com/egot3/supel/backend/timetable/internal/grpcUtils"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type ServerUtilsMock struct {
	UserFromContextFn func(ctx context.Context) (userUUID uuid.UUID, ok bool)
}

func (m *ServerUtilsMock) UserFromContext(ctx context.Context) (userUUID uuid.UUID, ok bool) {
	return m.UserFromContextFn(ctx)
}

func AllowAllUsers(i do.Injector) (grpcutils.ServerUtils, error) {
	return &ServerUtilsMock{
		UserFromContextFn: func(ctx context.Context) (userUUID uuid.UUID, ok bool) { return uuid.Nil, true },
	}, nil
}
