package adapters

import (
	"context"

	"github.com/malinatrash/egonez/internal/entity"
)

type (
	Bot interface {
		HandleMessage(ctx context.Context, chatID, userID int64, text string) error
		HandleSticker(ctx context.Context, chatID, userID int64, fileID, setName string) error
		GenerateResponse(ctx context.Context, chatID int64) (string, error)
		ClearChatHistory(ctx context.Context, chatID int64) error
		GetRandomSticker(ctx context.Context, chatID int64) (*entity.Sticker, error)
		GetChatStats(ctx context.Context, chatID int64) (*entity.ChatStats, error)
	}

	Markov interface {
		Train(chatID int64, text string) error
		Generate(chatID int64, prefix string, maxLength int) (string, error)
		Clear(chatID int64)
		Load(ctx context.Context, chatID int64) error
	}
)
