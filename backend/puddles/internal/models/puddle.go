package models

import (
	"context"
	"time"

	"github.com/egot3/supel/backend/puddles/internal/types"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Puddle struct {
	bun.BaseModel `bun:"puddles"`

	UUID        uuid.UUID        `bun:"uuid,type:uuid,pk,default:gen_random_uuid()"`
	Name        string           `bun:"name,notnull"`
	Description *string          `bun:"description"`
	PuddleType  types.PuddleType `bun:"puddle_type,notnull"`

	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	UpdatedAt time.Time  `bun:"updated_at,default:NOW(),notnull"`
	CreatedAt time.Time  `bun:"created_at,default:NOW(),notnull"`
}

type PuddlePatched struct {
	UUID        uuid.UUID
	Name        *string
	Description *string
}

var _ bun.BeforeAppendModelHook = (*Puddle)(nil)

func (g *Puddle) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		g.CreatedAt = time.Now()
		g.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		g.UpdatedAt = time.Now()
	}
	return nil
}
