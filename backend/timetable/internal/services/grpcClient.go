package services

import (
	"context"

	grpb "github.com/Egot3/supel/backend/contracts/group"
	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type Client interface {
	HasPermission(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb rbacpb.Verb) (bool, error)
	UsersGroups(ctx context.Context, userUUID string) ([]string, error)
}

type GRPCClient struct {
	RBACStub  rbacpb.RBACServiceClient
	GroupStub grpb.GroupServiceClient
}

func NewGRPCClient(i do.Injector) (Client, error) {
	RBACStub := do.MustInvoke[rbacpb.RBACServiceClient](i)
	GroupStub := do.MustInvoke[grpb.GroupServiceClient](i)

	return &GRPCClient{
		RBACStub:  RBACStub,
		GroupStub: GroupStub,
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

func (c *GRPCClient) UsersGroups(ctx context.Context, userUUID string) ([]string, error) {
	resp, err := c.GroupStub.MembersGroups(ctx, &grpb.MembersGroupsRequest{
		Page:       0,
		Size:       50,
		MemberUUID: userUUID,
	})
	if err != nil {
		return nil, err
	}

	groupsUUIDs := make([]string, len(resp.Groups))
	for i, group := range resp.Groups {
		groupsUUIDs[i] = group.UUID
	}
	return groupsUUIDs, nil
}
