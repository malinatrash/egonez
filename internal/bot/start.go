package bot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

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
