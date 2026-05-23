package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GroupsCurators struct {
	bun.BaseModel `bun:"table:groups_curators"`

	GroupUUID uuid.UUID `bun:"group_uuid,type:uuid,pk"`
	Group     *Group    `bun:"rel:belongs-to,join:group_uuid=uuid"`

	CuratorUUID uuid.UUID `bun:"curator_uuid,type:uuid,pk"`

	AssignedAt time.Time `bun:"assigned_at,default:NOW(),notnull"`
}
