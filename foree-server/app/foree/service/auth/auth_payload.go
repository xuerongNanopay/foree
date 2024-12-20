package foree_auth_service

import (
	"fmt"
	"regexp"
	"time"

	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/server/transport"
	"xue.io/go-pay/validation"
)

type SignUpReq struct {
	Email             string `json:"email" validate:"required,email"`
	Password          string `json:"password" validate:"required"`
	ReferrerReference string `json:"referrerReference"`
}

func (q SignUpReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid sign up request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type UpdatePasswdReq struct {
	transport.SessionReq
	OldPasswd string `json:"oldPasswd" validate:"required,min=8,max=16"`
	NewPasswd string `json:"newPasswd" validate:"required,min=8,max=16"`
}

func (q UpdatePasswdReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid change passwd request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type VerifyEmailReq struct {
	transport.SessionReq
	Code string `json:"code"`
}

func (q VerifyEmailReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid verify email request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type LoginReq struct {
	transport.SessionReq
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (q LoginReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid login request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type CreateUserReq struct {
	transport.SessionReq
	FirstName           string          `json:"firstName" validate:"required"`
	MiddleName          string          `json:"middleName"`
	LastName            string          `json:"lastName" validate:"required"`
	Age                 int             `json:"age"`
	Dob                 types.ForeeDate `json:"dob"`
	Pob                 string          `json:"pob" validate:"required"`
	Nationality         string          `json:"nationality" validate:"required"`
	Address1            string          `json:"address1" validate:"required"`
	Address2            string          `json:"address2"`
	City                string          `json:"city" validate:"required"`
	Province            string          `json:"province" validate:"required"`
	Country             string          `json:"country" validate:"required"`
	PostalCode          string          `json:"postalCode" validate:"required"`
	PhoneNumber         string          `json:"phoneNumber" validate:"required"`
	IdentificationType  string          `json:"identificationType" validate:"required"`
	IdentificationValue string          `json:"identificationValue" validate:"required"`
	AvatarUrl           string          `json:"avatarUrl"`
}

// TODO: trim name, and use allowText
func (q CreateUserReq) Validate() *transport.BadRequestError {
	ret := validation.ValidateStruct(q, "Invalid user creation request")

	// Age
	age := q.Dob.Time.Unix() / (int64(time.Hour/time.Second) * 24 * 365)

	if age < 19 || age > 120 {
		ret.AddDetails("dob", "illegal age")
	}

	q.Age = int(age)

	// Country
	if q.Country != "CA" {
		ret.AddDetails("country", fmt.Sprintf("invalid country `%v`", q.Country))
	}

	// Province
	_, ok := constant.Regions["CA"][q.Province]
	if !ok {
		ret.AddDetails("province", fmt.Sprintf("invalid province `%v`", q.Province))
	}

	country := constant.Countires[q.Country]

	// Postal code
	ok, _ = regexp.MatchString(country.PostalCodeRegex, q.PostalCode)
	if !ok {
		ret.AddDetails("postalCode", fmt.Sprintf("invalid postal code `%v`", q.PostalCode))
	}

	// Phone number
	ok, _ = regexp.MatchString(country.PhoneRegex, q.PhoneNumber)
	if !ok {
		ret.AddDetails("phoneNumber", fmt.Sprintf("invalid phone number `%v`", q.PhoneNumber))
	}

	// Identification type
	_, ok = foree_constant.AllowIdentificationTypes[foree_auth.IdentificationType(q.IdentificationType)]
	if !ok {
		ret.AddDetails("identificationType", fmt.Sprintf("invalid identificationType `%v`", q.IdentificationType))
	}

	if len(ret.Details) > 0 {
		return ret
	}

	return nil
}

type ForgetPasswdReq struct {
	Email string `json:"email" validate:"email,required"`
}

func (q ForgetPasswdReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid forget password request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type ForgetPasswdVerifyReq struct {
	Email        string `json:"email" validate:"email,required"`
	RetrieveCode string `json:"retrieveCode" validate:"required"`
}

func (q ForgetPasswdVerifyReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid forget password request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type ForgetPasswdUpdateReq struct {
	Email        string `json:"email" validate:"email,required"`
	RetrieveCode string `json:"retrieveCode" validate:"required"`
	NewPassword  string `json:"newPassword" validate:"required"`
}

func (q ForgetPasswdUpdateReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid new password request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type UpdateAddressReq struct {
	transport.SessionReq
	Address1   string `json:"address1" validate:"required"`
	Address2   string `json:"address2"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	Country    string `json:"country" validate:"required"`
	PostalCode string `json:"postalCode" validate:"required"`
}

func (q UpdateAddressReq) Validate() *transport.BadRequestError {
	ret := validation.ValidateStruct(q, "Invalid update address request")

	// Country
	if q.Country != "CA" {
		ret.AddDetails("country", fmt.Sprintf("invalid country `%v`", q.Country))
	}

	// Province
	_, ok := constant.Regions["CA"][q.Province]
	if !ok {
		ret.AddDetails("province", fmt.Sprintf("invalid province `%v`", q.Province))
	}

	country := constant.Countires[q.Country]

	// Postal code
	ok, _ = regexp.MatchString(country.PostalCodeRegex, q.PostalCode)
	if !ok {
		ret.AddDetails("postalCode", fmt.Sprintf("invalid postal code `%v`", q.PostalCode))
	}

	// Phone number
	// ok, _ = regexp.MatchString(country.PhoneRegex, q.PhoneNumber)
	// if !ok {
	// 	ret.AddDetails("phoneNumber", fmt.Sprintf("invalid phone number `%v`", q.PhoneNumber))
	// }

	if len(ret.Details) > 0 {
		return ret
	}

	return nil
}

type UpdatePhoneNumberReq struct {
	transport.SessionReq
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

func (q UpdatePhoneNumberReq) Validate() *transport.BadRequestError {
	ret := validation.ValidateStruct(q, "Invalid update phoneNumber request")

	country := constant.Countires["CA"]

	// Phone number
	ok, _ := regexp.MatchString(country.PhoneRegex, q.PhoneNumber)
	if !ok {
		ret.AddDetails("phoneNumber", fmt.Sprintf("invalid phone number `%v`", q.PhoneNumber))
	}

	if len(ret.Details) > 0 {
		return ret
	}

	return nil
}

type UpdateUserSetting struct {
	transport.SessionReq
	IsInAppNotificationEnable  bool `json:"isInAppNotificationEnable"`
	IsPushNotificationEnable   bool `json:"isPushNotificationEnable"`
	IsEmailNotificationsEnable bool `json:"isEmailNotificationsEnable"`
}

func (q UpdateUserSetting) Validate() *transport.BadRequestError {
	ret := validation.ValidateStruct(q, "Invalid update userSetting request")

	if len(ret.Details) > 0 {
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

type UserDetailDTO struct {
	ID          int64           `json:"id"`
	Status      auth.UserStatus `json:"status"`
	FirstName   string          `json:"firstName"`
	MiddleName  string          `json:"middleName"`
	LastName    string          `json:"lastName"`
	Age         int             `json:"age"`
	Dob         *time.Time      `json:"dob"`
	Address1    string          `json:"address1"`
	Address2    string          `json:"address2"`
	City        string          `json:"city"`
	Province    string          `json:"province"`
	Country     string          `json:"country"`
	PostalCode  string          `json:"postalCode"`
	PhoneNumber string          `json:"phoneNumber"`
	Email       string          `json:"email"`
	AvatarUrl   string          `json:"avatarUrl"`
	CreatedAt   int64           `json:"createdAt"`
}

func NewUserDetailDTO(u *auth.User) *UserDetailDTO {
	ret := &UserDetailDTO{
		ID:          u.ID,
		Status:      u.Status,
		FirstName:   u.FirstName,
		MiddleName:  u.MiddleName,
		LastName:    u.LastName,
		Age:         u.Age,
		Dob:         u.Dob,
		Address1:    u.Address1,
		Address2:    u.Address2,
		City:        u.City,
		Province:    u.Province,
		Country:     u.Country,
		PostalCode:  u.PostalCode,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		AvatarUrl:   u.AvatarUrl,
	}
	if u.CreatedAt != nil {
		ret.CreatedAt = u.CreatedAt.UnixMilli()
	}
	return ret
}

type UserSettingDTO struct {
	IsInAppNotificationEnable  bool  `json:"isInAppNotificationEnable"`
	IsPushNotificationEnable   bool  `json:"isPushNotificationEnable"`
	IsEmailNotificationsEnable bool  `json:"isEmailNotificationsEnable"`
	OwnerId                    int64 `json:"ownerId"`
}

func NewUserSettingDTO(us *auth.UserSetting) *UserSettingDTO {
	return &UserSettingDTO{
		IsInAppNotificationEnable:  us.IsInAppNotificationEnable,
		IsPushNotificationEnable:   us.IsPushNotificationEnable,
		IsEmailNotificationsEnable: us.IsEmailNotificationsEnable,
		OwnerId:                    us.OwnerId,
	}
}

type UserExtraDTO struct {
	UserReference      string `json:"userReference,omitempty"`
	Pob                string `json:"pob,omitempty"`
	Cor                string `json:"cor,omitempty"`
	Nationality        string `json:"nationality,omitempty"`
	OccupationCategory string `json:"occupationCategory,omitempty"`
	OccupationName     string `json:"occupationName,omitempty"`
}

func NewUserExtraDTO(ue *foree_auth.UserExtra) *UserExtraDTO {
	return &UserExtraDTO{
		UserReference:      ue.UserReference,
		Pob:                ue.Pob,
		Cor:                ue.Cor,
		Nationality:        ue.Nationality,
		OccupationCategory: ue.OccupationCategory,
		OccupationName:     ue.OccupationName,
	}
}
