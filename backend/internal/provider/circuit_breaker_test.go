package provider

import (
	"context"
	"errors"
	"finder/internal/model"
	"testing"
	"time"

	"github.com/sony/gobreaker/v2"
)

func TestCircuitBreakerProvider_ClosedState(t *testing.T) {
	p := NewCircuitBreakerProvider("test", &mockProvider{items: []model.Content{{ID: "1"}}}, time.Second, 3)

	got, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error in closed state: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d items, want 1", len(got))
	}
}

func TestCircuitBreakerProvider_OpensAfterThreshold(t *testing.T) {
	const threshold = 3
	mp := &mockProvider{err: errors.New("provider down")}
	p := NewCircuitBreakerProvider("test", mp, time.Minute, threshold)

	// threshold kadar ardışık hata → circuit açılır
	for i := range threshold {
		_, err := p.Fetch(context.Background())
		if err == nil {
			t.Fatalf("call %d: expected error, got nil", i+1)
		}
	}

	// Artık Open state'de olmalı — ErrOpenState sarmalı dönmeli
	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected circuit breaker open error")
	}
	if !errors.Is(err, gobreaker.ErrOpenState) {
		t.Errorf("expected ErrOpenState, got: %v", err)
	}
}

func TestCircuitBreakerProvider_HalfOpenRecovery(t *testing.T) {
	const threshold = 2
	mp := &mockProvider{err: errors.New("provider down")}
	p := NewCircuitBreakerProvider("test", mp, 10*time.Millisecond, threshold)

	// Circuit'i aç
	for range threshold {
		p.Fetch(context.Background()) //nolint:errcheck
	}

	// Timeout bekle → HalfOpen'a geçsin
	time.Sleep(20 * time.Millisecond)

	// Provider'ı düzelt ve başarılı istek yap → Closed'a döner
	mp.err = nil
	mp.items = []model.Content{{ID: "recovered"}}

	got, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("expected recovery after half-open, got: %v", err)
	}
	if len(got) != 1 || got[0].ID != "recovered" {
		t.Errorf("unexpected items after recovery: %+v", got)
	}
}

func TestCircuitBreakerProvider_PartialFailureDoesNotOpen(t *testing.T) {
	const threshold = 3
	failProvider := &mockProvider{err: errors.New("fail")}
	p := NewCircuitBreakerProvider("test", failProvider, time.Minute, threshold)

	// threshold-1 hata → hâlâ Closed olmalı
	for range threshold - 1 {
		p.Fetch(context.Background()) //nolint:errcheck
	}

	// Başarılı istek → consecutive failure sayacı sıfırlanır
	failProvider.err = nil
	failProvider.items = []model.Content{{ID: "ok"}}
	got, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("circuit should still be closed: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d items, want 1", len(got))
	}
}
