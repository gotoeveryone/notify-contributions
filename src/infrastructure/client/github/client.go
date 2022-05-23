package github

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"

	"gotoeveryone/notify-github-contributions/src/domain/client"
	"gotoeveryone/notify-github-contributions/src/domain/entity"
)

type githubClient struct{}

func NewClient() client.Contribution {
	return &githubClient{}
}

// Get is find contribution by identifier
func (c *githubClient) Get(identifier string, baseDate time.Time) (*entity.Contribution, error) {
	res, err := http.Get(fmt.Sprintf("https://github.com/users/%s/contributions", identifier))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	e := entity.Contribution{
		UserName: identifier,
		BaseDate: baseDate,
	}

	yesterday := baseDate.AddDate(0, 0, -1)
	doc.Find(fmt.Sprintf("rect[data-date=\"%s\"]", baseDate.Format("2006-01-02"))).Each(func(i int, s *goquery.Selection) {
		if v, err := strconv.Atoi(s.AttrOr("data-count", "0")); err == nil {
			e.BaseDateCount = v
		}
	})
	doc.Find(fmt.Sprintf("rect[data-date=\"%s\"]", yesterday.Format("2006-01-02"))).Each(func(i int, s *goquery.Selection) {
		if v, err := strconv.Atoi(s.AttrOr("data-count", "0")); err == nil {
			e.YesterdayCount = v
		}
	})

	return &e, nil
}
