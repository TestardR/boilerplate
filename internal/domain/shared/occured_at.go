package shared

import (
	"time"
)

type OccurredAt struct {
	value string
}

// OccurredAtFrom will:
// - assure conversion of any TZ to uniform UTC
// - assure equal time precision across all places in application
func OccurredAtFrom(t time.Time) OccurredAt {
	return OccurredAt{
		value: t.UTC().Format(time.RFC3339),
	}
}

func (t OccurredAt) AsTime() time.Time {
	value, err := time.Parse(time.RFC3339, t.value)
	if err != nil {
		panic("OccurredAt time value is not parsable to time.Time")
	}

	return value
}
