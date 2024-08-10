package auth

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/constant"
)

// 3600 * 24 * 365 *19
const Second_In_Year = 31536000

// letters, spaces, number and extended latin
const NameReg = `^[a-zA-Z_0-9\u00C0-\u017F][a-zA-Z_0-9\u00C0-\u017F\s]*$`
const NineDigitReg = `^\\d{9}$`

var phoneNumberReplayer = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "", "+", "")
var validate = validator.New()

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

type SignUpReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ReferralCode string `json:"referralCode"`
}

type SessionReq struct {
	SessionId string `json:"sessionId"`
}

type VerifyEmailReq struct {
	SessionReq
	Code string `json:"code"`
}

type LoginReq struct {
	SessionReq
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgetPasswordUpdateReq struct {
	RetrieveCode string
	NewPassword  string
}

type CreateUserReq struct {
	FirstName           string    `json:"firstName" validate:"required"`
	MiddleName          string    `json:"middleName"`
	LastName            string    `json:"lastName" validate:"required"`
	Age                 int       `json:"age" validate:"gte=19,lte=130"`
	Dob                 ForeeDate `json:"dob"`
	Nationality         string    `json:"nationality" validate:"required"`
	Address1            string    `json:"address1" validate:"required"`
	Address2            string    `json:"address2"`
	City                string    `json:"city" validate:"required"`
	Province            string    `json:"province" validate:"required"`
	Country             string    `json:"country" validate:"required"`
	PhoneNumber         string    `json:"phoneNumber" validate:"required"`
	IdentificationType  string    `json:"identificationType" validate:"required"`
	IdentificationValue string    `json:"identificationValue" validate:"required"`
	AvatarUrl           string    `json:"avatarUrl"`
}

// TODO: trim name, and use allowText
func (q *CreateUserReq) Validate() *transport.BadRequestError {
	ret := transport.NewFormError("Invalid user creation request")
	if err := validate.Struct(q); err != nil {
		errors := err.(validator.ValidationErrors)
		ret = transport.NewFormError("Invalid user creation request")
		for _, e := range errors {
			ret.AddDetails(e.Field(), e.Error())
		}
	}

	age := q.Dob.Time.Unix() / int64(Second_In_Year)

	if age < 19 || age > 120 {
		ret.AddDetails("dob", "illegal age")
	}

	q.Age = int(age)

	if q.Country != "CA" {
		ret.AddDetails("country", fmt.Sprintf("invalid country `%v`", q.Country))
	}

	_, ok := constant.Regions["CA"][q.Province]
	if !ok {
		ret.AddDetails("province", fmt.Sprintf("invalid province `%v`", q.Province))
	}

	if len(ret.Details) > 0 {
		return ret
	}

	phoneNumber := phoneNumberReplayer.Replace(q.PhoneNumber)
	ok, _ = regexp.MatchString(NineDigitReg, phoneNumber)
	if !ok {
		ret.AddDetails("phoneNumber", fmt.Sprintf("invalid phone number `%v`", q.PhoneNumber))
	}
	q.PhoneNumber = phoneNumber

	return nil
}

// func allowText(input string) bool {

// }
