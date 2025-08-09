package ports

import (
	"context"
	"time"

	"github.com/malinatrash/egonez/internal/entity"
)

type (
	StickerRepository interface {
		Create(ctx context.Context, sticker *entity.Sticker) error
		GetRandom(ctx context.Context, chatID int64) (*entity.Sticker, error)
		CountByChatID(ctx context.Context, chatID int64) (int, error)
		DeleteAll(ctx context.Context, chatID int64) (int64, error)
	}

	MessageRepository interface {
		Create(ctx context.Context, message *entity.Message) error
		GetByChatID(ctx context.Context, chatID int64, limit, offset int) ([]*entity.Message, error)
		CountByChatID(ctx context.Context, chatID int64) (int, error)
		DeleteOlderThan(ctx context.Context, chatID int64, beforeTime time.Time) (int64, error)
		GetRandom(ctx context.Context, chatID int64) (*entity.Message, error)
		GetAllChatIDs(ctx context.Context) ([]int64, error)
	}
)
