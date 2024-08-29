package service

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"xue.io/go-pay/server/transport"
)

// TODO: testing
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

func validateStruct(s any, errMsg string) *transport.BadRequestError {
	ret := transport.NewFormError(errMsg)
	if err := validate.Struct(s); err != nil {
		errors := err.(validator.ValidationErrors)
		for _, e := range errors {
			ret.AddDetails(e.Field(), e.Error())
		}
	}
	return ret
}
