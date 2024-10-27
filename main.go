package main

import (
	"log/slog"

	// "github.com/stianeikeland/go-rpio/v4"
	"github.com/victorpigmeo/greenhouse/internal/api"
)

func main() {
	// slog.Info("Setting up GPIO Pins...")
	// rpio.Open()
	// rpio.Pin(24).Output()
	// rpio.Pin(22).Output()

	slog.Info("Starting API...")
	api.Run()
}
