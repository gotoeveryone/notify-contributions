package mock

import (
	"time"

	"gotoeveryone/notify-contributions/src/domain/entity"
)

type ContributionClient struct {
	Entity *entity.Contribution
	Err    error
}

func (c *ContributionClient) Get(baseDate time.Time) (*entity.Contribution, error) {
	return c.Entity, c.Err
}

type NotificationClient struct {
	Err error
}

func (c *NotificationClient) Exec(message string) error {
	return c.Err
}
