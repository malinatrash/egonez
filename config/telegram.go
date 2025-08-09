package config

type (
	TelegramConfig struct {
		Token string `envconfig:"BOT_TOKEN" default:""`
	}
)
