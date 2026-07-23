package echo

import (
	"strconv"
	"time"

	"finder/internal/metrics"

	"github.com/labstack/echo/v5"
)

func metricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			err := next(c)

			_, status := echo.ResolveResponseStatus(c.Response(), err)

			metrics.HTTPRequestsTotal.WithLabelValues(
				c.Request().Method,
				c.Path(),
				strconv.Itoa(status),
			).Inc()

			metrics.HTTPRequestDuration.WithLabelValues(
				c.Request().Method,
				c.Path(),
			).Observe(time.Since(start).Seconds())

			return err
		}
	}
}
