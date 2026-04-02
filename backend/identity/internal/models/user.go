package models

import (
	"github.com/Egot3/supel/backend/identity/internal/types"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"` //во имя святого ORM
	UUID          string              `bun:"uuid,pk,type:uuid,default:gen_random_uuid()"`
	Role          types.UserRole      `bun:"role,notnull,default:'USER'"`
	PasswordHash  string              `bun:"password_hash,notnull"`
	Email         string              `bun:"email,notnull"`
	IsActive      bool                `bun:"is_active,default:true"`

	_ struct{} `bun:"unique:uuid_idx,type:unique"`
}
