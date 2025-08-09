package repository

import (
	"context"

	"github.com/go-telegram/bot/models"
	"github.com/malinatrash/egonez/internal/entity"
	"github.com/malinatrash/egonez/internal/ports"
	"github.com/uptrace/bun"
)

var _ ports.StickerRepository = (*Sticker)(nil)

type Sticker struct {
	db *bun.DB
}

func NewSticker(db *bun.DB) ports.StickerRepository {
	return &Sticker{db: db}
}

func (r *Sticker) Create(ctx context.Context, sticker *entity.Sticker) error {
	_, err := r.db.NewInsert().
		Model(sticker).
		On("CONFLICT (file_id) DO NOTHING").
		Exec(ctx)
	return err
}

func (r *Sticker) GetRandom(ctx context.Context, chatID int64) (*entity.Sticker, error) {
	var sticker entity.Sticker

	err := r.db.NewSelect().
		Model(&sticker).
		Where("chat_id = ?", chatID).
		Order("RANDOM()").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return &sticker, nil
}

func (r *Sticker) CountByChatID(ctx context.Context, chatID int64) (int, error) {
	return r.db.NewSelect().
		Model((*models.Sticker)(nil)).
		Where("chat_id = ?", chatID).
		Count(ctx)
}

func (r *Sticker) DeleteAll(ctx context.Context, chatID int64) (int64, error) {
	res, err := r.db.NewDelete().
		Model((*models.Sticker)(nil)).
		Where("chat_id = ?", chatID).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
