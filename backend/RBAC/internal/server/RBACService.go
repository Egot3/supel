package server

import (
	"context"

	rbacpb "github.com/Egot3/supel/backend/contracts/rbac"
	"github.com/egot3/supel/backend/rbac/internal/database/repositories"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"google.golang.org/grpc/metadata"
)

type RBACService struct {
	rbacpb.UnimplementedRBACServiceServer
	ActionRepository repositories.ActionRepository
	RoleRepository   repositories.RoleRepository
}

func UserFromContext(ctx context.Context) (userID uuid.UUID, role string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, "", false
	}
	rawUserID := md.Get("user-uuid")[0]
	userID, err := uuid.Parse(rawUserID)
	role = md.Get("user-role")[0]
	return userID, role, !(len(userID) == 0 && len(role) == 0 && err == nil)
}

func NewRBACService(i do.Injector) (*RBACService, error) {
	actionRepo, err := do.Invoke[repositories.ActionRepository](i)
	if err != nil {
		return nil, err
	}

	roleRepository, err := do.Invoke[repositories.RoleRepository](i)
	if err != nil {
		return nil, err
	}

	return &RBACService{
		ActionRepository: actionRepo,
		RoleRepository:   roleRepository,
	}, nil
}
