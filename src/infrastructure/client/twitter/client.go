package twitter

import (
	"bytes"
	"encoding/json"

	"github.com/dghubble/oauth1"

	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/domain/entity"
)

type twitterClient struct {
	auth entity.TwitterAuth
}

type Tweet struct {
	Text string `json:"text"`
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

	client := config.Client(oauth1.NoContext, token)

	t := Tweet{Text: message}
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(j)
	if _, err := client.Post("https://api.twitter.com/2/tweets", "application/json", body); err != nil {
		return err
	}

	return nil
}
