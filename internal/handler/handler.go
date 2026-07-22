package handler

import (
	"context"
	"finder/internal/model"
	"net/http"

	"github.com/labstack/echo/v5"
)

type searcher interface {
	Search(ctx context.Context, params model.SearchParams) (model.SearchResult, error)
}

type handler struct {
	searcher searcher
}

func New(s searcher) *handler {
	return &handler{searcher: s}
}

func (h *handler) Search(c *echo.Context) error {

	return c.JSON(http.StatusOK, map[string]any{
		"data": "data",
	})
}
