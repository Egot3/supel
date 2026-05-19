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

	UUID            uuid.UUID       `bun:"uuid,pk"`
	Name            string          `bun:"name,notnull"`
	Description     *string         `bun:"description"`
	GroupType       types.GroupType `bun:"group_type,notnull"`
	ParentGroupUUID uuid.UUID       `bun:"parent_group_uuid"`

	ParentGroup *Group `bun:"rel:belongs-to,join:parent_group_uuid=uuid"`

	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	UpdatedAt time.Time  `bun:"updated_at,default:NOW(),notnull"`
	CreatedAt time.Time  `bun:"created_at,default:NOW(),notnull"`
}

type GroupPatched struct {
	UUID        uuid.UUID
	Name        string
	Description *string
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
