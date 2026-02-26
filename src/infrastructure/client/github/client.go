package github

import (
	"context"
	"fmt"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.token},
	)
	httpClient := oauth2.NewClient(ctx, src)
	ghClient := githubv4.NewClient(httpClient)

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

	if err := ghClient.Query(ctx, &query, variable); err != nil {
		return nil, fmt.Errorf(
			"github contributions query failed (login=%s, from=%s, to=%s): %w",
			c.username,
			yesterday.Format("2006-01-02"),
			baseDate.Format("2006-01-02"),
			err,
		)
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
