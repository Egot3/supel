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

	return bunGroupRepository{db: db}, nil
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
	err := r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		insert := models.Group{
			Name:        name,
			Description: description,
			GroupType:   groupType,
		}
		_, err := tx.NewInsert().Model(&insert).Exec(ctx)
		if err != nil {
			return err
		}

		/* getting curatee there*/

		return nil
	})
	return err
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

func (r *bunGroupRepository) DeleteGroup(ctx context.Context, groupUUID uuid.UUID) error
