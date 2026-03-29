package models

import (
	"github.com/Egot3/supel/backend/identity/types"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`  //во имя святого ORM
	UUID            string               `bun:"uuid,pk,type:uuid,default:gen_random_uuid()"`
	Role          	types.UserRole       `bun:"role,notnull"`
	IsActive      	bool                 `bun:"isActive,default:true"`

	_ struct{} `bun:"unique:uuid_idx,type:unique"`
}