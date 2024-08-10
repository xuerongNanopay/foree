package transport

import (
	"time"

	"xue.io/go-pay/auth"
)

// See: https://www.bugsnag.com/blog/go-errors/ for error stacktrace
// Define generic format for HTTP transfermation.
// StatusCode should follow standard http status code
// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status

// List codes that are used often.
// 200 OK
// 201 Created
// 400 Bad Request: Mainly used in form submit.
// 401 Unauthorized: client need get new token or re-login.
// 412 Precondition Failed
//
// 403 Forbidden: The client does not have access rights to the content
// 503 Service Unavailable

type HTTPResponse struct {
	StatusCode int      `json:"statusCode"`
	Message    string   `json:"message"`
	EPStatus   string   `json:"epStatus"`
	User       *UserDTO `json:"user"`
	Data       any      `json:"data"`
	Error      any      `json:"error"`
}

type UserDTO struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"firstName"`
	MiddleName  string    `json:"middleName"`
	LastName    string    `json:"lastName"`
	Status      string    `json:"status"`
	Age         int       `json:"age"`
	Dob         time.Time `json:"dob"`
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

func NewUserDTO(user *auth.User) *UserDTO {
	if user == nil {
		return nil
	}
	return &UserDTO{
		ID:         user.ID,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Status:     string(user.Status),
		AvatarUrl:  user.AvatarUrl,
	}
}
