package app

import (
	"context"

	"github.com/malinatrash/egonez/internal/bot"
)

func startBot(handler *bot.Handler) error {

	bgCtx := context.Background()

	return handler.Start(bgCtx)
}
