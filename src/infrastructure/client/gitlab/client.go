package gitlab

import (
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
		return nil, err
	}
	yc, err := c.response(baseDate.AddDate(0, 0, -1))
	if err != nil {
		return nil, err
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
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return 0, fmt.Errorf("gitlab api returned %s: %s", res.Status, strings.TrimSpace(string(b)))
	}

	r := []any{}
	if err := json.Unmarshal(b, &r); err != nil {
		return 0, err
	}

	return len(r), nil
}
