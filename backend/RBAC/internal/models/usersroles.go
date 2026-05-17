package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserRoles struct {
	bun.BaseModel `bun:"table:users_roles"`

	UserUUID     uuid.UUID `bun:"user_uuid,notnull,pk"`
	RoleUUID     uuid.UUID `bun:"role_uuid,pk"`
	AssignorUUID uuid.UUID `bun:"assignor_uuid"`

	Role *Role `bun:"rel:belongs-to,join:role_uuid=role_uuid"`

	AssignedAt time.Time  `bun:"assigned_at,default:now()"`
	ExpiresAt  *time.Time `bun:"expires_at"`
}

//no user in microservice = insane JOINs
