package usecase

import (
	"github.com/malinatrash/egonez/config"
	"github.com/malinatrash/egonez/internal/repository"
	"github.com/malinatrash/egonez/internal/usecase/adapters"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Repo   *repository.Repository
	Config *config.Config
}

type Service struct {
	BotService adapters.Bot
}

func NewService(params Params) *Service {
	f := NewServiceFactory(params)

	return &Service{
		BotService: f.NewBotService(),
	}
}
