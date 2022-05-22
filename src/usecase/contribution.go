package usecase

import (
	"gotoeveryone/notify-github-contributions/src/domain/client"
	"time"
)

type Contribution interface {
	Exec(identifier string, baseDate time.Time) error
}

type contributionUsecase struct {
	source       client.Contribution
	notification client.Notification
}

func NewContributionUsecase(source client.Contribution, notification client.Notification) Contribution {
	return &contributionUsecase{
		source:       source,
		notification: notification,
	}
}

// Exec is get contribution and notify
func (u *contributionUsecase) Exec(identifier string, baseDate time.Time) error {
	c, err := u.source.Get(identifier, baseDate)
	if err != nil {
		return err
	}
	return u.notification.Exec(c.Message())
}
