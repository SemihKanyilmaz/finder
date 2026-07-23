package echo

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(metricsMiddleware())

	e.GET("/health", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
