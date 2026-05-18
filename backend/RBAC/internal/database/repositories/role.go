package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/egot3/supel/backend/rbac/internal/logctx"
	"github.com/egot3/supel/backend/rbac/internal/models"
	"github.com/egot3/supel/backend/rbac/types"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunRoleRepository struct {
	db *bun.DB
}

func NewRoleRepository(i do.Injector) (RoleRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunRoleRepository{db: db}, nil
}

func (r *bunRoleRepository) HasPermission(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb types.Verb) (bool, error) {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "repository"),
		slog.String("table", "users_roles,roles,actions"),
		slog.String("operation", "select"),
	)

	logger.InfoContext(ctx, "starting db request", slog.String("userUUID", userUUID.String()),
		slog.String("action", fmt.Sprintf("%v.%v:%v", scope, subScope, verb)),
	)
	start := time.Now()
	anchor := r.db.NewSelect().TableExpr("users_roles").
		ColumnExpr("users_roles.role_uuid").
		Where("users_roles.user_uuid = ?", userUUID).
		Where("users_roles.expires_at IS NULL OR users_roles.expires_at > now()").
		Join("JOIN roles AS role ON role.role_uuid = users_roles.role_uuid AND role.deleted_at IS NULL")

	recursive := r.db.NewSelect().TableExpr("roles").ColumnExpr("roles.extended_role_uuid").
		Join("JOIN role_deps AS rd ON rd.role_uuid = roles.role_uuid").
		Where("roles.extended_role_uuid IS NOT NULL").
		Where("roles.deleted_at IS NULL")

	cte := anchor.UnionAll(recursive)

	query := r.db.NewSelect().
		With("role_deps", cte).
		TableExpr("role_deps").
		Join("JOIN roles_actions AS ra ON ra.role_uuid = role_deps.role_uuid").
		Join("JOIN actions AS action ON action.action_uuid = ra.action_uuid").
		Where("action.scope = ?", scope).
		Where("action.verb = ?", verb)

	if subScope != nil {
		query = query.Where("action.sub_scope = ?", *subScope)
	} else {
		query = query.Where("action.sub_scope IS NULL")
	}

	can, err := query.ColumnExpr("1").Exists(ctx) //effieciency

	elapsed := time.Since(start)
	if err != nil {
		logger.ErrorContext(ctx, "db request failed",
			slog.Duration("elapsed", elapsed),
			slog.String("err", err.Error()),
		)
		return false, err
	}

	logger.InfoContext(ctx, "db request succeeded", slog.Duration("elapsed", elapsed))
	return can, nil
}

func (r *bunRoleRepository) UserRolePriority(ctx context.Context, userUUID uuid.UUID) (priority int16, err error) {
	err = r.db.NewSelect().
		TableExpr("users_roles AS user").
		Where("user.user_uuid = ?", userUUID).
		Join("JOIN roles AS role ON role.role_uuid = user.role_uuid").
		Column("role.priority").
		Order("role.priority DESC").Limit(1).Scan(ctx, &priority)
	if err != nil {
		return 0, err
	}

	return priority, nil
}

func (r *bunRoleRepository) RolePriority(ctx context.Context, roleUUID uuid.UUID) (priority int16, err error) {
	err = r.db.NewSelect().
		TableExpr("roles AS role").
		Where("role.role_uuid = ?", roleUUID).
		Column("role.priority").
		Scan(ctx, &priority)
	if err != nil {
		return 0, err
	}

	return priority, nil
}

func (r *bunRoleRepository) AssignRole(ctx context.Context, assigneeUUID, assignorUUID, roleUUID uuid.UUID, expiresAt *time.Time) error {
	assignedRole := models.UserRoles{AssignorUUID: assignorUUID, UserUUID: assigneeUUID, RoleUUID: roleUUID, ExpiresAt: expiresAt}

	_, err := r.db.NewInsert().Model(&assignedRole).Exec(ctx)

	return err
}

