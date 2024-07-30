package nbp

import (
	"testing"
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

		_, err := parseTokenExpiryDate(test1)
		if err == nil {
			t.Errorf("expected parse failed, but got %s", err.Error())
		}
	})
}
