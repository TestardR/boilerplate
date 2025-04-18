package testshared

import "time"

type FixedClock struct {
	fixedTime time.Time
}

func NewFixedClock(fixedTime time.Time) FixedClock {
	return FixedClock{fixedTime: fixedTime}
}

func (c FixedClock) Now() time.Time {
	return c.fixedTime
}
