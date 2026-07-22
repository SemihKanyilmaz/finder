package cmd

import (
	"finder/internal/config"
	"finder/pkg/http/echo"
	"log/slog"
)

func Execute() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return
	}

	e := echo.New()

	if err := e.Start(":" + cfg.Port); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}
