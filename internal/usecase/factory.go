package usecase

import (
	"github.com/malinatrash/egonez/config"
	"github.com/malinatrash/egonez/internal/repository"
	"github.com/malinatrash/egonez/internal/usecase/adapters"
	"github.com/malinatrash/egonez/pkg/markov"
	"go.uber.org/zap"
)

type ServiceFactory struct {
	logger     *zap.Logger
	repository *repository.Repository
	config     *config.Config
}

func NewServiceFactory(params Params) *ServiceFactory {
	return &ServiceFactory{
		logger:     params.Logger,
		repository: params.Repo,
		config:     params.Config,
	}
}

func (f *ServiceFactory) NewBotService() adapters.Bot {
	return NewBotService(
		f.repository.MessageRepository,
		f.repository.StickerRepository,
		f.newMarkovService(),
	)
}

func (f *ServiceFactory) newMarkovService() adapters.Markov {
	return markov.NewService(f.config.MarkovConfig.Order, f.repository.MessageRepository, f.logger)
}
