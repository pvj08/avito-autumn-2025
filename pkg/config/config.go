package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/postgres"
	httpserver "github.com/pvj08/avito-autumn-2025/pkg/server"
)

type Config struct {
	DB     postgres.Config   `env-prefix:"POSTGRES"`
	Server httpserver.Server `env-prefix:"HTTP"`
}

func MustLoad() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(fmt.Errorf("config error: %w", err))
	}
	return cfg
}
