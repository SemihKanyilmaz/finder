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

// Search godoc
// @Summary      Search content
// @Description  Fetches content from providers, scores them, and returns paginated results
// @Tags         search
// @Produce      json
// @Param        q        query  string  false  "Keyword"
// @Param        type     query  string  false  "Content type (video, article)"
// @Param        sortBy   query  string  false  "Sort field (score, freshness)"
// @Param        page     query  int     false  "Page number"  default(1)
// @Param        pageSize query  int     false  "Page size"    default(10)
// @Success      200  {object}  model.SearchResult
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /search [get]
type errorResponse struct {
	Message string `json:"message"`
}

func (h *handler) Search(c *echo.Context) error {
	var params model.SearchParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request params"})
	}

	result, err := h.searcher.Search(c.Request().Context(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: "failed to search"})
	}

	return c.JSON(http.StatusOK, result)
}
