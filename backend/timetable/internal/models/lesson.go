package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Lesson struct {
	bun.BaseModel `bun:"table:lessons,alias:l"`

	UUID               uuid.UUID `bun:"uuid,pk,type:uuid"`
	TimetableEntryUUID uuid.UUID `bun:"timetable_entry_uuid,type:uuid,unique"`
	Date               time.Time `bun:"date,notnull,unique"`
	Cancelled          bool      `bun:"cancelled,type:boolean,default:false"`

	CreatedAt time.Time `bun:"created_at,default:NOW(),notnull"`

	TimetableEntry     TimetableEntry       `bun:"rel:belongs-to,join:uuid=timetable_entry_uuid"`
	HomeworkAttachment []HomeworkAttachment `bun:"rel:has-many,join:uuid=lesson_uuid"`
}

var _ bun.BeforeAppendModelHook = (*Lesson)(nil)

func (l *Lesson) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		l.CreatedAt = time.Now()
	}
	return nil
}

type LessonPatched struct {
	UUID               uuid.UUID
	Date               *time.Time
	Cancelled          *bool
	TimetableEntryUUID uuid.UUID
}
