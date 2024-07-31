package idm

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEntityMarshal(t *testing.T) {
	t.Run("IDMRequest should marshal correctly", func(t *testing.T) {
		test1 := `{"tea":"dummy@xrw.io","dob":"1989-06-04","amt":112.11,"tags":["aaa","bbb"],"memo19":"usa"}`
		d := IDMDate(time.Date(1989, time.June, 4, 0, 0, 0, 0, time.UTC))

		req := IDMRequest{
			UserEmail:   "dummy@xrw.io",
			Dob:         &d,
			Amount:      112.114,
			Tags:        []string{"aaa", "bbb"},
			RemitterCOR: "usa",
		}

		s, err := json.Marshal(req)
		if err != nil {
			t.Errorf("IDMRequest marshal should not get err, but got `%v`", err.Error())
		}
		if test1 != string(s) {
			t.Errorf("IDMRequest marshal incorrect, got `%v`", string(s))
		}
	})
}
