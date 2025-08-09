package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) handleGenerate(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleGenerate"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil {
		logger.Error("update.Message is nil")
		return
	}

	chatID := update.Message.Chat.ID
	msg, err := h.service.BotService.GenerateResponse(ctx, chatID)
	if err != nil {
		h.sendMessage(ctx, chatID, "❌ Failed to generate response. Please try again later.")
		return
	}
	h.sendMessage(ctx, chatID, msg)
}
