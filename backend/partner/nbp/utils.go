package nbp

import "time"

const PAKISTAN_STANDARD_TIME_ZONE = "+05:00"

func parseTokenExpiryDate(rawExpiry string) (time.Time, error) {
	pkTime := rawExpiry + PAKISTAN_STANDARD_TIME_ZONE

	t, err := time.Parse(time.RFC3339, pkTime)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
