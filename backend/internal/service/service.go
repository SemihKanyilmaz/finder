package service

import (
	"context"
	"finder/internal/model"
	"finder/internal/provider"
	"finder/internal/repository"
	"finder/internal/scorer"
	"log/slog"
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
		slog.Error("error while fetching all contents", "err", err)
		return model.SearchResult{}, err
	}

	for i := range contents {
		contents[i].Score = scorer.Score(contents[i])
	}

	if err := s.repo.Upsert(ctx, contents); err != nil {
		slog.Error("error while upserting contents", "err", err)
		return model.SearchResult{}, err
	}

	res, err := s.repo.Search(ctx, params)
	if err != nil {
		slog.Error("error while searching contents", "err", err)
	}

	return res, err
}
