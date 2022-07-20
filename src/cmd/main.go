package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"gotoeveryone/notify-github-contributions/src/domain/entity"
	"gotoeveryone/notify-github-contributions/src/registry"
)

func notify() error {
	userName := os.Getenv("USER_NAME")
	t := entity.TwitterAuth{
		ConsumerKey:       os.Getenv("TWITTER_COMSUMER_KEY"),
		ConsumerSecret:    os.Getenv("TWITTER_COMSUMER_SECRET"),
		AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
	}

	baseDate := time.Now().UTC()

	cc := registry.NewGitHubClient()
	tc, err := registry.NewTwitterClient(t)
	if err != nil {
		if err != nil {
			return err
		}
	}

	u := registry.NewContributionUsecase(cc, tc)
	if err := u.Exec(userName, baseDate); err != nil {
		return err
	}
	return nil
}

func main() {
	if os.Getenv("DEBUG") == "1" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	if err := notify(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
