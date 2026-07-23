package provider

import (
	"context"
	"errors"
	"finder/internal/model"
	"testing"
)

func TestRateLimitedProvider_WithinLimit(t *testing.T) {
	rl := NewRateLimitedProvider("test", &mockProvider{items: []model.Content{{ID: "1"}}}, &mockCache{incrVal: 1}, 10)

	got, err := rl.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 item, got %d", len(got))
	}
}

func TestRateLimitedProvider_ExceedsLimit(t *testing.T) {
	rl := NewRateLimitedProvider("test", &mockProvider{items: []model.Content{{ID: "1"}}}, &mockCache{incrVal: 11}, 10)

	_, err := rl.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected rate limit error")
	}
}

func TestRateLimitedProvider_IncrementError(t *testing.T) {
	rl := NewRateLimitedProvider("test", &mockProvider{items: []model.Content{{ID: "1"}}}, &mockCache{incrErr: errors.New("redis down")}, 10)

	_, err := rl.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error on increment failure")
	}
}
