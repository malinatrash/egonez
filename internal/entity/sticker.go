package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Sticker struct {
	bun.BaseModel `bun:"table:stickers,alias:s"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	ChatID    int64     `bun:"chat_id,notnull" json:"chat_id"`
	FileID    string    `bun:"file_id,notnull" json:"file_id"`
	SetName   string    `bun:"set_name" json:"set_name"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
}
