package bot

import (
	"context"
	"math/rand"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	if update.Message.BoostAdded != nil {
		h.logger.Info("Boost added")

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "нихуя себе!",
		})
	}

	if update.Message.Sticker != nil {
		h.logger.Info("Sticker message received")

		h.handleStickerMessage(ctx, b, update)
	}

	if strings.Contains(update.Message.Text, "соси") {
		h.logger.Info("Command message received")

		h.sendMessage(ctx, update.Message.Chat.ID, "сам соси")
	}

	if strings.Contains(update.Message.Text, "сосал?") {
		h.logger.Info("Command message received")

		h.sendMessage(ctx, update.Message.Chat.ID, "сосал")
	}

	if update.Message.Video != nil || update.Message.VideoNote != nil || update.Message.Voice != nil {
		h.logger.Info("Video message received")

		chance := rand.Intn(100)
		if chance < 15 {
			h.sendMessage(ctx, update.Message.Chat.ID, "суки я не умею слушать")
		}
	}
}
