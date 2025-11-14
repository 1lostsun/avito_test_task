package app

import (
	"avito_test_task/internal/config"
	"avito_test_task/internal/repository/pg"
	"context"
	"fmt"
)

func Run() error {
	ctx := context.Background()

	cfg := config.New(ctx)
	pgRepo := pg.New(ctx, cfg.Postgres)

	fmt.Println(pgRepo)
	return nil
}
