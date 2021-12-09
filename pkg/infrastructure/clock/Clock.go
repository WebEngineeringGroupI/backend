package clock

import (
	"time"
)

type System struct {
}

func (c *System) Now() time.Time {
	return time.Now().UTC()
}
