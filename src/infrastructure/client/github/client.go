package github

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"

	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/domain/entity"
)

type githubClient struct {
	username string
}

func NewClient(username string) client.Contribution {
	return &githubClient{
		username: username,
	}
}

// Get is find contribution by identifier
func (c *githubClient) Get(baseDate time.Time) (*entity.Contribution, error) {
	res, err := http.Get(fmt.Sprintf("https://github.com/users/%s/contributions", c.username))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	e := entity.Contribution{
		Type:     "GitHub",
		BaseDate: baseDate,
	}

	r := regexp.MustCompile(`(\d+) contributions? on`)

	yesterday := baseDate.AddDate(0, 0, -1)
	doc.Find(fmt.Sprintf("rect[data-date=\"%s\"]", baseDate.Format("2006-01-02"))).Each(func(i int, s *goquery.Selection) {
		if r.MatchString(s.Text()) {
			t := r.FindStringSubmatch(s.Text())[1]
			if v, err := strconv.Atoi(t); err == nil {
				e.BaseDateCount = v
			}
		}
	})
	doc.Find(fmt.Sprintf("rect[data-date=\"%s\"]", yesterday.Format("2006-01-02"))).Each(func(i int, s *goquery.Selection) {
		if r.MatchString(s.Text()) {
			t := r.FindStringSubmatch(s.Text())[1]
			if v, err := strconv.Atoi(t); err == nil {
				e.YesterdayCount = v
			}
		}
	})

	return &e, nil
}
