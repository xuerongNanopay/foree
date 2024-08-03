package nbp

import (
	"testing"
	"time"
)

func TestParseTokenExpiryDate(t *testing.T) {

	t.Run("1989-06-04T00:00:00 should parse successfully", func(t *testing.T) {
		test1 := "1989-06-04T00:00:00"
		timestamp1 := int64(612903600000)

		d, err := parseTokenExpiryDate(test1)
		if err != nil {
			t.Errorf("expected parse successfully, but got %s", err.Error())
		}

		if d.UnixMilli() != timestamp1 {
			t.Errorf("expected %v, but got %v", timestamp1, d.UnixMilli())
		}

	})

	t.Run("1989-06-04T00:00:000 should parse successfully", func(t *testing.T) {
		test1 := "1989-06-04T00:00:00"
		timestamp1 := int64(612903600000)

		d, err := parseTokenExpiryDate(test1)
		if err != nil {
			t.Errorf("expected parse successfully, but got %s", err.Error())
		}

		if d.UnixMilli() != timestamp1 {
			t.Errorf("expected %v, but got %v", timestamp1, d.UnixMilli())
		}

	})

	t.Run("1989-6-04T00:00:000 should parse fail", func(t *testing.T) {
		test1 := "1989-6-04T00:00:00"

		p, err := parseTokenExpiryDate(test1)
		if err == nil {
			t.Errorf("expected parse failed, but got %v", p)
		}
	})
}

func TestIsTokenAvailable(t *testing.T) {

	t.Run("nil authCache should return false", func(t *testing.T) {
		var cache *tokenData
		if !isValidToken(cache, 0) {
			t.Errorf("nil tokenData expect false, but got true")
		}
	})

	t.Run("empty token should return false", func(t *testing.T) {
		cache := &tokenData{}

		if !isValidToken(cache, 0) {
			t.Errorf("empty token expect false, but got true")
		}
	})

	t.Run("2100-6-04T00:00:000UTC should return true", func(t *testing.T) {
		d := time.Date(2100, time.June, 4, 0, 0, 0, 0, time.UTC)
		cache := &tokenData{
			token:       "dummy",
			tokenExpiry: &d,
		}
		if !isValidToken(cache, 0) {
			t.Errorf("2100-6-04T00:00:000UTC expect true, but got false")
		}
	})

	t.Run("1989-6-04T00:00:000UTC should return true", func(t *testing.T) {
		d := time.Date(1989, time.June, 4, 0, 0, 0, 0, time.UTC)
		cache := &tokenData{
			token:       "dummy",
			tokenExpiry: &d,
		}
		if !isValidToken(cache, 0) {
			t.Errorf("1989-6-04T00:00:000UTC expect false, but got true")
		}
	})
}
