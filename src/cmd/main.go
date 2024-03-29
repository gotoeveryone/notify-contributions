package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"

	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/registry"
)

func notify() error {
	baseDate := time.Now().UTC()

	ghc := registry.NewGitHubClient(os.Getenv("GITHUB_USER_NAME"))
	glc := registry.NewGitlabClient(os.Getenv("GITLAB_USER_ID"), os.Getenv("GITLAB_TOKEN"))
	// t := entity.TwitterAuth{
	// 	ConsumerKey:       os.Getenv("TWITTER_COMSUMER_KEY"),
	// 	ConsumerSecret:    os.Getenv("TWITTER_COMSUMER_SECRET"),
	// 	AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
	// 	AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	// }
	// nc, err := registry.NewTwitterClient(t)
	// if err != nil {
	// 	return err
	// }
	nc := registry.NewSlackClient()

	u := registry.NewContributionUsecase([]client.Contribution{ghc, glc}, nc)
	if err := u.Exec(baseDate); err != nil {
		return err
	}
	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if os.Getenv("DEBUG") != "" {
		err := godotenv.Load()
		if err != nil {
			slog.Error("Error loading .env file", err)
			os.Exit(1)
		}
	}

	if err := notify(); err != nil {
		slog.Error("", err)
		os.Exit(1)
	}
}
