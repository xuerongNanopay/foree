package time_util

import "time"

func NowInToronto() time.Time {
	if location, err := time.LoadLocation("America/Toronto"); err == nil {
		return time.Now().In(location)
	}
	return time.Now()
}

func UTCInToronto(t time.Time) time.Time {
	if location, err := time.LoadLocation("America/Toronto"); err == nil {
		return t.In(location)
	}
	return t
}
