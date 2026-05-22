package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CuratorsHierarchy struct {
	bun.BaseModel `bun:"table:curators_hierarchy"`

	SeniorUUID      uuid.UUID `bun:"senior_uuid,notnull,pk"`
	SubordinateUUID uuid.UUID `bun:"subordinate_uuid,notnull,pk"`

	CreatedAt time.Time `bun:"created_at,default:NOW(),notnull"`
}
