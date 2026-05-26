package testutils

import (
	"context"

	pb "github.com/Egot3/supel/backend/contracts"
	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/timetable/internal/services"
	"github.com/samber/do/v2"

	"github.com/google/uuid"
)

type StubRBACClient struct {
	Rules map[string]bool
	Err   error
}

type StubClient struct {
	HasPermissionFunc  func(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb rbacpb.Verb) (bool, error)
	CheckExistanceFunc func(ctx context.Context, userUUID uuid.UUID) (bool, error)
	GetUserFunc        func(ctx context.Context, userUUID uuid.UUID) (*pb.User, error)
}

func (s *StubClient) HasPermission(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb rbacpb.Verb) (bool, error) {
	return s.HasPermissionFunc(ctx, userUUID, scope, subScope, verb)
}

func (s *StubClient) CheckExistance(ctx context.Context, userUUID uuid.UUID) (bool, error) {
	return s.CheckExistanceFunc(ctx, userUUID)
}

func (s *StubClient) GetUser(ctx context.Context, userUUID uuid.UUID) (*pb.User, error) {
	return s.GetUserFunc(ctx, userUUID)
}

func AllowAllStub(i do.Injector) (services.Client, error) {
	return &StubClient{
		HasPermissionFunc: func(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb rbacpb.Verb) (bool, error) {
			return true, nil
		},
		CheckExistanceFunc: func(ctx context.Context, userUUID uuid.UUID) (bool, error) {
			return true, nil
		},
		GetUserFunc: func(ctx context.Context, userUUID uuid.UUID) (*pb.User, error) {
			return &pb.User{Uuid: userUUID.String()}, nil
		},
	}, nil
}
