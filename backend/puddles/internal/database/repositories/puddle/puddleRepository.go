package puddle

import (
	"context"

	"github.com/egot3/supel/backend/puddles/internal/models"
	"github.com/google/uuid"
)

type PuddleRepository interface {
	Puddle(ctx context.Context, puddleUUID uuid.UUID) (*models.Puddle, error)
	CreateOneOnOnePuddle(ctx context.Context, name string, description *string, startingUserUUIDs [2]uuid.UUID) error
	CreateGroup(ctx context.Context, name string, description *string, startingUserUUIDs []uuid.UUID, ownerUUID uuid.UUID) error
	CreateChannel(ctx context.Context, name string, description *string, ownerUUID uuid.UUID) error
	PatchPuddle(ctx context.Context, patched models.PuddlePatched) error
	DeletePuddle(ctx context.Context, puddleUUID uuid.UUID) error
	ListPuddleMembers(ctx context.Context, puddleUUID uuid.UUID, page, size uint32) (uuid.UUIDs, error)
}
