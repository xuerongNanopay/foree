package auth

import (
	"strings"
	"time"
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
	FirstName   string    `json:"firstName"`
	MiddleName  string    `json:"middleName"`
	LastName    string    `json:"lastName"`
	Age         int       `json:"age"`
	Dob         ForeeDate `json:"dob"`
	Nationality string    `json:"nationality"`
	Address1    string    `json:"address1"`
	Address2    string    `json:"address2"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	Country     string    `json:"country"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
	AvatarUrl   string    `json:"avatarUrl"`
}
