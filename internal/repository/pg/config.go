package pg

import (
	"fmt"
	"os"
)

type Config struct {
	Username string `envconfig:"PG_USERNAME"`
	Password string `envconfig:"PG_PASSWORD"`
	Host     string `envconfig:"PG_HOST"`
	Port     string `envconfig:"PG_PORT"`
	Database string `envconfig:"PG_DATABASE"`
}

func (c *Config) PostgresDSN() string {
	fmt.Println(os.Getenv("PG_PORT"))
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Username, c.Password, c.Host, c.Port, c.Database)
}
