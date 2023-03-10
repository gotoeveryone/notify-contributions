package gitlab

import (
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGitHubGet(t *testing.T) {
	userID := "12345678"
	token := "test_token"
	c := gitlabClient{userID, token}
	httpmock.Activate()

	reqUrlForToday := fmt.Sprintf("https://gitlab.com/api/v4/users/%s/events?private_token=%s&before=2006-01-02&after=2005-12-31&per_page=100", userID, token)
	httpmock.RegisterResponder("GET", reqUrlForToday,
		httpmock.NewStringResponder(200, "[ { \"id\": 1 }, { \"id\": 2 }, { \"id\": 3 } ]"))

	reqUrlForYesterday := fmt.Sprintf("https://gitlab.com/api/v4/users/%s/events?private_token=%s&before=2006-01-03&after=2006-01-01&per_page=100", userID, token)
	httpmock.RegisterResponder("GET", reqUrlForYesterday,
		httpmock.NewStringResponder(200, "[ { \"id\": 1 } ]"))

	r, err := c.Get(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
		return
	}
	if r.YesterdayCount != 3 {
		t.Errorf("Failed: Name is not matched, actual: [%d], expected: [%d]", r.YesterdayCount, 2)
	}
	if r.BaseDateCount != 1 {
		t.Errorf("Failed: Name is not matched, actual: [%d], expected: [%d]", r.BaseDateCount, 1)
	}
	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("GET %s", reqUrlForToday)] != 1 {
		t.Errorf("Failed: GET %s is not called", reqUrlForToday)
	}
	if info[fmt.Sprintf("GET %s", reqUrlForYesterday)] != 1 {
		t.Errorf("Failed: GET %s is not called", reqUrlForYesterday)
	}
}
