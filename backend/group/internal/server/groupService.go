package server

import (
	"context"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"

	pb "github.com/Egot3/supel/backend/contracts"
	grpb "github.com/Egot3/supel/backend/contracts/group"

	"github.com/egot3/supel/backend/group/internal/database/repositories/curator"
	"github.com/egot3/supel/backend/group/internal/database/repositories/group"
	"github.com/egot3/supel/backend/group/internal/database/repositories/member"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"google.golang.org/grpc/metadata"
)

type GroupService struct {
	grpb.UnimplementedGroupServiceServer
	GroupRepository   group.GroupRepository
	CuratorRepository curator.CuratorRepository
	MemberRepository  member.MemberRepository
	RBACClient        rbacpb.RBACServiceClient
	IdentityClient    pb.IdentityServiceClient
	UserClient        pb.UserServiceClient
}

func UserFromContext(ctx context.Context) (userUUID uuid.UUID, ok bool) {
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

func NewGroupService(i do.Injector) (*GroupService, error) {
	groupRepository, err := do.Invoke[group.GroupRepository](i)
	if err != nil {
		return nil, err
	}

	curatorRepository, err := do.Invoke[curator.CuratorRepository](i)
	if err != nil {
		return nil, err
	}

	memberRepository, err := do.Invoke[member.MemberRepository](i)
	if err != nil {
		return nil, err
	}

	RBACClient, err := do.Invoke[rbacpb.RBACServiceClient](i)
	if err != nil {
		return nil, err
	}

	IdentityClient, err := do.Invoke[pb.IdentityServiceClient](i)
	if err != nil {
		return nil, err
	}

	UserClient, err := do.Invoke[pb.UserServiceClient](i)
	if err != nil {
		return nil, err
	}

	return &GroupService{
		GroupRepository:   groupRepository,
		CuratorRepository: curatorRepository,
		MemberRepository:  memberRepository,
		RBACClient:        RBACClient,
		IdentityClient:    IdentityClient,
		UserClient:        UserClient,
	}, nil
}
