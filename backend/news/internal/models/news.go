package models

import (
	"time"

	"github.com/uptrace/bun"
)

type New struct {
	bun.BaseModel `bun:"table:news"`
	NewUUID       string    `bun:"new_uuid,pk,type:uuid,default:gen_random_uuid()"`
	UserUUID      string    `bun:"user_uuid,notnull"`
	Caption       string    `bun:"caption,notnull"`
	Body          string    `bun:"body"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:now()"`

	_ struct{} `bun:"unique:new_uuid_idx"`
}
