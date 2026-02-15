package github

import (
	"context"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"

	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/domain/entity"
)

type githubClient struct {
	username string
	token    string
}

func NewClient(username string, token string) client.Contribution {
	return &githubClient{
		username: username,
		token:    token,
	}
}

// Get is find contribution by identifier
func (c *githubClient) Get(baseDate time.Time) (*entity.Contribution, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var query struct {
		User struct {
			ContrbutionsCollection struct {
				ContributionCalendar struct {
					Weeks []struct {
						ContributionDays []struct {
							ContributionCount githubv4.Int
							Date              githubv4.String
						}
					}
				}
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $login)"`
	}

	yesterday := baseDate.AddDate(0, 0, -1)
	variable := map[string]interface{}{
		"login": githubv4.String(c.username),
		"from":  githubv4.DateTime{Time: yesterday},
		"to":    githubv4.DateTime{Time: baseDate},
	}

	err := client.Query(context.Background(), &query, variable)
	if err != nil {
		return nil, err
	}

	e := entity.Contribution{
		Type:     "GitHub",
		BaseDate: baseDate,
	}

	for _, w := range query.User.ContrbutionsCollection.ContributionCalendar.Weeks {
		for _, d := range w.ContributionDays {
			if string(d.Date) == yesterday.Format("2006-01-02") {
				e.YesterdayCount = int(d.ContributionCount)
			} else if string(d.Date) == baseDate.Format("2006-01-02") {
				e.BaseDateCount = int(d.ContributionCount)
			}
		}
	}

	return &e, nil
}
