package curator

import (
	"context"
	"database/sql"

	"github.com/egot3/supel/backend/group/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunCuratorRepository struct {
	db *bun.DB
}

func NewGroupRepository(i do.Injector) (CuratorRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunCuratorRepository{db: db}, nil
}

func (r *bunCuratorRepository) AssignCuratorToSenior(ctx context.Context, seniorUUID, subordinateUUID uuid.UUID) error {
	_, err := r.db.NewInsert().
		Model(&models.CuratorsHierarchy{SeniorUUID: seniorUUID, SubordinateUUID: subordinateUUID}).
		Exec(ctx)
	return err
}

func (r *bunCuratorRepository) RevokeCuratorFromSenior(ctx context.Context, seniorUUID, subordinateUUID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model(&models.CuratorsHierarchy{SeniorUUID: seniorUUID, SubordinateUUID: subordinateUUID}).
		WherePK().Exec(ctx)
	return err
}

func (r *bunCuratorRepository) AssignCuratorToGroup(ctx context.Context, curatorUUID, groupUUID uuid.UUID) error {
	_, err := r.db.NewInsert().
		Model(&models.GroupsCurators{CuratorUUID: curatorUUID, GroupUUID: groupUUID}).
		Exec(ctx)
	return err
}

func (r *bunCuratorRepository) RevokeCuratorFromGroup(ctx context.Context, curatorUUID, groupUUID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model(&models.GroupsCurators{CuratorUUID: curatorUUID, GroupUUID: groupUUID}).
		WherePK().
		Exec(ctx)
	return err
}

func (r *bunCuratorRepository) AddCurator(ctx context.Context, seniorUUID, curatorUUID, groupUUID uuid.UUID) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().
			Model(&models.GroupsCurators{CuratorUUID: curatorUUID, GroupUUID: groupUUID}).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().
			Model(&models.GroupsCurators{CuratorUUID: curatorUUID, GroupUUID: groupUUID}).
			Exec(ctx)
		return err
	})
}

func (r *bunCuratorRepository) CanEdit(ctx context.Context, curatorUUID, groupUUID uuid.UUID) (bool, error) {
	subAnchor := r.db.NewSelect().
		TableExpr("curators_hierarchy AS ch").
		Where("ch.senior_uuid = ?", curatorUUID).
		ColumnExpr("ch.subordinate_uuid")

	subRecursive := r.db.NewSelect().
		TableExpr("curators_hierarchy AS ch").
		Join("JOIN subordinates AS s ON s.subordinate_uuid = ch.senior_uuid").
		ColumnExpr("ch.subordinate_uuid")
	subordinatesCTE := subAnchor.UnionAll(subRecursive)

	can, err := r.db.NewSelect().
		WithRecursive("subordinates", subordinatesCTE).
		TableExpr("groups_curators AS gc").
		Where("gc.group_uuid = ?", groupUUID).
		WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("gc.curator_uuid = ?", curatorUUID).
				Where("gc.curator_uuid IN (SELECT subordinate_uuid FROM subordinates)")
		}).
		ColumnExpr("1").Exists(ctx)

	if err != nil {
		return false, err
	}

	return can, nil
}

// the first and only danger-detection in this project
func (r *bunCuratorRepository) WillCycle(ctx context.Context, seniorUUID, subordinateUUID uuid.UUID) (bool, error) {
	will, err := r.db.NewSelect().Table("subordinates").WithRecursive("subordinates",
		r.db.NewSelect().
			TableExpr("curators_hierarchy AS ch").
			Where("ch.senior_uuid = ?", subordinateUUID).
			ColumnExpr("ch.subordinate_uuid").UnionAll(r.db.NewSelect().
			TableExpr("curators_hierarchy AS ch").
			Join("JOIN subordinates AS s ON s.subordinate_uuid = ch.senior_uuid").
			ColumnExpr("ch.subordinate_uuid"))).
		Where("subordinate_uuid = ?", seniorUUID).Exists(ctx)

	if err != nil {
		return true, err
	}

	return will, nil
}

func (r *bunCuratorRepository) RevokeCurator(ctx context.Context, curatorUUID uuid.UUID) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		senTable := tx.NewSelect().Table("curators_hierarchy").Column("senior_uuid").Where("subordinate_uuid = ?", curatorUUID)
		subTable := tx.NewSelect().Table("curators_hierarchy").Column("subordinate_uuid").Where("senior_uuid = ?", curatorUUID)

		_, err := tx.NewInsert().
			Model((*models.CuratorsHierarchy)(nil)).
			With("sen", senTable).
			With("sub", subTable).
			Column("senior_uuid", "subordinate_uuid").
			TableExpr("sen, sub").
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewDelete().Model((*models.CuratorsHierarchy)(nil)).WhereGroup(" OR ", func(dq *bun.DeleteQuery) *bun.DeleteQuery {
			return dq.Where("senior_uuid = ?", curatorUUID).Where("subordinate_uuid = ?", curatorUUID)
		}).Exec(ctx)
		return err
	})
}

func (r *bunCuratorRepository) GroupsCurators(ctx context.Context, groupUUID uuid.UUID) (uuid.UUIDs, error) {
	var curUUIDs uuid.UUIDs
	err := r.db.NewSelect().
		Model((*models.GroupsCurators)(nil)).
		Where("group_uuid = ?", groupUUID).Column("curator_uuid").Scan(ctx, &curUUIDs)
	if err != nil {
		return nil, err
	}

	return curUUIDs, nil
}

func (r *bunCuratorRepository) IsCurator(ctx context.Context, userUUID uuid.UUID) (bool, error) {
	is, err := r.db.NewSelect().
		Model((*models.CuratorsHierarchy)(nil)).
		WhereGroup(" OR ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.WhereOr("subordinate_uuid = ?", userUUID).WhereOr("senior_uuid = ?", userUUID)
		}).Exists(ctx)

	if err != nil {
		return false, err
	}

	return is, nil
}
