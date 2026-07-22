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

func (h *handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/search", h.Search)
}

func (h *handler) Search(c *echo.Context) error {
	var params model.SearchParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	result, err := h.searcher.Search(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
