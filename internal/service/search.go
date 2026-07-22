package service

import (
	"context"
	"finder/internal/model"
)

type searchService struct {
}

func New() *searchService {
	return &searchService{}
}

func (s *searchService) Search(_ context.Context, _ model.SearchParams) (model.SearchResult, error) {
	return model.SearchResult{}, nil
}
