package handler

import (
	"context"
	"encoding/json"
	"errors"
	"finder/internal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

type mockSearcher struct {
	result model.SearchResult
	err    error
}

func (m *mockSearcher) Search(_ context.Context, _ model.SearchParams) (model.SearchResult, error) {
	return m.result, m.err
}

func TestRegisterRoutes(t *testing.T) {
	h := New(&mockSearcher{})
	e := echo.New()
	h.RegisterRoutes(e) // should not panic
}

func TestSearch_OK(t *testing.T) {
	want := model.SearchResult{
		Items:      []model.Content{{ID: "1", Title: "Go"}},
		TotalCount: 1,
		Page:       1,
		PageSize:   10,
	}
	h := New(&mockSearcher{result: want})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/search?q=go&page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Search(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var got model.SearchResult
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.TotalCount != 1 || len(got.Items) != 1 {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestSearch_ServiceError(t *testing.T) {
	h := New(&mockSearcher{err: errors.New("db error")})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Search(c); err != nil {
		t.Fatalf("handler should write response, not return error: %v", err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestSearch_BindError(t *testing.T) {
	h := New(&mockSearcher{})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/search?page=not-a-number", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := h.Search(c); err != nil {
		t.Fatalf("handler should write response, not return error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
