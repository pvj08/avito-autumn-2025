package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/postgres"
	"github.com/pvj08/avito-autumn-2025/pkg/logger"
	httpserver "github.com/pvj08/avito-autumn-2025/pkg/server"
)

type Config struct {
	Postgres postgres.Config   `env-prefix:"POSTGRES"`
	Server   httpserver.Config `env-prefix:"HTTP"`
	Logger   logger.Config     `env-prefix:"LOG"`
}

func MustLoad() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(fmt.Errorf("config error: %w", err))
	}
	return cfg
}
