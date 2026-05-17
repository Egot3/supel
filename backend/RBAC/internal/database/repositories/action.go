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

func NewActionRepository(i do.Injector) (RoleRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunRoleRepository{db: db}, nil
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
