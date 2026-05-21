package group

import (
	"context"
	"database/sql"
	"time"

	"github.com/egot3/supel/backend/group/internal/carefulness"
	"github.com/egot3/supel/backend/group/internal/models"
	"github.com/egot3/supel/backend/group/internal/types"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunGroupRepository struct {
	db *bun.DB
}

func NewGroupRepository(i do.Injector) (GroupRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunGroupRepository{db: db}, nil
}

func (r *bunGroupRepository) Group(ctx context.Context, groupUUID uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.NewSelect().Model(&group).Where("uuid = ?", groupUUID).WhereAllWithDeleted().Scan(ctx)
	if err != nil {
		return nil, err
	}

	if group.DeletedAt != nil && time.Since(*group.DeletedAt) > 0 {
		return nil, carefulness.Gone
	}

	return &group, err
}

func (r *bunGroupRepository) CreateGroup(ctx context.Context, requestor uuid.UUID, name string, description *string, groupType types.GroupType) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		insert := models.Group{
			Name:        name,
			Description: description,
			GroupType:   groupType,
		}
		_, err := tx.NewInsert().Model(&insert).Returning("uuid").Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().Model(&models.GroupsCurators{GroupUUID: insert.UUID, CuratorUUID: requestor}).Exec(ctx)

		return err
	})
}

func (r *bunGroupRepository) Search(ctx context.Context, sample string, limit int) ([]models.Group, error) {
	var groupBuckets []struct {
		models.Group
		Score float64 `bun:"score"`
	}

	err := r.db.NewSelect().
		TableExpr("groups").
		ColumnExpr("*, similarity(name, ?) AS score", sample).
		Where("name % ?", sample).
		OrderExpr("score DESC").
		Limit(limit).
		Scan(ctx, &groupBuckets)
	if err != nil {
		return nil, err
	}

	groups := make([]models.Group, len(groupBuckets))
	for i, bucket := range groupBuckets {
		groups[i] = bucket.Group
	}

	return groups, nil
}

func (r *bunGroupRepository) DeleteGroup(ctx context.Context, groupUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.Group{UUID: groupUUID}).WherePK().Exec(ctx)
	return err
}

func (r *bunGroupRepository) PatchGroup(ctx context.Context, patched models.GroupPatched) error {
	isUpdated := false

	query := r.db.NewUpdate().Model((*models.Group)(nil)).Where("uuid = ?", patched.UUID)
	if patched.Name != nil {
		query = query.Set("name = ?", *patched.Name)
		isUpdated = true
	}
	if patched.Description != nil {
		query = query.Set("description = ?", patched.Description)
		isUpdated = true
	}

	if !isUpdated {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

/* func (r *bunGroupRepository) UsersGroups(ctx context.Context, userUUID uuid.UUID) (uuid.UUIDs, error) {
	groups := uuid.UUIDs{}
	err := r.db.NewSelect().Model((*models.GroupsMembers)(nil)).Where("member_uuid = ?", userUUID).Column("group_uuid").Scan(ctx, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
} */

func (r *bunGroupRepository) ListGroups(ctx context.Context, page, size uint32, groupType types.GroupType, order bun.Order) ([]models.Group, uint64, error) {
	groups := []models.Group{}
	query := r.db.NewSelect().Model(&groups).Limit(int(size)).Offset(int(size*page)).Order("created_at", string(order))
	if groupType != types.UNSPECIFIED {
		query = query.Where("group_type = ?", groupType)
	}

	total, err := query.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return groups, uint64(total), nil
}
