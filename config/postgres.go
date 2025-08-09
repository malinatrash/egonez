package config

type (
	PostgresConfig struct {
		Host     string `envconfig:"DB_HOST" default:"postgres"`
		Port     int    `envconfig:"DB_PORT" default:"5432"`
		User     string `envconfig:"DB_USER" default:"postgres"`
		Password string `envconfig:"DB_PASSWORD" default:"postgres"`
		Name     string `envconfig:"DB_NAME" default:"egonez"`
		SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	}
)
