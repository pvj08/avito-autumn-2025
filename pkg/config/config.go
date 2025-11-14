package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/pvj08/effective-mobile/pkg/logger"
	"github.com/spf13/viper"
)

type DB struct {
	Host       string        `mapstructure:"host"`
	Port       string        `mapstructure:"port"`
	User       string        `mapstructure:"user"`
	Password   string        `mapstructure:"password"`
	DBName     string        `mapstructure:"name"`
	SSLMode    string        `mapstructure:"ssl_mode"`
	RetryCount int           `mapstructure:"retry_count"`
	RetryDelay time.Duration `mapstructure:"retry_delay"`
}

type Server struct {
	Addr            string        `mapstructure:"addr"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Config struct {
	DB     DB     `mapstructure:"db"`
	Server Server `mapstructure:"server"`
}

// pkg/config/config.go
func LoadConfig(logger logger.Logger) (*Config, error) {
	// .env опционален — если нет файла, просто логируем и продолжаем
	_ = godotenv.Load()

	viper.SetConfigType("yaml")

	// 1) Явный путь из окружения
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		viper.SetConfigFile(p)
	} else {
		viper.AddConfigPath(".")    // текущая рабочая директория
		viper.AddConfigPath("/src") // твой WORKDIR в Docker
		if _, b, _, ok := runtime.Caller(0); ok && filepath.IsAbs(b) {
			basePath := filepath.Join(filepath.Dir(b), "../../")
			viper.AddConfigPath(basePath)
			logger.Info("Using config file from base path", "path", basePath)
		}
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	bindEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Разрешаем отсутствие файла: живём на дефолтах + env
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			logger.Error("Error reading config.yaml", "error", err)
			return nil, fmt.Errorf("ошибка чтения config.yaml: %w", err)
		}
		logger.Info("Config file not found, using only env and defaults")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Error("Error unmarshalling config", "error", err)
		return nil, fmt.Errorf("ошибка декодирования config: %w", err)
	}
	return &cfg, nil
}

func bindEnv() {
	// Секреты
	_ = viper.BindEnv("db.user", "POSTGRES_USER")
	_ = viper.BindEnv("db.password", "POSTGRES_PASSWORD")

	// Если захочешь переопределять не-секреты через env:
	_ = viper.BindEnv("db.host", "POSTGRES_HOST")
	_ = viper.BindEnv("db.port", "POSTGRES_PORT")
	_ = viper.BindEnv("db.name", "POSTGRES_DB")
	// _ = viper.BindEnv("db.ssl_mode", "POSTGRES_SSLMODE")
	// _ = viper.BindEnv("server.addr", "SERVER_ADDR")
}

func (c *DB) DSN() string {
	sslmode := c.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, sslmode,
	)
}
