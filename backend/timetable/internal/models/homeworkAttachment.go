package models

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type HomeworkAttachment struct {
	bun.BaseModel `bun:"table:homework_attachments"`

	FileUUID     string `bun:"file_uuid,pk"`
	ConcreteUUID string `bun:"concrete_uuid,notnull"`

	Name string `bun:"name,notnull"`
	Mime string `bun:"mime,notnull"`

	StorageKey string    `bun:"storage_key,notnull"`
	CreatedAt  time.Time `bun:"created_at,notnull,default:now()"`

	_ struct{} `bun:"unique:uuid_idx"`
}

func (h *HomeworkAttachment) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		h.CreatedAt = time.Now()
	}
	return nil
}
