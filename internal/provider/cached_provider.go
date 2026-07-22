package provider

import (
	"context"
	"encoding/json"
	"errors"
	"finder/internal/cache"
	"finder/internal/model"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type cachedProvider struct {
	name     string
	provider ContentProvider
	cache    cache.Cache
	ttl      time.Duration
}

func NewCachedProvider(name string, p ContentProvider, c cache.Cache, ttl time.Duration) ContentProvider {
	return &cachedProvider{name: name, provider: p, cache: c, ttl: ttl}
}

func (p *cachedProvider) Fetch(ctx context.Context) ([]model.Content, error) {
	key := fmt.Sprintf("provider:%s", p.name)

	data, err := p.cache.Get(ctx, key)
	if err == nil {
		var contents []model.Content
		if err := json.Unmarshal(data, &contents); err == nil {
			return contents, nil
		}
		slog.Error("provider cache unmarshal failed", "provider", p.name, "error", err)
	} else if !errors.Is(err, redis.Nil) {
		slog.Error("provider cache get failed", "provider", p.name, "error", err)
	}

	contents, err := p.provider.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	if encoded, err := json.Marshal(contents); err == nil {
		if err := p.cache.Set(ctx, key, encoded, p.ttl); err != nil {
			slog.Error("provider cache set failed", "provider", p.name, "error", err)
		}
	}

	return contents, nil
}
