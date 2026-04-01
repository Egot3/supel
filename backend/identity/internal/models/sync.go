package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Sync struct {
	bun.BaseModel `bun:"table:sync_state"` //во имя святого ORM
	Source        string                   `bun:"source_name,pk"`
	LastSync      time.Time                `bun:"last_synced_at,notnull"`
	LastMessageId string                   `bun:"last_message_id"`

	_ struct{} `bun:"unique:uuid_idx,type:unique"`
}
