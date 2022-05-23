package client

import (
	"gotoeveryone/notify-github-contributions/src/domain/entity"
	"time"
)

type Contribution interface {
	Get(identifier string, baseDate time.Time) (*entity.Contribution, error)
}
