package pg

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log/slog"
)

type Config struct {
	Username string `envconfig:"PG_USERNAME"`
	Password string `envconfig:"PG_PASSWORD"`
	Host     string `envconfig:"PG_HOST"`
	Port     string `envconfig:"PG_PORT"`
	Database string `envconfig:"PG_DATABASE"`
}

func (c *Config) PostgresDSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.Database)
	slog.Info("connecting to DB", slog.String("dsn", dsn))
	return dsn
}

func RunMigrations(dsn string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		return
	}

	if err := m.Up(); err != nil && !errors.Is(migrate.ErrNoChange, err) {
		slog.Error(fmt.Sprintf("error running migrations: %v", err))
		return
	}

	slog.Info("Migrations applied successfully")
}
