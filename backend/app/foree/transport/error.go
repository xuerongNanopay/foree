package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Severity string

const FormErrorSignUpMsg = "Invaild Signup Request"

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

type RequireAction string

const (
	RequireActionLogin       RequireAction = "LOGIN"
	RequireActionVerifyEmail RequireAction = "VERIFY_EMAIL"
	RequireActionCreateUser  RequireAction = "CREATE_USER"
)

const (
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
