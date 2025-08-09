package usecase

import (
	"github.com/malinatrash/egonez/config"
	"github.com/malinatrash/egonez/internal/repository"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Repo   repository.Repository
	Config *config.Config
}

type Service struct {
	BotService Bot
}

func NewService(params Params) *Service {
	f := NewServiceFactory(params)

	return &Service{
		BotService: f.NewBotService(),
	}
}
