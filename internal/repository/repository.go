package repository

import (
	"github.com/malinatrash/egonez/internal/ports"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	DB     *bun.DB
	Logger *zap.Logger
}

type Repository struct {
	MessageRepository ports.MessageRepository
	StickerRepository ports.StickerRepository
}

func NewRepository(deps Params) *Repository {
	f := newFactory(deps)

	return &Repository{
		MessageRepository: f.newMessageRepository(),
		StickerRepository: f.newStickerRepository(),
	}
}
