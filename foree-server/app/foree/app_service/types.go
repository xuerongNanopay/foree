package foree_service

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ForeeDate struct {
	time.Time
}

func (d *ForeeDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d.Time)
	s := t.Format(time.DateOnly)
	return []byte("\"" + s + "\""), nil
}

func (d *ForeeDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		d.Time = time.Time{}
		return
	}
	d.Time, err = time.Parse(time.DateOnly, s)
	return
}
