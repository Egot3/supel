package models

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/group/internal/types"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Group struct {
	bun.BaseModel `bun:"groups"`

	UUID        uuid.UUID       `bun:"uuid,type:uuid,pk,default:gen_random_uuid()"`
	Name        string          `bun:"name,notnull"`
	Description *string         `bun:"description"`
	GroupType   types.GroupType `bun:"group_type,notnull"`

	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	UpdatedAt time.Time  `bun:"updated_at,default:NOW(),notnull"`
	CreatedAt time.Time  `bun:"created_at,default:NOW(),notnull"`
}

type GroupPatched struct {
	UUID        uuid.UUID
	Name        *string
	Description *string
	GroupType   types.GroupType
}

var _ bun.BeforeAppendModelHook = (*Group)(nil)

func (g *Group) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		g.CreatedAt = time.Now()
		g.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		g.UpdatedAt = time.Now()
	}
	return nil
}
