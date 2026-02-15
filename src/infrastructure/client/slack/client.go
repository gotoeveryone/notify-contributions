package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gotoeveryone/notify-contributions/src/domain/client"
)

type slackClient struct {
	webhookURL string
}

func NewClient(webhookURL string) client.Notification {
	return &slackClient{
		webhookURL: webhookURL,
	}
}

func (c *slackClient) Exec(message string) error {
	m := map[string]string{
		"text":     message,
		"username": "notify-contributions",
	}
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := http.Post(c.webhookURL, "application/json", bytes.NewBuffer(j)); err != nil {
		return err
	}

	return nil
}
