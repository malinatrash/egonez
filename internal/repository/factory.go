package repository

import (
	"github.com/malinatrash/egonez/internal/ports"
)

type factory struct {
	deps Params
}

func newFactory(deps Params) *factory {
	return &factory{deps: deps}
}

func (f *factory) newMessageRepository() ports.MessageRepository {
	return NewMessage(f.deps.DB)
}

func (f *factory) newStickerRepository() ports.StickerRepository {
	return NewSticker(f.deps.DB)
}
