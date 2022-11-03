package entity

import (
	"fmt"
	"time"
)

type Contribution struct {
	Type           string
	BaseDate       time.Time
	YesterdayCount int
	BaseDateCount  int
}

func (c *Contribution) Message() string {
	return fmt.Sprintf("[%s]\nContribute count to %s: %d.\nDifference from yesterday: %d.",
		c.Type, c.BaseDate.Format("2006-01-02"), c.BaseDateCount, c.BaseDateCount-c.YesterdayCount)
}
