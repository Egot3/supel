package repositories

import (
	"context"

	"github.com/egot3/supel/backend/rbac/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunActionRepository struct {
	db *bun.DB
}

func NewActionRepository(i do.Injector) (ActionRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunActionRepository{db: db}, nil
}

func (r *bunActionRepository) AddActions(ctx context.Context, actionUUIDs []uuid.UUID, roleUUID uuid.UUID) error {
	if len(actionUUIDs) == 0 {
		return nil
	}

	rolesActionsUUIDs := make([]models.RolesActions, len(actionUUIDs))
	for i, actionUUID := range actionUUIDs {
		rolesActionsUUIDs[i] = models.RolesActions{
			RoleUUID:   roleUUID,
			ActionUUID: actionUUID,
		}
	}

	_, err := r.db.NewInsert().Ignore().Model(&rolesActionsUUIDs).Exec(ctx)
	return err
}

func (r *bunActionRepository) RevokeActions(ctx context.Context, actionUUIDs []uuid.UUID, roleUUID uuid.UUID) error {
	if len(actionUUIDs) == 0 {
		return nil
	}

	_, err := r.db.NewDelete().Model((*models.RolesActions)(nil)).
		Where("role_uuid = ?", roleUUID).
		Where("action_uuid IN (?)", bun.List(actionUUIDs)).
		Exec(ctx)
	return err
}

func (r *bunActionRepository) ListActions(ctx context.Context, page, size int) ([]models.Action, int, error) {
	actions := []models.Action{}
	total, err := r.db.NewSelect().Model(&actions).Limit(size).Offset(size * page).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return actions, total, nil
}

func (r *bunActionRepository) ActionsByRole(ctx context.Context, roleUUID uuid.UUID) ([]models.Action, int, error) {
	actions := []models.Action{}
	total, err := r.db.NewSelect().
		TableExpr("actions AS a").
		Join("JOIN roles_actions AS ra ON ra.action_uuid = a.action_uuid").
		Where("ra.role_uuid = ?", roleUUID).
		ScanAndCount(ctx, &actions)
	if err != nil {
		return nil, 0, err
	}

	return actions, total, nil
}

func (r *bunActionRepository) Action(ctx context.Context, actionUUID uuid.UUID) (*models.Action, error) {
	action := models.Action{}
	err := r.db.NewSelect().Model(&action).Where("action_uuid = ?", actionUUID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &action, nil
}
