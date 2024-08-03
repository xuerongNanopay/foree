package nbp

import "time"

const PAKISTAN_STANDARD_TIME_ZONE = "+05:00"

func parseTokenExpiryDate(rawExpiry string) (*time.Time, error) {
	pkTime := rawExpiry + PAKISTAN_STANDARD_TIME_ZONE

	t, err := time.Parse(time.RFC3339, pkTime)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func isValidToken(auth *tokenData, threshold int64) bool {
	if auth == nil || auth.token == "" || auth.tokenExpiry == nil {
		return false
	}

	if time.Now().Unix()+threshold >= auth.tokenExpiry.Unix() {
		return false
	}

	return true
}
