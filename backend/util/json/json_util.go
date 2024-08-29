package json_util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DeserializeJsonFromHttpRequest(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func SerializeToResponseWriter(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
