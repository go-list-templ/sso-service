package vo

import "time"

var TTL = 7 * 24 * time.Hour

type ExpiresAt struct {
	value time.Time
}

func NewExpiresAt(timeNow time.Time) ExpiresAt {
	return ExpiresAt{
		value: timeNow.Add(TTL),
	}
}

func UnsafeExpiresAt(time time.Time) ExpiresAt {
	return ExpiresAt{value: time}
}

func (e *ExpiresAt) Value() time.Time {
	return e.value
}

func (e *ExpiresAt) Expired() bool {
	return time.Now().After(e.value)
}
