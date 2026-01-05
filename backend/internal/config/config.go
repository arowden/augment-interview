package config

import (
	"github.com/arowden/augment-fund/internal/otel"
	"github.com/arowden/augment-fund/internal/postgres"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database  postgres.Config
	Server    Server
	Telemetry otel.Config
}

type Server struct {
	Host        string   `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port        int      `envconfig:"SERVER_PORT" default:"8080"`
	CORSOrigins []string `envconfig:"CORS_ORIGINS" default:"http://localhost:*,http://127.0.0.1:*"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg.Database); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg.Server); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg.Telemetry); err != nil {
		return nil, err
	}

	return &cfg, nil
}
