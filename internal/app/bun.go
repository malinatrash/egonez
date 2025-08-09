package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/malinatrash/egonez/config"
	"github.com/malinatrash/egonez/internal/entity"
	"go.uber.org/zap"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var tables = []interface{}{
	(*entity.Message)(nil),
	(*entity.Sticker)(nil),
	(*entity.ChatStats)(nil),
}

func NewDatabase(logger *zap.Logger, cfg *config.Config) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", cfg.PostgresConfig.Host, cfg.PostgresConfig.Port)),
		pgdriver.WithUser(cfg.PostgresConfig.User),
		pgdriver.WithPassword(cfg.PostgresConfig.Password),
		pgdriver.WithDatabase(cfg.PostgresConfig.Name),
		pgdriver.WithTLSConfig(nil),
	))

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(25)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())

	if cfg.LoggerConfig.Env == "development" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.WithEnabled(true),
		))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	for _, table := range tables {
		if _, err := db.NewCreateTable().Model(table).Exec(ctx); err != nil {
			logger.Error("failed to create table", zap.Error(err))
		} else {
			logger.Warn("table %s created", zap.String("table", fmt.Sprintf("%T", table)))
		}
	}

	return db, nil
}
