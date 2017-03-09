package testdata

import "time"

type Clock interface {
	Now() (now time.Time)
}
