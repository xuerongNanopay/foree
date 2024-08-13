package service

import (
	"fmt"
	"regexp"
	"strings"

	fAuth "xue.io/go-pay/app/foree/auth"
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
)

type SignUpReq struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	ReferralCode string `json:"referralCode"`
}

func (q *SignUpReq) TrimSpace() {
	q.Email = strings.TrimSpace(q.Email)
	q.Password = strings.TrimSpace(q.Password)
}

func (q *SignUpReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	if ret := validateStruct(q, "Invalid sign up request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type ChangePasswdReq struct {
	transport.SessionReq
	Password string `json:"password" validate:"required,min=8,max=12"`
}

func (q *ChangePasswdReq) TrimSpace() {
	q.Password = strings.TrimSpace(q.Password)
}

func (q *ChangePasswdReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	if ret := validateStruct(q, "Invalid change password request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type VerifyEmailReq struct {
	transport.SessionReq
	Code string `json:"code"`
}

func (q *VerifyEmailReq) TrimSpace() {
	q.Code = strings.TrimSpace(q.Code)
}

func (q *VerifyEmailReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	if ret := validateStruct(q, "Invalid verify email request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type LoginReq struct {
	transport.SessionReq
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=12"`
}

func (q *LoginReq) TrimSpace() {
	q.Email = strings.TrimSpace(q.Email)
	q.Password = strings.TrimSpace(q.Password)
}

func (q *LoginReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	if ret := validateStruct(q, "Invalid login request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type ForgetPasswordUpdateReq struct {
	RetrieveCode string
	NewPassword  string
}

type CreateUserReq struct {
	transport.SessionReq
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
	PostalCode          string    `json:"postalCode" validate:"required"`
	PhoneNumber         string    `json:"phoneNumber" validate:"required"`
	IdentificationType  string    `json:"identificationType" validate:"required"`
	IdentificationValue string    `json:"identificationValue" validate:"required"`
	AvatarUrl           string    `json:"avatarUrl"`
}

func (q *CreateUserReq) TrimSpace() {
	q.FirstName = strings.TrimSpace(q.FirstName)
	q.MiddleName = strings.TrimSpace(q.MiddleName)
	q.LastName = strings.TrimSpace(q.LastName)
	q.Nationality = strings.TrimSpace(q.Nationality)
	q.Address1 = strings.TrimSpace(q.Address1)
	q.Address2 = strings.TrimSpace(q.Address2)
	q.City = strings.TrimSpace(q.City)
	q.Province = strings.TrimSpace(q.Province)
	q.Country = strings.TrimSpace(q.Country)
	q.PostalCode = strings.TrimSpace(q.PostalCode)
	q.PhoneNumber = strings.TrimSpace(q.PhoneNumber)
	q.IdentificationType = strings.TrimSpace(q.IdentificationType)
	q.IdentificationValue = strings.TrimSpace(q.IdentificationValue)
	q.AvatarUrl = strings.TrimSpace(q.AvatarUrl)
}

// TODO: trim name, and use allowText
func (q *CreateUserReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid user creation request")

	// Age
	age := q.Dob.Time.Unix() / int64(Second_In_Year)

	if age < 19 || age > 120 {
		ret.AddDetails("dob", "illegal age")
	}

	q.Age = int(age)

	// Country/Region
	if q.Country != "CA" {
		ret.AddDetails("country", fmt.Sprintf("invalid country `%v`", q.Country))
	}

	_, ok := constant.Regions["CA"][q.Province]
	if !ok {
		ret.AddDetails("province", fmt.Sprintf("invalid province `%v`", q.Province))
	}

	//TODO: Postal Code

	// Phone number
	phoneNumber := phoneNumberReplayer.Replace(q.PhoneNumber)
	ok, _ = regexp.MatchString(NineDigitReg, phoneNumber)
	if !ok {
		ret.AddDetails("phoneNumber", fmt.Sprintf("invalid phone number `%v`", q.PhoneNumber))
	}
	q.PhoneNumber = phoneNumber

	// Identification type
	_, ok = allowIdentificationTypes[fAuth.IdentificationType(q.IdentificationType)]
	if !ok {
		ret.AddDetails("identificationType", fmt.Sprintf("invalid identificationType `%v`", q.IdentificationType))
	}

	if len(ret.Details) > 0 {
		return ret
	}

	return nil
}

// --------------- Response ------------------

type UserDTO struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Status     string `json:"status"`
	AvatarUrl  string `json:"avatarUrl"`
}

func NewUserDTO(user *auth.User) *UserDTO {
	return &UserDTO{
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Status:     string(user.Status),
		AvatarUrl:  user.AvatarUrl,
	}
}
