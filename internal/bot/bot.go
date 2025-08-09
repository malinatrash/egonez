package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/malinatrash/egonez/config"
	"github.com/malinatrash/egonez/internal/usecase"
	"go.uber.org/zap"

	"math/rand"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	service *usecase.Service
	logger  *zap.Logger
	bot     *bot.Bot
}

func NewHandler(config *config.Config, service *usecase.Service, logger *zap.Logger) (*Handler, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(config.TelegramConfig.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	h := &Handler{
		service: service,
		logger:  logger,
		bot:     b,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.handleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, h.handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/gen", bot.MatchTypeExact, h.handleGenerate)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/clear", bot.MatchTypeExact, h.handleClear)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/sticker", bot.MatchTypeExact, h.handleSticker)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/stats", bot.MatchTypeExact, h.handleStats)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, h.handleTextMessage)

	return h, nil
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

}
func (h *Handler) Start(ctx context.Context) error {

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	h.bot.Start(ctx)
	return nil
}

func (h *Handler) handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleStart"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil {
		logger.Error("update.Message is nil")
		return
	}

	msg := "👋 Привет, я Egonez!\n\n"
	msg += "Я могу учиться на ваших сообщениях и генерировать ответы. Вот что я могу делать:\n"
	msg += "- /gen чтобы сгенерировать ответ\n"
	msg += "- /clear чтобы очистить историю чата\n"
	msg += "- /sticker чтобы получить случайную стикеровую эмодзи\n"
	msg += "- /stats чтобы посмотреть статистику чата"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}

func (h *Handler) handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "bot/handler.handleHelp"

	logger := h.logger.With(zap.String("op", op))

	if update.Message == nil {
		logger.Error("update.Message is nil")
		return
	}

	msg := "🤖 *Ебанез helper*\n\n"
	msg += "*Команды:*\n"
	msg += "/start - Показывает приветственное сообщение\n"
	msg += "/help - Показывает эту помощь\n"
	msg += "/gen - Генерирует ответ на основе учтенных сообщений\n"
	msg += "/clear - Очищает историю чата и сбрасывает обучение\n"
	msg += "/sticker - Получает случайную стикеровую эмодзи\n"
	msg += "/stats - Показывает статистику чата\n\n"
	msg += "Отправьте мне текстовые сообщения и я буду учиться!"

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      msg,
		ParseMode: models.ParseModeMarkdown,
	})
}

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

func (h *Handler) sendMessage(ctx context.Context, chatID int64, text string) {
	h.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
}
