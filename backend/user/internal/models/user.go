package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	UUID          string  `bun:"uuid,pk,type:uuid"`
	Nickname      string  `bun:"nickname,notnull"`
	Description   *string `bun:"description"`
	AvatarKey     *string `bun:"avatar_key"`

	Status            *string    `bun:"status,null"`
	StatusExpiration  *time.Time `bun:"status_expiration,null"`
	StatusReactionKey *string    `bun:"status_reaction_key,null"`

	CreatedAt time.Time  `bun:"created_at,notnull,default:now()"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`

	_ struct{} `bun:"unique:new_uuid_idx"`
}
