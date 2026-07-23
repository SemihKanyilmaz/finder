package provider

import (
	"context"
	"encoding/json"
	"errors"
	"finder/internal/model"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// mockCache is shared across cached_provider_test.go and ratelimiter_test.go
type mockCache struct {
	getData  []byte
	getErr   error
	setErr   error
	incrVal  int64
	incrErr  error
	invalErr error
}

func (m *mockCache) Get(_ context.Context, _ string) ([]byte, error) {
	return m.getData, m.getErr
}

func (m *mockCache) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error {
	return m.setErr
}

func (m *mockCache) Increment(_ context.Context, _ string, _ time.Duration) (int64, error) {
	return m.incrVal, m.incrErr
}

func (m *mockCache) Invalidate(_ context.Context, _ string) error {
	return m.invalErr
}

func TestCachedProvider_CacheHit(t *testing.T) {
	items := []model.Content{{ID: "cached"}}
	data, _ := json.Marshal(items)

	cp := NewCachedProvider("test", &mockProvider{items: []model.Content{{ID: "fresh"}}}, &mockCache{getData: data}, time.Minute)

	got, err := cp.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "cached" {
		t.Errorf("expected cached item, got: %+v", got)
	}
}

func TestCachedProvider_CacheHitInvalidJSON(t *testing.T) {
	cp := NewCachedProvider("test", &mockProvider{items: []model.Content{{ID: "fresh"}}}, &mockCache{getData: []byte("not-json")}, time.Minute)

	got, err := cp.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "fresh" {
		t.Errorf("expected fresh item after invalid cache, got: %+v", got)
	}
}

func TestCachedProvider_CacheMiss(t *testing.T) {
	cp := NewCachedProvider("test", &mockProvider{items: []model.Content{{ID: "fresh"}}}, &mockCache{getErr: redis.Nil}, time.Minute)

	got, err := cp.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "fresh" {
		t.Errorf("expected fresh item, got: %+v", got)
	}
}

func TestCachedProvider_CacheGetError(t *testing.T) {
	cp := NewCachedProvider("test", &mockProvider{items: []model.Content{{ID: "fresh"}}}, &mockCache{getErr: errors.New("connection refused")}, time.Minute)

	got, err := cp.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "fresh" {
		t.Errorf("expected fresh item, got: %+v", got)
	}
}

func TestCachedProvider_ProviderError(t *testing.T) {
	cp := NewCachedProvider("test", &mockProvider{err: errors.New("provider down")}, &mockCache{getErr: redis.Nil}, time.Minute)

	_, err := cp.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error from provider")
	}
}

func TestCachedProvider_EmptyProviderResult(t *testing.T) {
	// Provider boş slice dönünce hata olmadan cache'e yazılmalı ve boş dönmeli
	cp := NewCachedProvider("test", &mockProvider{items: []model.Content{}}, &mockCache{getErr: redis.Nil}, time.Minute)

	got, err := cp.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error on empty result: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %d items, want 0", len(got))
	}
}

func TestCachedProvider_CacheSetError(t *testing.T) {
	cp := NewCachedProvider("test", &mockProvider{items: []model.Content{{ID: "fresh"}}}, &mockCache{getErr: redis.Nil, setErr: errors.New("redis full")}, time.Minute)

	got, err := cp.Fetch(context.Background())
	if err != nil {
		t.Fatalf("cache set error should not propagate: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected items despite set error, got: %+v", got)
	}
}
