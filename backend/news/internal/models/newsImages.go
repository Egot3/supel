package models

import (
	"time"

	"github.com/uptrace/bun"
)

type NewsImages struct {
	bun.BaseModel
	ImageUUID string    `bun:"image_uuid,pk,default:gen_random_uuid()"`
	NewUUID   string    `bun:"new_uuid,notnull,rel:belongs-to,join:new_uuid=new_uuid,index"` //о боже
	FileKey   string    `bun:"file_key,notnull"`
	Position  int       `bun:"position,default:0,index"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()"`
}
