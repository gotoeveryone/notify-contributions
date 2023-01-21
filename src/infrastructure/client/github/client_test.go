package github

import (
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGitHubGet(t *testing.T) {
	username := "test"
	c := githubClient{username: username}
	httpmock.Activate()

	reqUrl := fmt.Sprintf("https://github.com/users/%s/contributions", username)
	httpmock.RegisterResponder("GET", reqUrl,
		httpmock.NewStringResponder(200, "<rect data-date=\"2006-01-02\">1 contribution on January 2, 2006</rect><rect data-date=\"2006-01-01\">12 contributions on January 1, 2006</rect>"))

	r, err := c.Get(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
		return
	}
	if r.YesterdayCount != 12 {
		t.Errorf("Failed: Count is not matched, actual: [%d], expected: [%d]", r.YesterdayCount, 12)
	}
	if r.BaseDateCount != 1 {
		t.Errorf("Failed: Count is not matched, actual: [%d], expected: [%d]", r.BaseDateCount, 1)
	}
	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("GET %s", reqUrl)] != 1 {
		t.Errorf("Failed: GET %s is not called", reqUrl)
	}
}
