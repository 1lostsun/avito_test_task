package app

import (
	"avito_test_task/internal/config"
	"avito_test_task/internal/handler"
	"avito_test_task/internal/repository/pg"
	"avito_test_task/internal/usecase"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Run() error {
	ctx := context.Background()

	cfg := config.New(ctx)
	repo := pg.New(ctx, cfg.Postgres)
	uc := usecase.New(repo)
	h := handler.New(uc)

	r := gin.Default()

	port := fmt.Sprintf(":%s", cfg.AppPort)
	h.InitRoutes(r)

	if err := r.Run(port); err != nil {
		return fmt.Errorf("error running server: %w", err)
	}

	repo.Close()
	return nil
}
