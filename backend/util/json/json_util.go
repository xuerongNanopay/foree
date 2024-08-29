package json_util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJsonFromHttpRequest(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}
