package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
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

// Load загружает конфиг из файла + ENV.
// Требует НЕ-nil логгер
func Load(log logger.Logger) (*Config, error) {
	if log == nil {
		return nil, errors.New("logger is required for config loading")
	}

	// .env опционален
	_ = godotenv.Load()

	viper.SetConfigType("yaml")

	// 1) Явный путь из CONFIG_PATH
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		viper.SetConfigFile(p)
	} else {
		// 2) По умолчанию — config.yaml из текущей директории
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	// ENV имеют приоритет над файлом
	viper.AutomaticEnv()
	bindEnv()
	setDefaults()

	// Пытаемся прочитать файл конфига
	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			log.Error("Error reading config file", "error", err)
			return nil, fmt.Errorf("ошибка чтения config.yaml: %w", err)
		}
		log.Info("Config file not found, using only env and defaults")
	} else {
		if path := viper.ConfigFileUsed(); path != "" {
			log.Info("Loaded config file", "path", path)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Error("Error unmarshalling config", "error", err)
		return nil, fmt.Errorf("ошибка декодирования config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		log.Error("Invalid config", "error", err)
		return nil, err
	}

	return &cfg, nil
}

// Привязка ключей к переменным окружения.
func bindEnv() {
	// Секреты
	_ = viper.BindEnv("db.user", "POSTGRES_USER")
	_ = viper.BindEnv("db.password", "POSTGRES_PASSWORD")

	// Не-секреты
	_ = viper.BindEnv("db.host", "POSTGRES_HOST")
	_ = viper.BindEnv("db.port", "POSTGRES_PORT")
	_ = viper.BindEnv("db.name", "POSTGRES_DB")
}

// Дефолты, если нет ни файла, ни ENV.
func setDefaults() {
	// DB
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "5432")
	viper.SetDefault("db.ssl_mode", "disable")
	viper.SetDefault("db.retry_count", 3)
	viper.SetDefault("db.retry_delay", time.Second)

	// Server
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.read_timeout", 5*time.Second)
	viper.SetDefault("server.write_timeout", 5*time.Second)
	viper.SetDefault("server.idle_timeout", 60*time.Second)
	viper.SetDefault("server.shutdown_timeout", 10*time.Second)
}

// Минимальная валидация, чтобы не получить пустой DSN.
func (c *Config) validate() error {
	if c.DB.Host == "" || c.DB.Port == "" || c.DB.User == "" || c.DB.DBName == "" {
		return fmt.Errorf("неполный конфиг БД: host/port/user/name должны быть заданы")
	}
	return nil
}

// DSN собирает строку подключения к БД.
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
