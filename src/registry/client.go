package registry

import (
	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/domain/entity"
	"gotoeveryone/notify-contributions/src/infrastructure/client/github"
	"gotoeveryone/notify-contributions/src/infrastructure/client/gitlab"
	"gotoeveryone/notify-contributions/src/infrastructure/client/slack"
	"gotoeveryone/notify-contributions/src/infrastructure/client/twitter"
)

// NewGitHubClient create client for about contribution use GitHub
func NewGitHubClient(username string, token string) client.Contribution {
	return github.NewClient(username, token)
}

// NewGitlabClient create client for about contribution use Gitlab
func NewGitlabClient(userID string, token string) client.Contribution {
	return gitlab.NewClient(userID, token)
}

// NewTwitterClient is create client for about notification use twitter
func NewTwitterClient(auth entity.TwitterAuth) client.Notification {
	return twitter.NewClient(auth)
}

// NewSlackClient is create client for about notification use slack
func NewSlackClient(webhookURL string) client.Notification {
	return slack.NewClient(webhookURL)
}
