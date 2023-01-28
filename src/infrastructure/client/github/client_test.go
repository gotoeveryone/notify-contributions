package github

import (
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGitHubGet(t *testing.T) {
	username := "test"
	c := githubClient{username: username}
	httpmock.Activate()

	httpmock.RegisterResponder("POST", "https://api.github.com/graphql",
		httpmock.NewStringResponder(200, `
    {
      "data": {
        "user": {
          "contributionsCollection": {
            "contributionCalendar": {
              "weeks": [
                {
                  "contributionDays": [
                    {
                      "contributionCount": 9,
                      "date": "2006-01-01"
                    }
                  ]
                },
                {
                  "contributionDays": [
                    {
                      "contributionCount": 12,
                      "date": "2006-01-02"
                    }
                  ]
                }
              ]
            }
          }
        }
      }
    }`))

	r, err := c.Get(time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
		return
	}
	if r.YesterdayCount != 9 {
		t.Errorf("Failed: Count is not matched, actual: [%d], expected: [%d]", r.YesterdayCount, 9)
	}
	if r.BaseDateCount != 12 {
		t.Errorf("Failed: Count is not matched, actual: [%d], expected: [%d]", r.BaseDateCount, 12)
	}
}
