package transport

import (
	"encoding/json"
	"testing"
)

func TestHTTPResponseMarshar(t *testing.T) {
	type Foo struct {
		Bar string
	}

	t.Run("HTTPResponse should marshal correctly", func(t *testing.T) {
		test1 := `{"statusCode":200,"message":"success","data":{"Bar":"AAAA"},"error":null}`
		foo := &Foo{
			Bar: "AAAA",
		}

		response := HTTPResponse{
			StatusCode: 200,
			Message:    "success",
			Data:       foo,
		}
		s, err := json.Marshal(response)
		if err != nil {
			t.Errorf("HTTPResponse marshal should not get err, but got `%v`", err.Error())
		}
		if test1 != string(s) {
			t.Errorf("HTTPResponse marshal incorrect, got `%v`", string(s))
		}
	})
}
