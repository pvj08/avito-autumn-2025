package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type Config struct {
	Path    string `env:"PATH" env-default:"migrations"`
	Dialect string `env:"DIALECT" env-default:"postgres"`
}

func Up(db *sql.DB, cfg Config) error {
	goose.SetBaseFS(nil) // не используем embed
	if err := goose.SetDialect(cfg.Dialect); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	return goose.Up(db, cfg.Path)
}
