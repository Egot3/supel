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
	PatchRole(ctx context.Context, patchedModel models.PatchedRole) error
	Role(ctx context.Context, roleUUID uuid.UUID) (*models.Role, error)
	ListRoles(ctx context.Context, page, size int) ([]models.Role, int, error)
	CreateRole(ctx context.Context, role models.Role, actionUUIDs []uuid.UUID) error
	AssignorsAssignees(ctx context.Context, assignorUUID uuid.UUID) (assignees uuid.UUIDs, err error)
	RoleAssignor(ctx context.Context, roleUUID, assigneeUUID uuid.UUID) (assignorUUID uuid.UUID, err error)
	RoleAssignees(ctx context.Context, roleUUID uuid.UUID) (userUUIDs uuid.UUIDs, err error)
	AssignRole(ctx context.Context, assigneeUUID, assignorUUID, roleUUID uuid.UUID, expiresAt *time.Time) error
	RolePriority(ctx context.Context, roleUUID uuid.UUID) (priority int16, err error)
	UserRolePriority(ctx context.Context, userUUID uuid.UUID) (priority int16, err error)
	RevokeRole(ctx context.Context, assigneeUUID, roleUUID uuid.UUID) error
}
