package app

import (
	"github.com/malinatrash/egonez/config"
	"github.com/malinatrash/egonez/internal/bot"
	"github.com/malinatrash/egonez/internal/repository"
	"github.com/malinatrash/egonez/internal/usecase"
	"go.uber.org/fx"
)

func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.Load,
			newLogger,
			NewDatabase,
			repository.NewRepository,
			usecase.NewService,
			bot.NewHandler,
		),
		fx.Invoke(
			startBot,
		),
	)
}
