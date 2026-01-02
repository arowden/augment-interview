// Package config provides application configuration via environment variables.
package config

import (
	"github.com/arowden/augment-fund/internal/postgres"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration.
type Config struct {
	Database postgres.Config
	Server   Server
}

// Server holds HTTP server configuration.
type Server struct {
	Host string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port int    `envconfig:"SERVER_PORT" default:"8080"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg.Database); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg.Server); err != nil {
		return nil, err
	}

	return &cfg, nil
}