func (r *bunRoleRepository) RevokeRole(ctx context.Context, assigneeUUID, roleUUID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model(&models.UserRoles{
			UserUUID: assigneeUUID,
			RoleUUID: roleUUID,
		}).
		WherePK().
		Exec(ctx)
	return err
}

func (r *bunRoleRepository) RoleAssignees(ctx context.Context, roleUUID uuid.UUID) (userUUIDs uuid.UUIDs, err error) {
	err = r.db.NewSelect().
		Model((*models.UserRoles)(nil)).
		Where("role_uuid = ?", roleUUID).
		Column("user_uuid").
		Scan(ctx, &userUUIDs)
	if err != nil {
		return nil, err
	}

	return userUUIDs, nil
}

func (r *bunRoleRepository) RoleAssignor(ctx context.Context, roleUUID, assigneeUUID uuid.UUID) (assignorUUID uuid.UUID, err error) {
	err = r.db.
		NewSelect().
		Model((*models.RolesActions)(nil)).
		Where("role_uuid = ?", roleUUID).
		Where("user_uuid = ?", assigneeUUID).
		Column("assignor_uuid").
		Scan(ctx, &assignorUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return assignorUUID, nil
}

func (r *bunRoleRepository) AssignorsAssignees(ctx context.Context, assignorUUID uuid.UUID) (assignees uuid.UUIDs, err error) {
	err = r.db.NewSelect().
		Model((*models.UserRoles)(nil)).
		Where("assignor_uuid = ?", assignorUUID).
		WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.
				Where("expires_at > NOW()").
				WhereOr("expires_at IS NULL")
		}).
		Column("user_uuid").
		Scan(ctx, &assignees)
	if err != nil {
		return nil, err
	}

	return assignees, nil
}

func (r *bunRoleRepository) CreateRole(ctx context.Context, role models.Role, actionUUIDs []uuid.UUID) error {
	err := r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&role).Returning("role_uuid").Exec(ctx)
		if err != nil {
			return err
		}

		if len(actionUUIDs) == 0 {
			return nil
		}

		rolesActions := make([]models.RolesActions, len(actionUUIDs))
		for i, actionUUID := range actionUUIDs {
			rolesActions[i] = models.RolesActions{
				RoleUUID:   role.UUID,
				ActionUUID: actionUUID,
			}
		}

		_, err = tx.
			NewInsert().
			Model(&rolesActions).
			Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *bunRoleRepository) DeleteRole(ctx context.Context, roleUUID uuid.UUID) error {
	logger := logctx.LoggerFromContext(ctx).With(
		slog.String("layer", "repository"),
		slog.String("table", "roles_actions,roles"),
		slog.String("operation", "Delete"),
	)

	logger.InfoContext(ctx, "starting db request", slog.String("roleUUID", roleUUID.String()))

	_, err := r.db.NewDelete().Model(&models.Role{UUID: roleUUID}).Exec(ctx)
	return err
}

func (r *bunRoleRepository) PatchRole(ctx context.Context, patchedModel models.PatchedRole) error {
	query := r.db.NewUpdate().Model(&models.Role{UUID: patchedModel.UUID})
	if patchedModel.Name != nil {
		query = query.Set("name = ?", patchedModel.Name)
	}
	if patchedModel.Description != nil {
		query = query.Set("description = ?", patchedModel.Description)
	}
	if patchedModel.Priority != nil {
		query = query.Set("priority = ?", patchedModel.Priority)
	}
	if patchedModel.ExtendedUUID != nil {
		query = query.Set("extended_role_uuid = ?", patchedModel.ExtendedUUID)
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *bunRoleRepository) Role(ctx context.Context, roleUUID uuid.UUID) (*models.Role, error) {
	role := models.Role{}
	err := r.db.NewSelect().Model(&role).Where("role_uuid = ?", roleUUID).WhereAllWithDeleted().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *bunRoleRepository) ListRoles(ctx context.Context, page, size int) ([]models.Role, int, error) {
	var roles []models.Role

	total, err := r.db.NewSelect().
		Model(&roles).
		WhereAllWithDeleted().
		Limit(size).
		Offset(page * size).
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}
