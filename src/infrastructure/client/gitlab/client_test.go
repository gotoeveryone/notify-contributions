package gitlab

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGitHubGet(t *testing.T) {
	userID := "12345678"
	token := "test_token"
	c := gitlabClient{userID, token}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	reqUrlForToday := fmt.Sprintf("https://gitlab.com/api/v4/users/%s/events?after=2005-12-31&before=2006-01-02&per_page=100", userID)
	httpmock.RegisterResponder("GET", reqUrlForToday, func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("PRIVATE-TOKEN") != token {
			t.Errorf("Failed: PRIVATE-TOKEN header is not matched")
		}
		return httpmock.NewStringResponse(200, "[ { \"id\": 1 }, { \"id\": 2 }, { \"id\": 3 } ]"), nil
	})

	reqUrlForYesterday := fmt.Sprintf("https://gitlab.com/api/v4/users/%s/events?after=2006-01-01&before=2006-01-03&per_page=100", userID)
	httpmock.RegisterResponder("GET", reqUrlForYesterday, func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("PRIVATE-TOKEN") != token {
			t.Errorf("Failed: PRIVATE-TOKEN header is not matched")
		}
		return httpmock.NewStringResponse(200, "[ { \"id\": 1 } ]"), nil
	})

	r, err := c.Get(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
		return
	}
	if r.YesterdayCount != 3 {
		t.Errorf("Failed: Name is not matched, actual: [%d], expected: [%d]", r.YesterdayCount, 3)
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
