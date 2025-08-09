package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) handleSticker(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleSticker"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil {
		logger.Error("update.Message is nil")
		return
	}

	chatID := update.Message.Chat.ID
	sticker, err := h.service.BotService.GetRandomSticker(ctx, chatID)
	if err != nil {
		h.sendMessage(ctx, chatID, "❌ Не удалось получить стикер. Отправьте мне несколько стикеров!")
		return
	}

	h.bot.SendSticker(ctx, &bot.SendStickerParams{
		ChatID: chatID,
		Sticker: &models.InputFileString{
			Data: sticker.FileID,
		},
	})
}

func (h *Handler) handleStickerMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleStickerMessage"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil || update.Message.Sticker == nil {
		logger.Error("update.Message is nil or update.Message.Sticker is nil")
		return
	}

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	sticker := update.Message.Sticker

	err := h.service.BotService.(interface {
		HandleSticker(ctx context.Context, chatID, userID int64, fileID, setName string) error
	}).HandleSticker(ctx, chatID, userID, sticker.FileID, sticker.SetName)

	if err != nil {

		logger.Error("Failed to handle sticker", zap.Error(err))
	}
}
