package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

func NewClient(auth entity.TwitterAuth) client.Notification {
	return &twitterClient{
		auth: auth,
	}
}

// Exec is execute notification of summary to target
func (c *twitterClient) Exec(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := oauth1.NewConfig(c.auth.ConsumerKey, c.auth.ConsumerSecret)
	token := oauth1.NewToken(c.auth.AccessToken, c.auth.AccessTokenSecret)

	httpClient := config.Client(ctx, token)
	httpClient.Timeout = 10 * time.Second

	t := Tweet{Text: message}
	j, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal twitter payload: %w", err)
	}

	body := bytes.NewBuffer(j)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.twitter.com/2/tweets", body)
	if err != nil {
		return fmt.Errorf("failed to build twitter request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call twitter api: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		b, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return fmt.Errorf("twitter api returned %s and reading body failed: %w", res.Status, readErr)
		}
		return fmt.Errorf("twitter api returned %s: %s", res.Status, strings.TrimSpace(string(b)))
	}

	return nil
}
