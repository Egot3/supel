package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type HomeworkAttachment struct {
	bun.BaseModel `bun:"homework_attachments,alias:ha"`

	Name       string    `bun:"name,notnull"`
	Mime       string    `bun:"mime,notnull"`
	LessonUUID uuid.UUID `bun:"lesson_uuid,type:uuid,notnull"`
	StorageKey string    `bun:"storage_key,notnull"`

	CreatedAt time.Time `bun:"created_at,default:NOW(),notnull"`

	Lesson Lesson `bun:"rel:belongs-to,join:uuid=lesson_uuid"`
}

var _ bun.BeforeAppendModelHook = (*HomeworkAttachment)(nil)

func (l *HomeworkAttachment) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		l.CreatedAt = time.Now()
	}
	return nil
}
