package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"gotoeveryone/notify-github-contributions/src/domain/client"
	"gotoeveryone/notify-github-contributions/src/domain/entity"
)

type twitterClient struct {
	auth entity.TwitterAuth
}

func NewClient(auth entity.TwitterAuth) (client.Notification, error) {
	return &twitterClient{
		auth: auth,
	}, nil
}

// Exec is execute notification of summary to target
func (c *twitterClient) Exec(message string) error {
	config := oauth1.NewConfig(c.auth.ConsumerKey, c.auth.ConsumerSecret)
	token := oauth1.NewToken(c.auth.AccessToken, c.auth.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	_, _, err := client.Statuses.Update(message, nil)
	if err != nil {
		return err
	}

	return nil
}
