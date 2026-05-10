package models

import (
	"context"
	"time"

	"github.com/Egot3/supel/backend/timetable/internal/types"
	"github.com/uptrace/bun"
)

type Period struct {
	bun.BaseModel `bun:"table:periods"`

	UUID string `bun:"uuid,pk"`

	Position   uint16    `bun:"position"`
	WeekNumber uint16    `bun:"week_number,notnull"`
	DayOfWeek  types.Day `bun:"day_of_week,notnull"`
	Year       uint16    `bun:"year,notnull"`

	Start time.Time `bun:"start,unique,notnull"`
	End   time.Time `bun:"end,unique,notnull"`

	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()"`
}

func (p *Period) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		p.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		p.UpdatedAt = time.Now()
	}
	return nil
}

type PatchedPeriod struct {
	UUID string

	WeekNumber *uint16
	DayOfWeek  *types.Day
	Year       *uint16
	Start      *time.Time
	End        *time.Time
}
