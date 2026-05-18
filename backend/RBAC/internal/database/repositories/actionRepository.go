package repositories

import (
	"context"

	"github.com/egot3/supel/backend/rbac/internal/models"
	"github.com/google/uuid"
)

type ActionRepository interface {
	AddActions(ctx context.Context, actionUUIDs []uuid.UUID, roleUUID uuid.UUID) error
	RevokeActions(ctx context.Context, actionUUIDs []uuid.UUID, roleUUID uuid.UUID) error
	ListActions(ctx context.Context, page, size int) ([]models.Action, int, error)
	ActionsByRole(ctx context.Context, roleUUID uuid.UUID) ([]models.Action, int, error)
	Action(ctx context.Context, actionUUID uuid.UUID) (*models.Action, error)
}
