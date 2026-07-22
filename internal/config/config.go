package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port            string `envconfig:"SERVER_PORT" default:"8080"`
	DatabaseURL     string `envconfig:"DATABASE_URL" required:"true"`
	RedisURL        string `envconfig:"REDIS_URL" required:"true"`
	Provider1URL    string `envconfig:"PROVIDER1_URL" required:"true"`
	Provider2URL    string `envconfig:"PROVIDER2_URL" required:"true"`
	RateLimitPerSec  int           `envconfig:"RATE_LIMIT_PER_SEC" default:"10"`
	CacheTTL         time.Duration `envconfig:"CACHE_TTL" default:"5m"`
	ProviderCacheTTL time.Duration `envconfig:"PROVIDER_CACHE_TTL" default:"1m"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
