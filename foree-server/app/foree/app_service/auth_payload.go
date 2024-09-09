package foree_service

import (
	"fmt"
	"regexp"
	"time"

	fAuth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/server/transport"
)

type SignUpReq struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	ReferralCode string `json:"referralCode"`
}

func (q SignUpReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid sign up request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type ChangePasswdReq struct {
	transport.SessionReq
	OldPassword string `json:"oldPassword" validate:"required,min=8,max=16"`
	Password    string `json:"password" validate:"required,min=8,max=16"`
}

func (q ChangePasswdReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid change password request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type VerifyEmailReq struct {
	transport.SessionReq
	Code string `json:"code"`
}

func (q VerifyEmailReq) Validate() *transport.BadRequestError {
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

func (q LoginReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid login request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type CreateUserReq struct {
	transport.SessionReq
	FirstName           string    `json:"firstName" validate:"required"`
	MiddleName          string    `json:"middleName"`
	LastName            string    `json:"lastName" validate:"required"`
	Age                 int       `json:"age"`
	Dob                 ForeeDate `json:"dob"`
	Pob                 string    `json:"pob" validate:"required"`
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

// TODO: trim name, and use allowText
func (q CreateUserReq) Validate() *transport.BadRequestError {
	ret := validateStruct(q, "Invalid user creation request")

	// Age
	age := q.Dob.Time.Unix() / (int64(time.Hour/time.Second) * 24 * 365)

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
	ok, _ = regexp.MatchString(foree_constant.NineDigitReg, phoneNumber)
	if !ok {
		ret.AddDetails("phoneNumber", fmt.Sprintf("invalid phone number `%v`", q.PhoneNumber))
	}
	q.PhoneNumber = phoneNumber

	// Identification type
	_, ok = foree_constant.AllowIdentificationTypes[fAuth.IdentificationType(q.IdentificationType)]
	if !ok {
		ret.AddDetails("identificationType", fmt.Sprintf("invalid identificationType `%v`", q.IdentificationType))
	}

	if len(ret.Details) > 0 {
		return ret
	}

	return nil
}

type ForgetPasswordUpdateReq struct {
	Email        string `json:"email" validate:"email,required"`
	RetrieveCode string `json:"retrieveCode" validate:"required"`
	NewPassword  string `json:"newPassword" validate:"required"`
}

func (q ForgetPasswordUpdateReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid new password request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type ForgetPasswordReq struct {
	Email string `json:"email" validate:"email,required"`
}

func (q ForgetPasswordReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid forget password request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

// --------------- Response ------------------
func NewUserDTO(session *auth.Session) *UserDTO {
	ret := &UserDTO{
		SessionId: session.ID,
	}

	if session.EmailPasswd != nil {
		ret.LoginStatus = session.EmailPasswd.Status
	}

	if session.User != nil {
		ret.UserStatus = session.User.Status
		ret.FirstName = session.User.FirstName
		ret.MiddleName = session.User.MiddleName
		ret.LastName = session.User.LastName
		ret.AvatarUrl = session.User.AvatarUrl
	}

	return ret
}

type UserDTO struct {
	SessionId   string                 `json:"sessionId,omitempty"`
	LoginStatus auth.EmailPasswdStatus `json:"loginStatus,omitempty"`
	UserStatus  auth.UserStatus        `json:"userStatus,omitempty"`
	FirstName   string                 `json:"firstName,omitempty"`
	MiddleName  string                 `json:"middleName,omitempty"`
	LastName    string                 `json:"lastName,omitempty"`
	AvatarUrl   string                 `json:"avatarUrl,omitempty"`
}
