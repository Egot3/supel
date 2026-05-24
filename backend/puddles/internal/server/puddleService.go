package server

import (
	ppb "github.com/Egot3/supel/backend/contracts/puddle"
	"github.com/egot3/supel/backend/puddles/internal/database/repositories/member"
	"github.com/egot3/supel/backend/puddles/internal/database/repositories/moderator"
	"github.com/egot3/supel/backend/puddles/internal/database/repositories/puddle"
	grpcutils "github.com/egot3/supel/backend/puddles/internal/grpcUtils"
	"github.com/egot3/supel/backend/puddles/internal/services"
	"github.com/egot3/supel/backend/puddles/internal/token"
	"github.com/samber/do/v2"
)

type PuddleService struct {
	ppb.UnimplementedPuddleServiceServer

	moderatorRepository moderator.ModeratorRepsoitory
	memberRepository    member.MemberRepository
	puddleRepository    puddle.PuddleRepository

	grpcClient services.Client
	store      token.Store
	su         grpcutils.ServerUtils
}

func NewPuddleService(i do.Injector) (PuddleService, error) {
	moRepo := do.MustInvoke[moderator.ModeratorRepsoitory](i)
	meRepo := do.MustInvoke[member.MemberRepository](i)
	puRepo := do.MustInvoke[puddle.PuddleRepository](i)

	store := do.MustInvoke[token.Store](i)
	grpcClient := do.MustInvoke[services.Client](i)
	su := do.MustInvoke[grpcutils.ServerUtils](i)
	return PuddleService{
		moderatorRepository: moRepo,
		memberRepository:    meRepo,
		puddleRepository:    puRepo,

		grpcClient: grpcClient,
		store:      store,
		su:         su,
	}, nil
}
