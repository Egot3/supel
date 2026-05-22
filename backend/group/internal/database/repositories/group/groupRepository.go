package group

import (
	"context"

	"github.com/egot3/supel/backend/group/internal/models"
	"github.com/egot3/supel/backend/group/internal/types"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GroupRepository interface {
	Group(ctx context.Context, groupUUID uuid.UUID) (*models.Group, error)
	CreateGroup(ctx context.Context, requestor uuid.UUID, name string, description *string, groupType types.GroupType) error
	Search(ctx context.Context, sample string, limit int) ([]models.Group, error)
	DeleteGroup(ctx context.Context, groupUUID uuid.UUID) error
	ListGroups(ctx context.Context, page, size uint32, groupType types.GroupType, order bun.Order) ([]models.Group, uint64, error)
	PatchGroup(ctx context.Context, patched models.GroupPatched) error
	CuratorsGroups(ctx context.Context, curatorsUUID uuid.UUID) ([]models.Group, error)
}
