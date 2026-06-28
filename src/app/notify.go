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

	ghc, err := newGitHubClient(cfg)
	if err != nil {
		return err
	}

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

func newGitHubClient(cfg *Config) (client.Contribution, error) {
	if !cfg.UseGitHubApp() {
		return registry.NewGitHubClient(cfg.GitHubUser, cfg.GitHubToken), nil
	}

	ghc, err := registry.NewGitHubAppClient(
		cfg.GitHubUser,
		cfg.GitHubAppID,
		cfg.GitHubAppInstallationID,
		cfg.GitHubAppPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	return ghc, nil
}
