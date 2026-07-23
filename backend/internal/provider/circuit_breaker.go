package provider

import (
	"context"
	"finder/internal/metrics"
	"finder/internal/model"
	"fmt"
	"time"

	"github.com/sony/gobreaker/v2"
)

type circuitBreakerProvider struct {
	name     string
	cb       *gobreaker.CircuitBreaker[[]model.Content]
	provider ContentProvider
}

func NewCircuitBreakerProvider(name string, p ContentProvider, timeout time.Duration, threshold uint32) ContentProvider {
	st := gobreaker.Settings{
		Name:        name,
		MaxRequests: 1,
		Interval:    0,
		Timeout:     timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= threshold
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			var stateVal float64
			switch to {
			case gobreaker.StateClosed:
				stateVal = 0
			case gobreaker.StateOpen:
				stateVal = 1
			case gobreaker.StateHalfOpen:
				stateVal = 2
			}
			metrics.CircuitBreakerState.WithLabelValues(name).Set(stateVal)
		},
	}

	return &circuitBreakerProvider{
		name:     name,
		cb:       gobreaker.NewCircuitBreaker[[]model.Content](st),
		provider: p,
	}
}

func (c *circuitBreakerProvider) Fetch(ctx context.Context) ([]model.Content, error) {
	contents, err := c.cb.Execute(func() ([]model.Content, error) {
		return c.provider.Fetch(ctx)
	})
	if err != nil {
		if err == gobreaker.ErrOpenState {
			return nil, fmt.Errorf("circuit breaker open for %s: %w", c.name, err)
		}
		return nil, err
	}
	return contents, nil
}
