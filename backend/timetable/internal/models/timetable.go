package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Timetable struct {
	bun.BaseModel `bun:"tabletimetables:"`

	UUID      uuid.UUID  `bun:"uuid,type:uuid,pk"`
	GroupUUID uuid.UUID  `bun:"group_uuid,type:uuid,notnull"`
	Name      string     `bun:"name,notnull"`
	AssignAt  *time.Time `bun:"assign_at"`
	RevokeAt  *time.Time `bun:"revoke_at"`

	UpdatedAt time.Time `bun:"updated_at,default:NOW(),notnull"`
	CreatedAt time.Time `bun:"created_at,default:NOW(),notnull"`
}

var _ bun.BeforeAppendModelHook = (*Timetable)(nil)

func (t *Timetable) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		t.CreatedAt = time.Now()
		t.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		t.UpdatedAt = time.Now()
	}
	return nil
}

type TimetablePatched struct {
	UUID             uuid.UUID
	GroupUUID        uuid.UUID //uuid.Nil exists
	Name             *string
	AssignAtUpdated  bool
	RevokeAtUpdatged bool
	AssignAt         *time.Time
	RevokeAt         *time.Time
}
