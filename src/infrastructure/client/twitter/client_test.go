package twitter

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"

	"gotoeveryone/notify-contributions/src/domain/entity"
)

func TestNotificationExec(t *testing.T) {
	auth := entity.TwitterAuth{
		ConsumerKey:       "consumer_key",
		ConsumerSecret:    "consumer_secret",
		AccessToken:       "access_token",
		AccessTokenSecret: "access_token_secret",
	}
	c := twitterClient{auth: auth}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	url := "https://api.twitter.com/2/tweets"
	httpmock.RegisterResponder("POST", url, httpmock.NewStringResponder(201, `{"data":{"id":"1"}}`))

	if err := c.Exec("test message"); err != nil {
		t.Error(err)
	}

	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("POST %s", url)] != 1 {
		t.Errorf("Failed: POST %s is not called", url)
	}
}
