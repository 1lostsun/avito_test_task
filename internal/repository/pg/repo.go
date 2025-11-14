package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type Repository struct {
	pg *pgx.Conn
}

func New(ctx context.Context, cfg *Config) *Repository {
	conn, err := pgx.Connect(ctx, cfg.PostgresDSN())
	if err != nil {
		slog.Error(fmt.Sprintf("error connecting to the database: %v", err))
	}

	return &Repository{pg: conn}
}
