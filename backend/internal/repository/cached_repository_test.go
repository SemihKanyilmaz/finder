package repository

import (
	"context"
	"encoding/json"
	"errors"
	"finder/internal/model"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type mockCache struct {
	getData     []byte
	getErr      error
	setErr      error
	incrVal     int64
	incrErr     error
	invalErr    error
	invalidated []string
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

func (m *mockCache) Invalidate(_ context.Context, pattern string) error {
	m.invalidated = append(m.invalidated, pattern)
	return m.invalErr
}

type mockRepo struct {
	upsertErr error
	result    model.SearchResult
	searchErr error
}

func (m *mockRepo) Upsert(_ context.Context, _ []model.Content) error {
	return m.upsertErr
}

func (m *mockRepo) Search(_ context.Context, _ model.SearchParams) (model.SearchResult, error) {
	return m.result, m.searchErr
}

func TestCachedRepository_SearchCacheHit(t *testing.T) {
	want := model.SearchResult{TotalCount: 5, Page: 1, PageSize: 10}
	data, _ := json.Marshal(want)

	repo := NewCachedRepository(&mockRepo{}, &mockCache{getData: data}, time.Minute)

	got, err := repo.Search(context.Background(), model.SearchParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TotalCount != 5 {
		t.Errorf("expected cached result, got: %+v", got)
	}
}

func TestCachedRepository_SearchCacheHitInvalidJSON(t *testing.T) {
	want := model.SearchResult{TotalCount: 7}
	repo := NewCachedRepository(&mockRepo{result: want}, &mockCache{getData: []byte("not-json")}, time.Minute)

	got, err := repo.Search(context.Background(), model.SearchParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TotalCount != 7 {
		t.Errorf("expected repo result after invalid cache, got: %+v", got)
	}
}

func TestCachedRepository_SearchCacheMiss(t *testing.T) {
	want := model.SearchResult{TotalCount: 3}
	repo := NewCachedRepository(&mockRepo{result: want}, &mockCache{getErr: redis.Nil}, time.Minute)

	got, err := repo.Search(context.Background(), model.SearchParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TotalCount != 3 {
		t.Errorf("expected repo result, got: %+v", got)
	}
}

func TestCachedRepository_SearchCacheGetError(t *testing.T) {
	want := model.SearchResult{TotalCount: 2}
	repo := NewCachedRepository(&mockRepo{result: want}, &mockCache{getErr: errors.New("redis down")}, time.Minute)

	got, err := repo.Search(context.Background(), model.SearchParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.TotalCount != 2 {
		t.Errorf("expected repo result, got: %+v", got)
	}
}

func TestCachedRepository_SearchRepoError(t *testing.T) {
	repo := NewCachedRepository(&mockRepo{searchErr: errors.New("db down")}, &mockCache{getErr: redis.Nil}, time.Minute)

	_, err := repo.Search(context.Background(), model.SearchParams{})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestCachedRepository_SearchCacheSetError(t *testing.T) {
	want := model.SearchResult{TotalCount: 1}
	repo := NewCachedRepository(&mockRepo{result: want}, &mockCache{getErr: redis.Nil, setErr: errors.New("redis full")}, time.Minute)

	got, err := repo.Search(context.Background(), model.SearchParams{})
	if err != nil {
		t.Fatalf("cache set error should not propagate: %v", err)
	}
	if got.TotalCount != 1 {
		t.Errorf("expected result despite set error: %+v", got)
	}
}

func TestCachedRepository_UpsertSuccess(t *testing.T) {
	c := &mockCache{}
	repo := NewCachedRepository(&mockRepo{}, c, time.Minute)

	if err := repo.Upsert(context.Background(), []model.Content{{ID: "1"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.invalidated) != 1 || c.invalidated[0] != "search:*" {
		t.Errorf("expected search:* invalidated, got: %v", c.invalidated)
	}
}

func TestCachedRepository_UpsertError(t *testing.T) {
	repo := NewCachedRepository(&mockRepo{upsertErr: errors.New("constraint violation")}, &mockCache{}, time.Minute)

	if err := repo.Upsert(context.Background(), []model.Content{{ID: "1"}}); err == nil {
		t.Fatal("expected upsert error")
	}
}

func TestCachedRepository_UpsertInvalidateError(t *testing.T) {
	repo := NewCachedRepository(&mockRepo{}, &mockCache{invalErr: errors.New("scan failed")}, time.Minute)

	// invalidate error is logged but not propagated
	if err := repo.Upsert(context.Background(), []model.Content{{ID: "1"}}); err != nil {
		t.Fatalf("invalidate error should not propagate: %v", err)
	}
}
