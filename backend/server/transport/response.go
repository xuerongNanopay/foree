package transport

// See: https://www.bugsnag.com/blog/go-errors/ for error stacktrace
// Define generic format for HTTP transfermation.
// StatusCode should follow standard http status code
// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status

// List codes that are used often.
// 200 OK
// 201 Created
// 400 Bad Request: Mainly used in form submit.
// 401 Unauthorized: client need get new token or re-login.
// 428 Precondition Require
//
// 403 Forbidden: The client does not have access rights to the content
// 503 Service Unavailable

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func NewHttpResponse(statusCode int, message string, data any) *HttpResponse {
	return &HttpResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}
