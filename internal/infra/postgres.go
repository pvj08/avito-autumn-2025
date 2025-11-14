package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/pvj08/avito-autumn-2025/pkg/config"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
)

var ErrDBConnectionFailed = errors.New("db connection failed")

// ConnectDBWithRetry подключается к PostgreSQL с указанным количеством повторных попыток.
// Использует context с таймаутом для ping. Возвращает подключение *sqlx.DB или ошибку.
func ConnectDBWithRetry(ctx context.Context, cfg config.DB, log logger.Logger) (*sqlx.DB, error) {
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
