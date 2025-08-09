package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/malinatrash/egonez/config"
	"github.com/malinatrash/egonez/internal/usecase"
	"go.uber.org/zap"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	service *usecase.Service
	logger  *zap.Logger
	bot     *bot.Bot
}

func NewHandler(config *config.Config, service *usecase.Service, logger *zap.Logger) (*Handler, error) {

	h := &Handler{}

	opts := []bot.Option{
		bot.WithDefaultHandler(h.defaultHandler),
	}

	b, err := bot.New(config.TelegramConfig.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	h.service = service
	h.logger = logger
	h.bot = b

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.handleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, h.handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/gen", bot.MatchTypeExact, h.handleGenerate)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/clear", bot.MatchTypeExact, h.handleClear)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/sticker", bot.MatchTypeExact, h.handleSticker)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stats", bot.MatchTypeExact, h.handleStats)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, h.handleTextMessage)

	return h, nil
}

func (h *Handler) Start(ctx context.Context) error {

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	h.bot.Start(ctx)
	return nil
}

func (h *Handler) handleClear(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleClear"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil {
		logger.Error("update.Message is nil")
		return
	}

	chatID := update.Message.Chat.ID
	if err := h.service.BotService.ClearChatHistory(ctx, chatID); err != nil {
		h.sendMessage(ctx, chatID, "❌ Не удалось очистить историю чата. Попробуйте позже.")
		return
	}
	h.sendMessage(ctx, chatID, "🧹 История чата была очищена!")
}
