package provider

import (
	"context"
	"errors"
	"finder/internal/model"
	"testing"
)

type mockProvider struct {
	items []model.Content
	err   error
}

func (m *mockProvider) Fetch(_ context.Context) ([]model.Content, error) {
	return m.items, m.err
}

func TestAggregatorAllSuccess(t *testing.T) {
	p1 := &mockProvider{items: []model.Content{{ID: "1"}}}
	p2 := &mockProvider{items: []model.Content{{ID: "2"}}}

	agg := NewAggregator(p1, p2)
	items, err := agg.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("got %d items, want 2", len(items))
	}
}

func TestAggregatorPartialFailure(t *testing.T) {
	p1 := &mockProvider{items: []model.Content{{ID: "1"}}}
	p2 := &mockProvider{err: errors.New("connection refused")}

	agg := NewAggregator(p1, p2)
	items, err := agg.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("partial failure should not return error: %v", err)
	}
	if len(items) != 1 {
		t.Errorf("got %d items, want 1", len(items))
	}
}

func TestAggregatorAllFailed(t *testing.T) {
	p1 := &mockProvider{err: errors.New("timeout")}
	p2 := &mockProvider{err: errors.New("connection refused")}

	agg := NewAggregator(p1, p2)
	_, err := agg.FetchAll(context.Background())
	if err == nil {
		t.Fatal("expected error when all providers fail")
	}
}

func TestAggregatorSingleProvider(t *testing.T) {
	p1 := &mockProvider{items: []model.Content{{ID: "1"}, {ID: "2"}}}

	agg := NewAggregator(p1)
	items, err := agg.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("got %d items, want 2", len(items))
	}
}

func TestAggregatorEmpty(t *testing.T) {
	agg := NewAggregator()
	items, err := agg.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("got %d items, want 0", len(items))
	}
}
