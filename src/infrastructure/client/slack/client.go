package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"gotoeveryone/notify-contributions/src/domain/client"
)

type slackClient struct{}

func NewClient() client.Notification {
	return &slackClient{}
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
	if _, err := http.Post(os.Getenv("SLACK_WEBHOOK_URL"), "application/json", bytes.NewBuffer(j)); err != nil {
		return err
	}

	return nil
}
