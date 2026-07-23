package provider

import (
	"context"
	"encoding/json"
	"errors"
	"finder/internal/cache"
	"finder/internal/metrics"
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
		if unmarshalErr := json.Unmarshal(data, &contents); unmarshalErr == nil {
			metrics.CacheHitsTotal.WithLabelValues("provider", "hit").Inc()
			return contents, nil
		} else {
			slog.Error("provider cache unmarshal failed", "provider", p.name, "err", unmarshalErr)
		}
	} else if !errors.Is(err, redis.Nil) {
		slog.Error("provider cache get failed", "provider", p.name, "err", err)
	}
	metrics.CacheHitsTotal.WithLabelValues("provider", "miss").Inc()

	contents, err := p.provider.Fetch(ctx)
	if err != nil {
		metrics.ProviderFetchTotal.WithLabelValues(p.name, "failure").Inc()
		slog.Error("error while fetching data from provider.", " provider: ", p.name)
		return nil, err
	}
	metrics.ProviderFetchTotal.WithLabelValues(p.name, "success").Inc()

	if encoded, err := json.Marshal(contents); err == nil {
		if err := p.cache.Set(ctx, key, encoded, p.ttl); err != nil {
			slog.Error("provider cache set failed", "provider", p.name, "error", err)
		}
	}

	return contents, nil
}
