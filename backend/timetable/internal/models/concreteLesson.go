package models

import (
	"context"
	"time"

	"github.com/Egot3/supel/backend/timetable/internal/types"
	"github.com/uptrace/bun"
)

type ConcreteLesson struct {
	bun.BaseModel `bun:"table:abstract_lessons"`

	ConcreteUUID string `bun:"concrete_uuid,pk,type:uuid"`
	AbstractUUID string `bun:"abstract_uuid,notnull"`
	TeacherUUID  string `bun:"teacher_uuid,notnull"`
	GroupUUID    string `bun:"group_uuid,notnull"`

	HomeworkBodyKey string `bun:"homework_body_key,type:character varying GENERATED ALWAYS AS ('orgs/ETSEvilCorp/timetable/homework/body/' || concrete_uuid) STORED,scanonly"`

	WeekNumber uint16    `bun:"week_number"`
	DayOfWeek  types.Day `bun:"day_of_week"`
	Period     uint16    `bun:"period"`
	Year       uint16    `bun:"year"`

	CreatedAt time.Time  `bun:"created_at,notnull,default:now()"`
	UpdatedAt time.Time  `bun:"updated_ay,notnull,default:now()"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`

	Attachments    []*HomeworkAttachment `bun:"rel:has-many,join:concrete_uuid=concrete_uuid"`
	AbstractLesson AbstractLesson        `bun:"rel:belongs-to,join:abstract_uuid=uuid"`

	_ struct{} `bun:"unique:uuid_idx"`
}

func (c *ConcreteLesson) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		c.CreatedAt = time.Now()
		c.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		c.UpdatedAt = time.Now()
	}
	return nil
}

type PatchConcreteLesson struct {
	ConcreteUUID string
	AbstractUUID *string
	TeacherUUID  *string
	GroupUUID    *string
	Period       *uint16
	DayOfWeek    *types.Day
	Year         *uint16
	WeekNumber   *uint16
}
