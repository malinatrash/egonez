package app

import (
	"github.com/malinatrash/egonez/config"
	"go.uber.org/zap"
)

func newLogger(config *config.Config) *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}
