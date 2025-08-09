package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

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
