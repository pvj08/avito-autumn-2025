package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

type Config struct {
	User     string `env:"USER"     env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Port     string `env:"PORT"     env-required:"true"`
	Host     string `env:"HOST"     env-required:"true"`
	DBName   string `env:"DB_NAME"  env-required:"true"`

	SSLMode    string        `env:"SSL_MODE"  env-default:"disable"`
	RetryCount int           `env:"RETRY_CNT" env-default:"3"`
	RetryDelay time.Duration `env:"RETRY_DUR" env-default:"2s"`
}

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

// New подключается к PostgreSQL с указанным количеством повторных попыток.
// Использует context с таймаутом для ping. Возвращает подключение *sqlx.DB или ошибку.
func New(ctx context.Context, cfg Config, log logger.Logger) (*sqlx.DB, error) {
	dsn := cfg.DSN()
	log.Info("Connecting to database", "dsn", dsn)

	var db *sqlx.DB
	var err error

	for attempt := 1; attempt <= cfg.RetryCount; attempt++ {

		// Подключение к базе данных
		db, err = sqlx.Connect("pgx", dsn)
		if err == nil {
			log.Info("Database connection established", "attempt", attempt)
			return db, nil
		}

		log.Warn("db connect failed", "attempt", attempt, "error", err)

		// Проверяем, не отменён ли внешний контекст
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %v", ErrDBConnectionFailed, ctx.Err())
		default:
		}

		// Задержка между попытками
		time.Sleep(cfg.RetryDelay * time.Duration(attempt))
	}

	return nil, fmt.Errorf("%w: %v", ErrDBConnectionFailed, err)
}

func (c Config) DSN() string {
	sslmode := c.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, sslmode,
	)
}

var ErrDBConnectionFailed = errors.New("db connection failed")
