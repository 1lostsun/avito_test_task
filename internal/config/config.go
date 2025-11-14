package config

import (
	"avito_test_task/internal/repository/pg"
	"context"
	"log/slog"
)

type Config struct {
	Postgres *pg.Config
}

func New(ctx context.Context) *Config {
	cfg := Config{
		Postgres: &pg.Config{},
	}

	if err := loadSecretsFromVault(ctx, cfg); err != nil {
		slog.Error(err.Error())
	}

	return &cfg
}
