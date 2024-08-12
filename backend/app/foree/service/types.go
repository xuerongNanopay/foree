package service

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"xue.io/go-pay/app/foree/auth"
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

// 3600 * 24 * 365 *19
const Second_In_Year = 31536000

// letters, spaces, number and extended latin
const NameReg = `^[a-zA-Z_0-9\u00C0-\u017F][a-zA-Z_0-9\u00C0-\u017F\s]*$`
const NineDigitReg = `^\\d{9}$`

var phoneNumberReplayer = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "", "+", "")
var validate = validator.New()

var allowIdentificationTypes = map[auth.IdentificationType]bool{
	auth.IDTypePassport:      true,
	auth.IDTypeDriverLicense: true,
	auth.IDTypeProvincalId:   true,
	auth.IDTypeNationId:      true,
}
