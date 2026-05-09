package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type AbstractLesson struct {
	bun.BaseModel `bun:"table:abstract_lessons"`
	UUID          string `bun:"uuid,pk,type:uuid"`
	Name          string `bun:"name,notnull"`

	ConcreteLessons []*ConcreteLesson `bun:"rel:has-many,join:uuid=abstract_uuid"`

	CreatedAt time.Time  `bun:"created_at,notnull,default:now()"`
	UpdatedAt time.Time  `bun:"updated_ay,notnull,default:now()"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`

	_ struct{} `bun:"unique:uuid_idx"`
}

var _ bun.BeforeAppendModelHook = (*AbstractLesson)(nil)

func (a *AbstractLesson) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		a.CreatedAt = time.Now()
		a.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		a.UpdatedAt = time.Now()
	}
	return nil
}

type PatchAbstractLesson struct {
	UUID string
	Name *string
}
