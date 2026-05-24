package moderator

import (
	"context"

	"github.com/egot3/supel/backend/puddles/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunModeratorRepository struct {
	db *bun.DB
}

func NewModeratorRepository(i do.Injector) (ModeratorRepsoitory, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunModeratorRepository{db: db}, nil
}

func (r *bunModeratorRepository) AssignModerator(ctx context.Context, puddleUUID, memberUUID uuid.UUID, assignor *uuid.UUID) error {
	_, err := r.db.NewInsert().
		Model(&models.PuddlesModerators{PuddleUUID: puddleUUID, ModeratorUUID: memberUUID, AssignorUUID: assignor}).
		Exec(ctx)
	return err
}

func (r *bunModeratorRepository) IsModerator(ctx context.Context, puddleUUID, userUUID uuid.UUID) (bool, error) {
	is, err := r.db.NewSelect().Model(&models.PuddlesModerators{ModeratorUUID: userUUID, PuddleUUID: puddleUUID}).
		WherePK().Exists(ctx)
	if err != nil {
		return false, err
	}

	return is, nil
}

func (r *bunModeratorRepository) RevokeModerator(ctx context.Context, puddleUUID, moderatorUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.PuddlesModerators{PuddleUUID: puddleUUID, ModeratorUUID: moderatorUUID}).
		WherePK().Exec(ctx)
	return err
}

func (r *bunModeratorRepository) ListAssignorModerators(ctx context.Context, puddleUUID, assignorUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error) {
	var moderatorUUIDs uuid.UUIDs
	total, err := r.db.NewSelect().Model((*models.PuddlesModerators)(nil)).
		Where("assignor_uuid = ?", assignorUUID).Where("puddle_uuid = ?", puddleUUID).
		Limit(int(size)).Offset(int(size*page)).Column("moderator_uuid").ScanAndCount(ctx, &moderatorUUIDs)
	if err != nil {
		return nil, 0, err
	}

	return moderatorUUIDs, total, nil
}

func (r *bunModeratorRepository) Moderator(ctx context.Context, puddleUUID, moderatorUUID uuid.UUID) (*models.PuddlesModerators, error) {
	moderator := models.PuddlesModerators{
		PuddleUUID:    puddleUUID,
		ModeratorUUID: moderatorUUID,
	}
	err := r.db.NewSelect().Model(&moderator).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &moderator, nil
}
