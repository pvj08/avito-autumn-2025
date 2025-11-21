package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"

	"github.com/pvj08/avito-autumn-2025/internal/infrastructure/postgres"
	httpserver "github.com/pvj08/avito-autumn-2025/pkg/server"
)

type Config struct {
	DB     postgres.Config   `env-prefix:"POSTGRES"`
	Server httpserver.Server `env-prefix:"SERVER"`
}

func MustLoad() (*Config, error) {
	var cfg *Config

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %v\n", err)
	}

	return cfg, nil
}
