package slack

import (
	"fmt"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestNotificationExec(t *testing.T) {
	url := "https://hoge.example.com"
	os.Setenv("SLACK_WEBHOOK_URL", url)
	c := slackClient{}
	httpmock.Activate()

	// Exact URL match
	httpmock.RegisterResponder("POST", url,
		httpmock.NewStringResponder(201, `[{"result": "success"}]`))

	if err := c.Exec("test message"); err != nil {
		t.Error(err)
	}
	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("POST %s", url)] != 1 {
		t.Errorf("Failed: POST %s is not called", url)
	}
}
