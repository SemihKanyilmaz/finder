package repository

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

type cachedRepository struct {
	repo  ContentRepository
	cache cache.Cache
	ttl   time.Duration
}

func NewCachedRepository(repo ContentRepository, c cache.Cache, ttl time.Duration) ContentRepository {
	return &cachedRepository{repo: repo, cache: c, ttl: ttl}
}

func (r *cachedRepository) Upsert(ctx context.Context, contents []model.Content) error {
	if err := r.repo.Upsert(ctx, contents); err != nil {
		return err
	}

	if err := r.cache.Invalidate(ctx, "search:*"); err != nil {
		slog.Error("cache invalidate failed after upsert", "error", err)
	}

	return nil
}

func (r *cachedRepository) Search(ctx context.Context, params model.SearchParams) (model.SearchResult, error) {
	key := cacheKey(params)

	data, err := r.cache.Get(ctx, key)
	if err == nil {
		var result model.SearchResult
		if unmarshalErr := json.Unmarshal(data, &result); unmarshalErr == nil {
			metrics.CacheHitsTotal.WithLabelValues("repository", "hit").Inc()
			return result, nil
		} else {
			slog.Error("cache unmarshal failed", "key", key, "err", unmarshalErr)
		}
	} else if !errors.Is(err, redis.Nil) {
		slog.Error("cache get failed", "key", key, "err", err)
	}
	metrics.CacheHitsTotal.WithLabelValues("repository", "miss").Inc()

	result, err := r.repo.Search(ctx, params)
	if err != nil {
		return model.SearchResult{}, err
	}

	if encoded, err := json.Marshal(result); err == nil {
		if err := r.cache.Set(ctx, key, encoded, r.ttl); err != nil {
			slog.Error("cache set failed", "key", key, "error", err)
		}
	}

	return result, nil
}

func cacheKey(p model.SearchParams) string {
	return fmt.Sprintf("search:%s:%s:%s:%d:%d", p.Keyword, p.ContentType, p.SortBy, p.Page, p.PageSize)
}
