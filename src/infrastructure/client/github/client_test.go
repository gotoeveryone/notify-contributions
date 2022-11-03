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
		httpmock.NewStringResponder(200, "<rect data-date=\"2006-01-02\" data-count=\"1\" /><rect data-date=\"2006-01-01\" data-count=\"2\" />"))

	r, err := c.Get(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
		return
	}
	if r.YesterdayCount != 2 {
		t.Errorf("Failed: Name is not matched, actual: [%d], expected: [%d]", r.YesterdayCount, 2)
	}
	if r.BaseDateCount != 1 {
		t.Errorf("Failed: Name is not matched, actual: [%d], expected: [%d]", r.BaseDateCount, 1)
	}
	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("GET %s", reqUrl)] != 1 {
		t.Errorf("Failed: GET %s is not called", reqUrl)
	}
}
