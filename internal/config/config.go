package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port         string `envconfig:"SERVER_PORT" default:"8080"`
	Provider1URL string `envconfig:"PROVIDER1_URL" required:"true"`
	Provider2URL string `envconfig:"PROVIDER2_URL" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
