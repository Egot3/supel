package member

import (
	"context"
	"database/sql"
	"errors"

	"github.com/egot3/supel/backend/puddles/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunMemberRepository struct {
	db *bun.DB
}

func NewMemberRepository(i do.Injector) (MemberRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunMemberRepository{db: db}, nil
}

func (r *bunMemberRepository) AddMember(ctx context.Context, memberUUID uuid.UUID, adderUUID *uuid.UUID, puddleUUID uuid.UUID) error {
	_, err := r.db.NewInsert().
		Model(&models.PuddlesMembers{PuddleUUID: puddleUUID, AdderUUID: adderUUID, MemberUUID: memberUUID}).Exec(ctx)
	return err
}

func (r *bunMemberRepository) ListMembersPuddles(ctx context.Context, memberUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error) {
	var puddlesUUIDs uuid.UUIDs
	total, err := r.db.NewSelect().Model((*models.PuddlesMembers)(nil)).Where("member_uuid = ?", memberUUID).
		Limit(int(size)).Offset(int(size*page)).
		Column("puddle_uuid").ScanAndCount(ctx, &puddlesUUIDs)
	if err != nil {
		return nil, 0, err
	}

	return puddlesUUIDs, total, nil
}

func (r *bunMemberRepository) ListMembersPuddlesIntersections(ctx context.Context, memberAUUID, memberBUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error) {
	var puddlesUUIDs uuid.UUIDs

	memberBSubq := r.db.NewSelect().Model((*models.PuddlesMembers)(nil)).Where("member_uuid = ?", memberBUUID).
		Column("puddle_uuid")

	total, err := r.db.NewSelect().Model((*models.PuddlesMembers)(nil)).
		Where("member_uuid = ?", memberAUUID).
		Column("puddle_uuid").
		Join("INNER JOIN (?) AS sub ON sub.puddle_uuid = puddles_members.puddle_uuid", memberBSubq).
		Order("puddle_uuid DESC").
		Limit(int(size)).Offset(int(size*page)).
		ScanAndCount(ctx, &puddlesUUIDs)
	if err != nil {
		return nil, 0, err
	}

	return puddlesUUIDs, total, nil
}

func (r *bunMemberRepository) ListAddersAddors(ctx context.Context, adderUUID uuid.UUID, page, size uint32) (uuid.UUIDs, int, error) {
	var memberUUIDs uuid.UUIDs
	total, err := r.db.NewSelect().
		Model((*models.PuddlesMembers)(nil)).Where("adder_uuid = ?", adderUUID).
		ScanAndCount(ctx, &memberUUIDs)
	if err != nil {
		return nil, 0, err
	}

	return memberUUIDs, total, nil
}

func (r *bunMemberRepository) PuddleMember(ctx context.Context, puddleUUID, memberUUID uuid.UUID) (*models.PuddlesMembers, error) {
	pMembers := models.PuddlesMembers{
		PuddleUUID: puddleUUID,
		MemberUUID: memberUUID,
	}
	err := r.db.NewSelect().Model(&pMembers).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &pMembers, nil
}

func (r *bunMemberRepository) RemoveMember(ctx context.Context, puddleUUID, memberUUID uuid.UUID) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().Model(&models.PuddlesMembers{PuddleUUID: puddleUUID, MemberUUID: memberUUID}).
			WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().
			Model(&models.PuddlesModerators{PuddleUUID: puddleUUID, ModeratorUUID: memberUUID}).
			WherePK().Exec(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}

			return err
		}

		return nil
	})
}
