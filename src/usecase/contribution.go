package usecase

import (
	"errors"
	"gotoeveryone/notify-contributions/src/domain/client"
	"log"
	"os"
	"strings"
	"time"
)

var (
	ErrContributionNotFound = errors.New("contribution is nil")
)

type Contribution interface {
	Exec(baseDate time.Time) error
}

type contributionUsecase struct {
	sources      []client.Contribution
	notification client.Notification
}

func NewContributionUsecase(sources []client.Contribution, notification client.Notification) Contribution {
	return &contributionUsecase{
		sources:      sources,
		notification: notification,
	}
}

// Exec is get contribution and notify
func (u *contributionUsecase) Exec(baseDate time.Time) error {
	messages := []string{}
	for _, source := range u.sources {
		c, err := source.Get(baseDate)
		if err != nil {
			return err
		}
		if c == nil {
			return ErrContributionNotFound
		}
		messages = append(messages, c.Message())
	}
	message := strings.Join(messages, "\n\n")
	if os.Getenv("DEBUG") != "" {
		log.Println(message)
		return nil
	}
	return u.notification.Exec(message)
}
