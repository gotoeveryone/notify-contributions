package twitter

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"gotoeveryone/notify-github-contributions/src/domain/client"
	"gotoeveryone/notify-github-contributions/src/domain/entity"
	"gotoeveryone/notify-github-contributions/src/helper"
)

type twitterClient struct {
	auth entity.Twitter
}

func NewClient(c ssm.Client) (client.Notification, error) {
	consumerKey, err := helper.GetParameter(c, "twitter_consumer_key")
	if err != nil {
		return nil, err
	}

	consumerSecret, err := helper.GetParameter(c, "twitter_consumer_secret")
	if err != nil {
		return nil, err
	}

	accessToken, err := helper.GetParameter(c, "twitter_access_token")
	if err != nil {
		return nil, err
	}

	accessTokenSecret, err := helper.GetParameter(c, "twitter_access_token_secret")
	if err != nil {
		return nil, err
	}

	return &twitterClient{
		auth: entity.Twitter{
			ConsumerKey:       *consumerKey,
			ConsumerSecret:    *consumerSecret,
			AccessToken:       *accessToken,
			AccessTokenSecret: *accessTokenSecret,
		},
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
