package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m := map[string]string{
		"text":     message,
		"username": "notify-contributions",
	}
	j, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewBuffer(j))
	if err != nil {
		return fmt.Errorf("failed to build slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call slack webhook: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		b, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			return fmt.Errorf("slack webhook returned %s and reading body failed: %w", res.Status, readErr)
		}
		return fmt.Errorf("slack webhook returned %s: %s", res.Status, strings.TrimSpace(string(b)))
	}

	return nil
}
