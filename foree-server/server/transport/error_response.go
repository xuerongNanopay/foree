package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HError interface {
	GetStatusCode() int
	GetMessage() string
	Error() string
}

type Severity string

const (
	SeverityInfo  Severity = "INFO"
	SeverityWarn  Severity = "WARN"
	SeverityError Severity = "Error"
)

func NewFormError(message string, errors ...string) *BadRequestError {
	details := make([]ErrorDetail, len(errors)/2)
	for i := 0; i < len(errors); i += 2 {
		details = append(details, ErrorDetail{
			Severity: SeverityError,
			Field:    errors[i],
			Message:  errors[i+1],
		})
	}
	return &BadRequestError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Details:    details,
		Timestamp:  time.Now(),
	}
}

// 400 Bad Request
// eg: malformed request syntax, invalid request message framing, or deceptive request routing
type BadRequestError struct {
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Details    []ErrorDetail `json:"details"`
	Timestamp  time.Time     `json:"timestamp"`
}

type ErrorDetail struct {
	Severity Severity `json:"dseverity"`
	Field    string   `json:"field"`
	Message  string   `json:"message"`
}

func (b *BadRequestError) GetStatusCode() int {
	return b.StatusCode
}

func (b *BadRequestError) GetMessage() string {
	return b.Message
}

func (b *BadRequestError) Error() string {
	return serializeError(b)
}

func (b *BadRequestError) AddDetails(errors ...string) *BadRequestError {
	for i := 0; i < len(errors); i += 2 {
		b.Details = append(b.Details, ErrorDetail{
			Severity: SeverityError,
			Field:    errors[i],
			Message:  errors[i+1],
		})
	}
	return b
}

func (b *BadRequestError) AddDetailWithSeverity(severity Severity, field, message string) *BadRequestError {
	b.Details = append(b.Details, ErrorDetail{
		Severity: severity,
		Field:    field,
		Message:  message,
	})
	return b
}

type UnauthorizedRequestError struct {
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Require    RequireAction `json:"require"`
}

func NewUnauthorizedRequestError() *UnauthorizedRequestError {
	return &UnauthorizedRequestError{
		StatusCode: http.StatusUnauthorized,
		Message:    PreconditionRequireMsgLogin,
		Require:    RequireActionLogin,
	}
}

func (b *UnauthorizedRequestError) GetStatusCode() int {
	return b.StatusCode
}

func (b *UnauthorizedRequestError) GetMessage() string {
	return b.Message
}

func (b *UnauthorizedRequestError) Error() string {
	return serializeError(b)
}

type RequireAction string

const (
	RequireActionToMain      RequireAction = "TO_MAIN"
	RequireActionLogin       RequireAction = "LOGIN"
	RequireActionVerifyEmail RequireAction = "VERIFY_EMAIL"
	RequireActionCreateUser  RequireAction = "CREATE_USER"
)

const (
	PreconditionRequireMsgToMain      string = "Please navigate to main menu."
	PreconditionRequireMsgLogin       string = "Please login."
	PreconditionRequireMsgVerifyEmail string = "Please verify your email."
	PreconditionRequireMsgCreateUser  string = "Please fullfill your information."
)

type PreconditionRequireError struct {
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Require    RequireAction `json:"require"`
}

func NewPreconditionRequireError(message string, require RequireAction) *PreconditionRequireError {
	return &PreconditionRequireError{
		StatusCode: http.StatusPreconditionRequired,
		Message:    message,
		Require:    require,
	}
}

func (b *PreconditionRequireError) GetStatusCode() int {
	return b.StatusCode
}

func (b *PreconditionRequireError) GetMessage() string {
	return b.Message
}

func (b *PreconditionRequireError) Error() string {
	return serializeError(b)
}

func serializeError(e any) string {
	s, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%v", e)
	}
	return string(s)
}

// 403 Forbidden
type ForbiddenError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func NewForbiddenError(requirePermission string) *ForbiddenError {
	return &ForbiddenError{
		StatusCode: http.StatusForbidden,
		Message:    fmt.Sprintf("No `%v` permission.", requirePermission),
	}
}

func (b *ForbiddenError) GetStatusCode() int {
	return b.StatusCode
}

func (b *ForbiddenError) GetMessage() string {
	return b.Message
}

func (b *ForbiddenError) Error() string {
	return serializeError(b)
}

// 500 Internal Server Error
type InteralServerError struct {
	StatusCode    int    `json:"statusCode"`
	Message       string `json:"message"`
	OriginalError error  `json:"-"`
}

func (b *InteralServerError) GetStatusCode() int {
	return b.StatusCode
}

func (b *InteralServerError) GetMessage() string {
	return b.Message
}

func (b *InteralServerError) Error() string {
	return serializeError(b)
}

func WrapInteralServerError(e error) *InteralServerError {
	return &InteralServerError{
		StatusCode:    http.StatusInternalServerError,
		Message:       "Internal Server Error",
		OriginalError: e,
	}
}

func NewInteralServerError(format string, a ...any) *InteralServerError {
	e := fmt.Errorf(format, a...)
	return &InteralServerError{
		StatusCode:    http.StatusInternalServerError,
		Message:       "Internal Server Error",
		OriginalError: e,
	}
}
