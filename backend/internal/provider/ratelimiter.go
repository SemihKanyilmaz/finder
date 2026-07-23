package provider

import (
	"context"
	"finder/internal/cache"
	"finder/internal/model"
	"fmt"
	"log/slog"
	"time"
)

type rateLimitedProvider struct {
	name     string
	provider ContentProvider
	cache    cache.Cache
	limit    int
	window   time.Duration
}

func NewRateLimitedProvider(name string, p ContentProvider, c cache.Cache, limitPerSec int) ContentProvider {
	return &rateLimitedProvider{
		name:     name,
		provider: p,
		cache:    c,
		limit:    limitPerSec,
		window:   time.Second,
	}
}

func (r *rateLimitedProvider) Fetch(ctx context.Context) ([]model.Content, error) {
	key := fmt.Sprintf("ratelimit:%s", r.name)

	count, err := r.cache.Increment(ctx, key, r.window)
	if err != nil {
		return nil, fmt.Errorf("rate limit check: %w", err)
	}

	if count > int64(r.limit) {
		slog.Warn("rate limit exceeded", "provider", r.name, "count", count, "limit", r.limit)
		return nil, fmt.Errorf("rate limit exceeded for %s", r.name)
	}

	return r.provider.Fetch(ctx)
}
