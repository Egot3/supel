package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GroupsMembers struct {
	bun.BaseModel `bun:"table:groups_members"`

	GroupUUID uuid.UUID `bun:"group_uuid,type:uuid,pk"`
	Group     *Group    `bun:"rel:belongs-to,join:group_uuid=uuid"`

	MemberUUID uuid.UUID `bun:"member_uuid,type:uuid,pk"`

	JoinedAt time.Time `bun:"joined_at,default:NOW(),notnull"`
}
