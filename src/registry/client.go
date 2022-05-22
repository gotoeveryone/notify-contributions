package registry

import (
	"gotoeveryone/notify-github-contributions/src/domain/client"
	"gotoeveryone/notify-github-contributions/src/infrastructure/client/github"
	"gotoeveryone/notify-github-contributions/src/infrastructure/client/twitter"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// NewGitHubClient create client for about contribution use github
func NewGitHubClient() client.Contribution {
	return github.NewClient()
}

// NewTwitterClient is create client for about notification use twitter
func NewTwitterClient(c ssm.Client) (client.Notification, error) {
	return twitter.NewClient(c)
}
