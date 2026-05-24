package puddle

import (
	"context"
	"database/sql"
	"time"

	"github.com/egot3/supel/backend/puddles/internal/carefulness"
	"github.com/egot3/supel/backend/puddles/internal/models"
	"github.com/egot3/supel/backend/puddles/internal/types"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunPuddleRepository struct {
	db *bun.DB
}

func NewPuddleRepository(i do.Injector) (PuddleRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunPuddleRepository{db: db}, nil
}

func (r *bunPuddleRepository) Puddle(ctx context.Context, puddleUUID uuid.UUID) (*models.Puddle, error) {
	var puddle models.Puddle
	err := r.db.NewSelect().Model(&puddle).Where("puddle_uuid = ?", puddleUUID).WhereAllWithDeleted().Scan(ctx)
	if err != nil {
		return nil, err
	}
	if puddle.DeletedAt != nil && puddle.DeletedAt.Before(time.Now()) {
		return nil, carefulness.Gone
	}

	return &puddle, nil
}

func (r *bunPuddleRepository) CreateOneOnOnePuddle(ctx context.Context, name string, description *string, startingUserUUIDs [2]uuid.UUID) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		createdPuddle := models.Puddle{
			Name:        name,
			Description: description,
			PuddleType:  types.ONEONONE,
		}
		err := tx.NewInsert().Model(&createdPuddle).Returning("uuid").Scan(ctx)
		if err != nil {
			return err
		}

		var assignMemberModels [2]models.PuddlesMembers
		var assignModeratorModels [2]models.PuddlesModerators
		for i, memberUUID := range startingUserUUIDs {
			assignMemberModels[i] = models.PuddlesMembers{
				PuddleUUID: createdPuddle.UUID,
				MemberUUID: memberUUID,
			}
			assignModeratorModels[i] = models.PuddlesModerators{
				PuddleUUID:    createdPuddle.UUID,
				ModeratorUUID: memberUUID,
			}
		}

		res, err := tx.NewInsert().Model(&assignMemberModels).Exec(ctx)
		if err != nil {
			return err
		}
		if c, _ := res.RowsAffected(); c < 2 {
			return carefulness.Conflict
		}

		res, err = tx.NewInsert().Model(&assignModeratorModels).Exec(ctx)
		if err != nil {
			return err
		}
		if c, _ := res.RowsAffected(); c < 2 {
			return carefulness.Conflict
		}

		return nil
	})
}

func (r *bunPuddleRepository) CreateGroup(ctx context.Context, name string, description *string, startingUserUUIDs []uuid.UUID, ownerUUID uuid.UUID) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		createdPuddle := models.Puddle{
			Name:        name,
			Description: description,
			PuddleType:  types.GROUP,
		}
		err := tx.NewInsert().Model(&createdPuddle).Returning("uuid").Scan(ctx)
		if err != nil {
			return err
		}

		assignMemberModels := make([]models.PuddlesMembers, len(startingUserUUIDs))
		assignModeratorModels := models.PuddlesModerators{
			PuddleUUID:    createdPuddle.UUID,
			ModeratorUUID: ownerUUID,
		}
		for i, memberUUID := range startingUserUUIDs {
			assignMemberModels[i] = models.PuddlesMembers{
				PuddleUUID: createdPuddle.UUID,
				MemberUUID: memberUUID,
			}
		}

		_, err = tx.NewInsert().Model(&assignMemberModels).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().Model(&assignModeratorModels).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *bunPuddleRepository) CreateChannel(ctx context.Context, name string, description *string, ownerUUID uuid.UUID) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		createdPuddle := models.Puddle{
			Name:        name,
			Description: description,
			PuddleType:  types.CHANNEL,
		}
		err := tx.NewInsert().Model(&createdPuddle).Returning("uuid").Scan(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().Model(&models.PuddlesModerators{PuddleUUID: createdPuddle.UUID, ModeratorUUID: ownerUUID}).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *bunPuddleRepository) PatchPuddle(ctx context.Context, patched models.PuddlePatched) error {
	isChanged := false
	query := r.db.NewUpdate().Model((*models.Puddle)(nil)).Where("uuid = ?", patched.UUID)
	if patched.Name != nil {
		query = query.Set("name = ?", patched.Name)
		isChanged = true
	}
	if patched.Description != nil {
		query = query.Set("description = ?", patched.Description)
		isChanged = true
	}
	if !isChanged {
		return nil
	}

	_, err := query.Exec(ctx)

	return err
}

func (r *bunPuddleRepository) DeletePuddle(ctx context.Context, puddleUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model((*models.Puddle)(nil)).Where("uuid = ?", puddleUUID).Exec(ctx)
	return err
}

func (r *bunPuddleRepository) ListPuddleMembers(ctx context.Context, puddleUUID uuid.UUID, page, size uint32) (uuid.UUIDs, error) {
	var userUUIDs uuid.UUIDs
	err := r.db.NewSelect().Model((*models.PuddlesMembers)(nil)).
		Where("puddle_uuid = ?", puddleUUID).
		Limit(int(size)).Offset(int(size*page)).
		Column("member_uuid").Scan(ctx, &userUUIDs)
	if err != nil {
		return nil, err
	}

	return userUUIDs, nil
}

func (r *bunPuddleRepository) PuddleMemberCount(ctx context.Context, puddleUUID uuid.UUID) (int, error) {
	return r.db.NewSelect().Model((*models.PuddlesMembers)(nil)).Where("puddle_uuid = ?", puddleUUID).Count(ctx)
}
