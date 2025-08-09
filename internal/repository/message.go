package repository

import (
	"context"
	"time"

	"github.com/malinatrash/egonez/internal/entity"
	"github.com/malinatrash/egonez/internal/ports"
	"github.com/uptrace/bun"
)

var _ ports.MessageRepository = (*Message)(nil)

type Message struct {
	db *bun.DB
}

func NewMessage(db *bun.DB) *Message {
	return &Message{db: db}
}

func (r *Message) Create(ctx context.Context, message *entity.Message) error {
	_, err := r.db.NewInsert().
		Model(message).
		Exec(ctx)
	return err
}

func (r *Message) GetByChatID(ctx context.Context, chatID int64, limit, offset int) ([]*entity.Message, error) {
	var messages []*entity.Message
	err := r.db.NewSelect().
		Model(&messages).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	return messages, err
}

func (r *Message) CountByChatID(ctx context.Context, chatID int64) (int, error) {
	return r.db.NewSelect().
		Model((*entity.Message)(nil)).
		Where("chat_id = ?", chatID).
		Count(ctx)
}

func (r *Message) DeleteOlderThan(ctx context.Context, chatID int64, beforeTime time.Time) (int64, error) {
	res, err := r.db.NewDelete().
		Model((*entity.Message)(nil)).
		Where("chat_id = ? AND created_at < ?", chatID, beforeTime).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (r *Message) GetRandom(ctx context.Context, chatID int64) (*entity.Message, error) {
	var message entity.Message
	err := r.db.NewSelect().
		Model(&message).
		Where("chat_id = ?", chatID).
		Order("RANDOM()").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *Message) GetAllChatIDs(ctx context.Context) ([]int64, error) {
	var chatIDs []int64
	
	err := r.db.NewSelect().
		Model((*entity.Message)(nil)).
		ColumnExpr("DISTINCT chat_id").
		Order("chat_id").
		Scan(ctx, &chatIDs)

	if err != nil {
		return nil, err
	}

	return chatIDs, nil
}
