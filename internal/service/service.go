package service

import (
	"context"
	"finder/internal/model"
	"finder/internal/provider"
	"finder/internal/repository"
	"finder/internal/scorer"
)

type searchService struct {
	repo       repository.ContentRepository
	aggregator *provider.Aggregator
}

func New(repo repository.ContentRepository, aggregator *provider.Aggregator) *searchService {
	return &searchService{
		repo:       repo,
		aggregator: aggregator,
	}
}

func (s *searchService) Search(ctx context.Context, params model.SearchParams) (model.SearchResult, error) {
	contents, err := s.aggregator.FetchAll(ctx)
	if err != nil {
		return model.SearchResult{}, err
	}

	for i := range contents {
		contents[i].Score = scorer.Score(contents[i])
	}

	if err := s.repo.Upsert(ctx, contents); err != nil {
		return model.SearchResult{}, err
	}

	return s.repo.Search(ctx, params)
}
