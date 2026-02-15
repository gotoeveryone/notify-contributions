package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"

	"gotoeveryone/notify-contributions/src/app"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if os.Getenv("DEBUG") != "" {
		err := godotenv.Load()
		if err != nil {
			slog.Error("Error loading .env file", err)
			os.Exit(1)
		}
	}

	if err := app.Notify(time.Now().UTC()); err != nil {
		slog.Error("", err)
		os.Exit(1)
	}
}
