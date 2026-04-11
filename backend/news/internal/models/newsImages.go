package models

import (
	"time"

	"github.com/uptrace/bun"
)

type NewsImages struct {
	bun.BaseModel `bun:"table:news_images"`
	ImageUUID     string    `bun:"image_uuid,pk,type:uuid,default:gen_random_uuid()"`
	NewUUID       string    `bun:"new_uuid,type:uuid,notnull"` //о боже
	FileKey       string    `bun:"file_key,notnull"`
	Position      int       `bun:"position,default:0"`
	CreatedAt     time.Time `bun:"created_at,notnull,default:now()"`

	New *New `bun:"rel:belongs-to,join:new_uuid=new_uuid"`
}
