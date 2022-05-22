package entity

import (
	"fmt"
	"time"
)

type Contribution struct {
	UserName       string
	BaseDate       time.Time
	YesterdayCount int
	BaseDateCount  int
}

func (c *Contribution) Message() string {
	return fmt.Sprintf("Contribute count to GitHub on %s by %s: %d.\nDifference from yesterday: %d.",
		c.BaseDate.Format("2006-01-02"), c.UserName, c.BaseDateCount, c.BaseDateCount-c.YesterdayCount)
}
