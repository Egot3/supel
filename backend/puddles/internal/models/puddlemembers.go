package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PuddlesMembers struct {
	bun.BaseModel `bun:"table:puddles_members"`

	PuddleUUID uuid.UUID `bun:"puddle_uuid,type:uuid,pk"`
	Puddle     *Puddle   `bun:"rel:belongs-to,join:puddle_uuid=uuid"`

	MemberUUID uuid.UUID  `bun:"member_uuid,type:uuid,pk"`
	AdderUUID  *uuid.UUID `bun:"adder_uuid,type:uuid"`

	JoinedAt time.Time `bun:"joined_at,default:NOW(),notnull"`
}

var _ bun.BeforeAppendModelHook = (*PuddlesMembers)(nil)

func (g *PuddlesMembers) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		g.JoinedAt = time.Now()
	}
	return nil
}
