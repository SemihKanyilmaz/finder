package provider

import (
	"context"
	"finder/internal/model"
	"fmt"
	"log/slog"
)

type ContentProvider interface {
	Fetch(ctx context.Context) ([]model.Content, error)
}

type Aggregator struct {
	providers []ContentProvider
}

func NewAggregator(providers ...ContentProvider) *Aggregator {
	return &Aggregator{providers: providers}
}

func (a *Aggregator) FetchAll(ctx context.Context) ([]model.Content, error) {
	type result struct {
		items []model.Content
		err   error
	}

	ch := make(chan result, len(a.providers))

	for _, p := range a.providers {
		go func() {
			items, err := p.Fetch(ctx)
			ch <- result{items, err}
		}()
	}

	var (
		contents []model.Content
		errs     []error
	)

	for range a.providers {
		r := <-ch
		if r.err != nil {
			slog.Error("provider fetch failed", "error", r.err)
			errs = append(errs, r.err)
			continue
		}
		contents = append(contents, r.items...)
	}

	if len(errs) > 0 && len(contents) == 0 {
		return nil, fmt.Errorf("all providers failed: %v", errs)
	}

	if len(errs) > 0 {
		slog.Warn("some providers failed, returning partial results", "failed", len(errs), "succeeded", len(a.providers)-len(errs))
	}

	return contents, nil
}
