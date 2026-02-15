package app

import (
	"fmt"
	"os"

	"gotoeveryone/notify-contributions/src/domain/entity"
)

type Config struct {
	GitHubToken  string
	GitHubUser   string
	GitLabUserID string
	GitLabToken  string
	NotifyType   string
	SlackWebhook string
	TwitterAuth  entity.TwitterAuth
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		GitHubToken:  os.Getenv("GITHUB_TOKEN"),
		GitHubUser:   os.Getenv("GITHUB_USER_NAME"),
		GitLabUserID: os.Getenv("GITLAB_USER_ID"),
		GitLabToken:  os.Getenv("GITLAB_TOKEN"),
		NotifyType:   os.Getenv("NOTIFY_TYPE"),
		SlackWebhook: os.Getenv("SLACK_WEBHOOK_URL"),
		TwitterAuth: entity.TwitterAuth{
			ConsumerKey:       os.Getenv("TWITTER_COMSUMER_KEY"),
			ConsumerSecret:    os.Getenv("TWITTER_COMSUMER_SECRET"),
			AccessToken:       os.Getenv("TWITTER_ACCESS_TOKEN"),
			AccessTokenSecret: os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"),
		},
	}

	if err := require("GITHUB_TOKEN", cfg.GitHubToken); err != nil {
		return nil, err
	}
	if err := require("GITHUB_USER_NAME", cfg.GitHubUser); err != nil {
		return nil, err
	}
	if err := require("GITLAB_USER_ID", cfg.GitLabUserID); err != nil {
		return nil, err
	}
	if err := require("GITLAB_TOKEN", cfg.GitLabToken); err != nil {
		return nil, err
	}
	if err := require("NOTIFY_TYPE", cfg.NotifyType); err != nil {
		return nil, err
	}

	switch cfg.NotifyType {
	case "slack":
		if err := require("SLACK_WEBHOOK_URL", cfg.SlackWebhook); err != nil {
			return nil, err
		}
	case "twitter":
		if err := require("TWITTER_COMSUMER_KEY", cfg.TwitterAuth.ConsumerKey); err != nil {
			return nil, err
		}
		if err := require("TWITTER_COMSUMER_SECRET", cfg.TwitterAuth.ConsumerSecret); err != nil {
			return nil, err
		}
		if err := require("TWITTER_ACCESS_TOKEN", cfg.TwitterAuth.AccessToken); err != nil {
			return nil, err
		}
		if err := require("TWITTER_ACCESS_TOKEN_SECRET", cfg.TwitterAuth.AccessTokenSecret); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid NOTIFY_TYPE: %s", cfg.NotifyType)
	}

	return cfg, nil
}

func require(key string, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", key)
	}
	return nil
}
