package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LoggerConfig   LoggerConfig
	TelegramConfig TelegramConfig
	PostgresConfig PostgresConfig
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
