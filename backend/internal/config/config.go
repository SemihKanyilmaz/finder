package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ProviderConfig struct {
	Name    string `envconfig:"NAME" required:"true"`
	BaseURL string `envconfig:"BASE_URL" required:"true"`
}

type Config struct {
	Port             string         `envconfig:"SERVER_PORT" default:"8080"`
	DatabaseURL      string         `envconfig:"DATABASE_URL" required:"true"`
	RedisURL         string         `envconfig:"REDIS_URL" required:"true"`
	Provider1        ProviderConfig `envconfig:"PROVIDER1"`
	Provider2        ProviderConfig `envconfig:"PROVIDER2"`
	RateLimitPerSec  int            `envconfig:"RATE_LIMIT_PER_SEC" default:"10"`
	CacheTTL         time.Duration  `envconfig:"CACHE_TTL" default:"5m"`
	ProviderCacheTTL time.Duration  `envconfig:"PROVIDER_CACHE_TTL" default:"1m"`
	CBTimeout        time.Duration  `envconfig:"CB_TIMEOUT" default:"30s"`
	CBThreshold      uint32         `envconfig:"CB_THRESHOLD" default:"3"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
