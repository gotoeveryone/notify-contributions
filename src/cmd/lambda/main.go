package main

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/exp/slog"

	"gotoeveryone/notify-contributions/src/app"
)

func handler(_ context.Context) error {
	return app.Notify(time.Now().UTC())
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	lambda.Start(handler)
}
