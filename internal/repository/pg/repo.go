package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type Repository struct {
	pg *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) *Repository {
	dsn := cfg.PostgresDSN()
	RunMigrations(dsn)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("failed to parse config", "err", err)
		return nil
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("error connecting to the database", "error", err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		slog.Error("error pinging database", "error", err)
		return nil
	}

	slog.Info("Database connected successfully")
	return &Repository{pg: pool}
}

func (r *Repository) Close() {
	r.pg.Close()
}
