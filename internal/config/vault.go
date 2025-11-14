package config

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault-client-go"
	"log/slog"
)

func loadSecretsFromVault(ctx context.Context, cfg Config) error {
	client, err := vault.New(vault.WithEnvironment())
	if err != nil {
		slog.Error("error creating vault client")
		return fmt.Errorf("error creating vault client: %w", err)
	}

	secret, err := client.Read(ctx, "secret/data/avito-test-task")
	if err != nil {
		slog.Error("error reading secret/data/avito-test-task")
		return fmt.Errorf("error reading secret/data/avito-test-task: %w", err)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		slog.Error("error mapping secret/data/avito-test-task")
		return fmt.Errorf("error mapping secret/data/avito-test-task: %w", err)
	}

	if v, ok := data["PG_USERNAME"].(string); ok && v != "" {
		cfg.Postgres.Username = v
	}

	if v, ok := data["PG_PASSWORD"].(string); ok && v != "" {
		cfg.Postgres.Password = v
	}

	if v, ok := data["PG_DATABASE"].(string); ok && v != "" {
		cfg.Postgres.Database = v
	}

	if v, ok := data["PG_HOST"].(string); ok && v != "" {
		cfg.Postgres.Host = v
	}

	if v, ok := data["PG_PORT"].(string); ok && v != "" {
		cfg.Postgres.Port = v
	}

	return nil
}
