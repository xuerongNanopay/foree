package transport

// Define generic format for HTTP transfermation.
// StatusCode should follow standard http status code
// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
type HTTPResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	Error      any    `json:"error"`
}
