package models

import (
	"github.com/egot3/supel/backend/rbac/types"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Action struct {
	bun.BaseModel `bun:"table:actions"`

	UUID     uuid.UUID  `bun:"action_uuid,default:now(),pk"`
	Scope    string     `bun:"scope,notnull"`
	SubScope *string    `bun:"sub_scope"`
	Verb     types.Verb `bun:"verb,notnull"`
}
