package config

import (
	"fmt"

	"github.com/90poe/envconfig"

	httpv1 "boilerplate/internal/infrastructure/api/http_v1"
	eventstream "boilerplate/internal/infrastructure/event_stream"
	"boilerplate/internal/infrastructure/persistence/postgres"
)

type Config struct {
	HTTP        httpv1.Config
	EventStream eventstream.Config
	Postgres    postgres.Config
	LogLevel    int `envconfig:"LOG_LEVEL" default:"-4"`
}

func FromEnv() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to process config: %w", err)
	}

	return config, nil
}
