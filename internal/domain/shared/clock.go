package shared

import "time"

type CurrentTime interface {
	Now() time.Time
}
