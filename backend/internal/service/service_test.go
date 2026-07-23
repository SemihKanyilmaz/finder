package service

import (
	"context"
	"errors"
	"finder/internal/model"
	"finder/internal/provider"
	"testing"
)

type mockContentProvider struct {
	items []model.Content
	err   error
}

func (m *mockContentProvider) Fetch(_ context.Context) ([]model.Content, error) {
	return m.items, m.err
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

func TestService_Search_OK(t *testing.T) {
	p := &mockContentProvider{
		items: []model.Content{{ID: "1", Type: model.ContentTypeVideo, Views: 1000, Likes: 100}},
	}
	repo := &mockRepo{result: model.SearchResult{TotalCount: 1, Page: 1, PageSize: 10}}
	svc := New(repo, provider.NewAggregator(p))

	result, err := svc.Search(context.Background(), model.SearchParams{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalCount != 1 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestService_Search_AllProvidersFail(t *testing.T) {
	p := &mockContentProvider{err: errors.New("timeout")}
	svc := New(&mockRepo{}, provider.NewAggregator(p))

	_, err := svc.Search(context.Background(), model.SearchParams{})
	if err == nil {
		t.Fatal("expected error when all providers fail")
	}
}

func TestService_Search_UpsertError(t *testing.T) {
	p := &mockContentProvider{items: []model.Content{{ID: "1", Type: model.ContentTypeVideo}}}
	svc := New(&mockRepo{upsertErr: errors.New("db error")}, provider.NewAggregator(p))

	_, err := svc.Search(context.Background(), model.SearchParams{})
	if err == nil {
		t.Fatal("expected error on upsert failure")
	}
}

func TestService_Search_SearchError(t *testing.T) {
	p := &mockContentProvider{items: []model.Content{{ID: "1", Type: model.ContentTypeVideo}}}
	svc := New(&mockRepo{searchErr: errors.New("query failed")}, provider.NewAggregator(p))

	_, err := svc.Search(context.Background(), model.SearchParams{})
	if err == nil {
		t.Fatal("expected error on search failure")
	}
}
