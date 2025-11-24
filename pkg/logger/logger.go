package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level string `env:"LEVEL"     required:"true"`
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// Init инициализирует глобальный JSON логгер с уровнем из ENV (LOG_LEVEL).
func SetupLogger(c Config) *slog.Logger {
	lvl := parseLogLevel(c.Level)

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	})

	return slog.New(handler)
}

// parseLogLevel преобразует строку в slog.Level
func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
