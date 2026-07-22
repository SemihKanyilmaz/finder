package echo

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func New() *echo.Echo {

	e := echo.New()

	e.Use(middleware.Recover())

	e.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	return e
}
