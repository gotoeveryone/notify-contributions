package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/domain/entity"
)

type gitlabClient struct {
	userID string
	token  string
}

func NewClient(userID string, token string) client.Contribution {
	return &gitlabClient{
		userID: userID,
		token:  token,
	}
}

// Get is find contribution by identifier
func (c *gitlabClient) Get(baseDate time.Time) (*entity.Contribution, error) {
	tc, err := c.response(baseDate)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gitlab contributions for base date: %w", err)
	}
	yc, err := c.response(baseDate.AddDate(0, 0, -1))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gitlab contributions for yesterday: %w", err)
	}

	e := entity.Contribution{
		Type:           "GitLab",
		BaseDate:       baseDate,
		BaseDateCount:  tc,
		YesterdayCount: yc,
	}

	return &e, nil
}

func (c *gitlabClient) response(baseDate time.Time) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 前日～翌日を指定することで当日分を取得できる
	before := baseDate.AddDate(0, 0, 1).Format("2006-01-02")
	after := baseDate.AddDate(0, 0, -1).Format("2006-01-02")
	// 上限はひとまず100とする
	perPage := 100

	query := url.Values{}
	query.Set("before", before)
	query.Set("after", after)
	query.Set("per_page", fmt.Sprintf("%d", perPage))

	endpoint := fmt.Sprintf("https://gitlab.com/api/v4/users/%s/events?%s", url.PathEscape(c.userID), query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to build gitlab request (date=%s): %w", baseDate.Format("2006-01-02"), err)
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	res, err := httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call gitlab api (date=%s): %w", baseDate.Format("2006-01-02"), err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read gitlab response body (date=%s): %w", baseDate.Format("2006-01-02"), err)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return 0, fmt.Errorf(
			"gitlab api returned %s (date=%s): %s",
			res.Status,
			baseDate.Format("2006-01-02"),
			strings.TrimSpace(string(b)),
		)
	}

	r := []any{}
	if err := json.Unmarshal(b, &r); err != nil {
		return 0, fmt.Errorf("failed to decode gitlab response (date=%s): %w", baseDate.Format("2006-01-02"), err)
	}

	return len(r), nil
}
