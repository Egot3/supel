package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Period struct {
	bun.BaseModel `bun:"table:periods"`

	UUID     uuid.UUID `bun:"type:uuid,uuid,pk"`
	Name     string    `bun:"name,notnull"`
	Position int16     `bun:"position,notnull,default:0"`

	StartTime time.Time `bun:"start_time,type:time"`
	EndTime   time.Time `bun:"end_time,type:time"`

	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	UpdatedAt time.Time  `bun:"updated_at,default:NOW(),notnull"`
	CreatedAt time.Time  `bun:"created_at,default:NOW(),notnull"`
}

var _ bun.BeforeAppendModelHook = (*Period)(nil)

func (p *Period) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		p.CreatedAt = time.Now()
		p.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		p.UpdatedAt = time.Now()
	}
	return nil
}

type PeriodPatched struct {
	UUID      uuid.UUID
	Name      *string
	Position  *int32
	StartTime *time.Time
	EndTime   *time.Time
}
