package scotia

import "time"

func isValidToken(auth *tokenData, threshold int64) bool {
	if auth == nil || auth.token == "" || auth.tokenExpiry == nil {
		return false
	}

	if time.Now().Unix()+threshold >= auth.tokenExpiry.Unix() {
		return false
	}

	return true
}
