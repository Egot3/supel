package moderator

import (
	"context"

	"github.com/egot3/supel/backend/puddles/internal/models"
	"github.com/google/uuid"
)

type ModeratorRepsoitory interface {
	AssignModerator(ctx context.Context, puddleUUID, memberUUID uuid.UUID, assignor *uuid.UUID) error
	IsModerator(ctx context.Context, puddleUUID, userUUID uuid.UUID) (bool, error)
	RevokeModerator(ctx context.Context, puddleUUID, moderatorUUID uuid.UUID) error
	AssignorsModerator(ctx context.Context, puddleUUID, assignorUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error)
	Moderator(ctx context.Context, puddleUUID, moderatorUUID uuid.UUID) (*models.PuddlesModerators, error)
}
