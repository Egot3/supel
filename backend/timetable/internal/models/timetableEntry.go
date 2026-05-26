package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TimetableEntry struct {
	bun.BaseModel `bun:"table:timetable_entries,alias:te"`

	UUID          uuid.UUID `bun:"uuid,type:uuid,pk"`
	TimetableUUID uuid.UUID `bun:"timetable_uuid,type:uuid,notnull,"`
	PeriodUUID    uuid.UUID `bun:"period_uuid,type:uuid,notnull"`
	DayOfWeek     int16     `bun:"day_of_week,notnull"`
	SubjectUUID   uuid.UUID `bun:"subject_uuid,type:uuid,notnull"`
	Place         string    `bun:"place,notnull"`
	TeacherUUID   uuid.UUID `bun:"teacher_uuid,type:uuid"` //nullable

	Subject   Subject   `bun:"rel:belongs-to,join:uuid=subject_uuid"`
	Period    Period    `bun:"rel:belongs-to,join:uuid=period_uuid"`
	Timetable Timetable `bun:"rel:belongs-to,join:uuid=timetable_uuid"`

	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	UpdatedAt time.Time  `bun:"updated_at,default:NOW(),notnull"`
	CreatedAt time.Time  `bun:"created_at,default:NOW(),notnull"`
}

var _ bun.BeforeAppendModelHook = (*TimetableEntry)(nil)

func (s *TimetableEntry) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		s.CreatedAt = time.Now()
		s.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		s.UpdatedAt = time.Now()
	}
	return nil
}

type TimetableEntryPatched struct {
	UUID          uuid.UUID
	TimetableUUID uuid.UUID
	PeriodUUID    uuid.UUID
	DayOfWeek     int16 //0 == nil
	SubjectUUID   uuid.UUID
	Place         *string
	TeacherUUID   uuid.UUID
}
