package services

import (
	"context"

	pb "github.com/Egot3/supel/backend/contracts"
	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type Client interface {
	HasPermission(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb rbacpb.Verb) (bool, error)
	CheckExistance(ctx context.Context, userUUID uuid.UUID) (bool, error)
	GetUser(ctx context.Context, userUUID uuid.UUID) (*pb.User, error)
}

type GRPCClient struct {
	RBACStub     rbacpb.RBACServiceClient
	IdentityStub pb.IdentityServiceClient
	UserStub     pb.UserServiceClient
}

func NewGRPCClient(i do.Injector) (GRPCClient, error) {
	RBACStub := do.MustInvoke[rbacpb.RBACServiceClient](i)
	IdentityStub := do.MustInvoke[pb.IdentityServiceClient](i)
	UserStub := do.MustInvoke[pb.UserServiceClient](i)

	return GRPCClient{
		RBACStub:     RBACStub,
		IdentityStub: IdentityStub,
		UserStub:     UserStub,
	}, nil
}

func (c *GRPCClient) HasPermission(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb rbacpb.Verb) (bool, error) {
	resp, err := c.RBACStub.HasPermission(ctx, &rbacpb.HasPermissionQuestion{
		UserUUID: userUUID.String(),
		Scope:    scope,
		SubScope: subScope,
		Verb:     verb,
	})
	if err != nil {
		return false, err
	}

	return resp.Has, err
}

func (c *GRPCClient) CheckExistance(ctx context.Context, userUUID uuid.UUID) (bool, error) {
	resp, err := c.IdentityStub.CheckExistance(ctx, &pb.CheckExistanceRequest{
		Uuid: userUUID.String(),
	})
	if err != nil {
		return false, err
	}

	return resp.Exists, err
}

func (c *GRPCClient) GetUser(ctx context.Context, userUUID uuid.UUID) (*pb.User, error) {
	resp, err := c.UserStub.GetUser(ctx, &pb.GetUserRequest{
		Uuid: userUUID.String(),
	})

	if err != nil {
		return nil, err
	}

	return resp.User, nil
}
