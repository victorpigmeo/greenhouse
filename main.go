package main

import (
	"log/slog"

	"github.com/victorpigmeo/greenhouse/internal/api"
)

func main() {
	slog.Info("Starting API...")
	api.Run()
}
