package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles"`

	UUID         uuid.UUID  `bun:"role_uuid,pk"`
	Name         string     `bun:"role_name,unique,notnull"`
	Description  *string    `bun:"role_description"`
	ExtendedUUID *uuid.UUID `bun:"extended_role_uuid"`

	Priority int16 `bun:"priority,nullzero"`

	CreatedAt time.Time  `bun:"created_at,default:now()"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	UpdatedAt time.Time  `bun:"updated_at,default:now()"`

	ExtendedRole *Role     `bun:"rel:belongs-to,join:extended_role_uuid=role_uuid"`
	Actions      []*Action `bun:"m2m:roles_actions,join:role_uuid=role_uuid"`
}

var _ bun.BeforeAppendModelHook = (*Role)(nil)

func (r *Role) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		r.CreatedAt = time.Now()
		r.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		r.UpdatedAt = time.Now()
	}
	return nil
}

type PatchedRole struct {
	UUID         uuid.UUID
	Name         *string
	Description  *string
	ExtendedUUID *string
	Priority     *int16
}
