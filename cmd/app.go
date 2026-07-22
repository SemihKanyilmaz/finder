package cmd

import (
	"finder/pkg/http/echo"
	"log/slog"
)

func Execute() {

	e := echo.New()

	if err := e.Start(":8080"); err != nil {
		slog.Error("failed to start server", "error", err)
	}

}
