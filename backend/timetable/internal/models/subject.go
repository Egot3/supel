package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Subject struct {
	bun.BaseModel `bun:"table:subjects"`

	UUID uuid.UUID `bun:"uuid,pk,type:uuid"`
	Name string    `bun:"name,notnull"`

	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	UpdatedAt time.Time  `bun:"updated_at,default:NOW(),notnull"`
	CreatedAt time.Time  `bun:"created_at,default:NOW(),notnull"`
}

var _ bun.BeforeAppendModelHook = (*Subject)(nil)

func (s *Subject) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		s.CreatedAt = time.Now()
		s.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		s.UpdatedAt = time.Now()
	}
	return nil
}

type SubjectPatched struct {
	UUID uuid.UUID
	Name *string //there is just nothing more to it now
}
