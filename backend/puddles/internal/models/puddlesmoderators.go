package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PuddlesModerators struct {
	bun.BaseModel `bun:"table:puddles_moderators"`

	PuddleUUID uuid.UUID `bun:"puddle_uuid,type:uuid,pk"`
	Puddle     *Puddle   `bun:"rel:belongs-to,join:puddle_uuid=uuid"`

	ModeratorUUID uuid.UUID  `bun:"moderator_uuid,type:uuid,pk"`
	AssignorUUID  *uuid.UUID `bun:"assignor_uuid,type:uuid"`

	AssignedAt time.Time `bun:"assigned_at,default:NOW(),notnull"`
}

var _ bun.BeforeAppendModelHook = (*PuddlesModerators)(nil)

func (g *PuddlesModerators) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		g.AssignedAt = time.Now()
	}
	return nil
}
