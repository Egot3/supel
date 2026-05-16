package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type RolesActions struct {
	bun.BaseModel `bun:"table:roles_actions"`

	ActionUUID uuid.UUID `bun:"action_uuid,pk"`
	Action     *Action   `bun:"rel:belongs-to,join:action_uuid=action_uuid"`

	RoleUUID uuid.UUID `bun:"role_uuid,pk"`
	Role     *Role     `bun:"rel:belongs-to,join:role_uuid=role_uuid"`
}
