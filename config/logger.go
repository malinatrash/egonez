package config

type (
	LoggerConfig struct {
		Level string `envconfig:"LOG_LEVEL" default:"info"`
		Env   string `envconfig:"ENV" default:"development"`
	}
)
