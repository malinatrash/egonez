package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Message struct {
	bun.BaseModel `bun:"table:messages,alias:m"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	ChatID    int64     `bun:"chat_id,notnull" json:"chat_id"`
	UserID    int64     `bun:"user_id,notnull" json:"user_id"`
	Text      string    `bun:"text,notnull" json:"text"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
}
