package server

import (
	grpb "github.com/Egot3/supel/backend/contracts/group"

	"github.com/egot3/supel/backend/group/internal/database/repositories/curator"
	"github.com/egot3/supel/backend/group/internal/database/repositories/group"
	"github.com/egot3/supel/backend/group/internal/database/repositories/member"
	grpcutils "github.com/egot3/supel/backend/group/internal/grpcUtils"
	"github.com/egot3/supel/backend/group/internal/services"
	"github.com/samber/do/v2"
)

type GroupService struct {
	grpb.UnimplementedGroupServiceServer
	GroupRepository   group.GroupRepository
	CuratorRepository curator.CuratorRepository
	MemberRepository  member.MemberRepository
	Client            services.Client
	su                grpcutils.ServerUtils
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

	Client := do.MustInvoke[services.Client](i)
	Su := do.MustInvoke[grpcutils.ServerUtils](i)

	return &GroupService{
		GroupRepository:   groupRepository,
		CuratorRepository: curatorRepository,
		MemberRepository:  memberRepository,
		Client:            Client,
		su:                Su,
	}, nil
}
