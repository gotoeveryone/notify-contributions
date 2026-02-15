package app

import (
	"fmt"
	"time"

	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/registry"
)

// Notify fetches contribution counts and sends a notification.
func Notify(baseDate time.Time) error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	ghc := registry.NewGitHubClient(cfg.GitHubUser, cfg.GitHubToken)
	glc := registry.NewGitlabClient(cfg.GitLabUserID, cfg.GitLabToken)

	var nc client.Notification
	switch cfg.NotifyType {
	case "slack":
		nc = registry.NewSlackClient(cfg.SlackWebhook)
	case "twitter":
		nc = registry.NewTwitterClient(cfg.TwitterAuth)
	default:
		return fmt.Errorf("invalid NOTIFY_TYPE: %s", cfg.NotifyType)
	}

	u := registry.NewContributionUsecase([]client.Contribution{ghc, glc}, nc)
	if err := u.Exec(baseDate); err != nil {
		return err
	}
	return nil
}
