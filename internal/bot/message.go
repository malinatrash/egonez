package bot

import (
	"context"
	"math/rand"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) sendMessage(ctx context.Context, chatID int64, text string) {
	h.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
}

func (h *Handler) handleTextMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleTextMessage"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil || update.Message.Text == "" {
		logger.Error("update.Message is nil or update.Message.Text is empty")
		return
	}

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)

	if strings.HasPrefix(text, "/") {
		return
	}

	if err := h.service.BotService.HandleMessage(ctx, chatID, userID, text); err != nil {
		logger.Error("Failed to handle message", zap.Error(err))
	}

	if rand.Intn(100) > 70 {
		h.handleGenerate(ctx, b, update)
	}
}
