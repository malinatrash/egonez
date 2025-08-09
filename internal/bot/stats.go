package bot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) handleStats(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleStats"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil {
		logger.Error("update.Message is nil")
		return
	}

	chatID := update.Message.Chat.ID
	stats, err := h.service.BotService.GetChatStats(ctx, chatID)
	if err != nil {
		h.sendMessage(ctx, chatID, "❌ Не удалось получить статистику чата. Попробуйте позже.")
		return
	}

	msg := fmt.Sprintf("📊 Статистика чата\n\n"+
		"Сообщений: %d\n"+
		"Стикеров: %d",
		stats.MessageCount, stats.StickerCount)

	h.sendMessage(ctx, chatID, msg)
}
