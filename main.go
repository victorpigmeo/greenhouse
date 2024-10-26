package main

import (
	"fmt"
	"log/slog"

	"github.com/stianeikeland/go-rpio/v4"
	"github.com/victorpigmeo/greenhouse/internal/api"
)

func main() {
	slog.Info("Setting up GPIO Pins...")
	rpio.Open()
	rpio.Pin(24).Output()

	fmt.Println("Starting API...")
	api.Run()
}
