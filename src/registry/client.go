package registry

import (
	"gotoeveryone/notify-github-contributions/src/domain/client"
	"gotoeveryone/notify-github-contributions/src/domain/entity"
	"gotoeveryone/notify-github-contributions/src/infrastructure/client/github"
	"gotoeveryone/notify-github-contributions/src/infrastructure/client/twitter"
)

// NewGitHubClient create client for about contribution use github
func NewGitHubClient() client.Contribution {
	return github.NewClient()
}

// NewTwitterClient is create client for about notification use twitter
func NewTwitterClient(auth entity.TwitterAuth) (client.Notification, error) {
	return twitter.NewClient(auth)
}
