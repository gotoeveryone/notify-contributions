package client

import (
	"gotoeveryone/notify-contributions/src/domain/entity"
	"time"
)

type Contribution interface {
	Get(baseDate time.Time) (*entity.Contribution, error)
}
