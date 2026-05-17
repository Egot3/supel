package repositories

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/rbac/internal/models"
	"github.com/egot3/supel/backend/rbac/types"
	"github.com/google/uuid"
)

type RoleRepository interface {
	HasPermission(ctx context.Context, userUUID uuid.UUID, scope string, subScope *string, verb types.Verb) (bool, error)
	DeleteRole(ctx context.Context, roleUUID uuid.UUID) error
	CreateRole(ctx context.Context, role models.Role, actionUUIDs []uuid.UUID) error
	AssignorsAssignees(ctx context.Context, assignorUUID uuid.UUID) (assignees []uuid.UUID, err error)
	RoleAssignor(ctx context.Context, roleUUID, assigneeUUID uuid.UUID) (assignorUUID uuid.UUID, err error)
	RoleAssignees(ctx context.Context, roleUUID uuid.UUID) (userUUIDs []uuid.UUID, err error)
	AssignRole(ctx context.Context, assigneeUUID, assignorUUID, roleUUID uuid.UUID, expiresAt *time.Time) error
	UserRolePriority(ctx context.Context, userUUID string) (priority int16, err error)
	PatchRole(ctx context.Context, patchedModel models.PatchedRole) error
}
