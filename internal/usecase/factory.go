package usecase

import (
	"github.com/malinatrash/egonez/internal/repository"
	"github.com/malinatrash/egonez/pkg/markov"
	"go.uber.org/zap"
)

type ServiceFactory struct {
	logger     *zap.Logger
	repository repository.Repository
	markov     Markov
}

func NewServiceFactory(params Params) *ServiceFactory {
	return &ServiceFactory{
		logger:     params.Logger,
		repository: params.Repo,
		markov:     markov.NewService(3, params.Repo.MessageRepository, params.Logger),
	}
}

func (f *ServiceFactory) NewBotService() Bot {
	return NewBotService(
		f.repository.MessageRepository,
		f.repository.StickerRepository,
		f.markov,
	)
}
