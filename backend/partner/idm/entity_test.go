package idm

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestEntityMarshal(t *testing.T) {
	t.Run("IDMRequest should marshal correctly", func(t *testing.T) {
		test1 := `{"tea":"dummy@xrw.io","dob":"1989-06-04","amt":112.11,"tags":["aaa","bbb"],"memo19":"usa"}`
		d := IDMDate(time.Date(1989, time.June, 4, 0, 0, 0, 0, time.UTC))

		req := IDMRequest{
			UserEmail:   "dummy@xrw.io",
			Dob:         &d,
			Amount:      112.114,
			Tags:        []string{"aaa", "bbb"},
			RemitterCOR: "usa",
		}

		s, err := json.Marshal(req)
		if err != nil {
			t.Errorf("IDMRequest marshal should not get err, but got `%v`", err.Error())
		}
		if test1 != string(s) {
			t.Errorf("IDMRequest marshal incorrect, got `%v`", string(s))
		}
	})
}

func TestEntityUnMarshal(t *testing.T) {
	t.Run("IDMResponse should unmarshal correctly", func(t *testing.T) {
		test1 := `{"frp":"ACCEPT","ednaScoreCard":{}}`
		resp := &IDMResponse{
			ResponseCommon: ResponseCommon{
				StatusCode:  200,
				RawResponse: test1,
			},
		}

		err := json.NewDecoder(bytes.NewBuffer([]byte(test1))).Decode(resp)
		if err != nil {
			t.Errorf("IDMResponse unmarshal should not get err, but got `%v`", err.Error())
		}

		if resp.StatusCode != 200 {
			t.Errorf("expect `200` but got `%v`", resp.StatusCode)
		}

		if resp.FraudEvaluationResult != ResultStatusAccept {
			t.Errorf("expect `%v` but got `%v`", ResultStatusAccept, resp.FraudEvaluationResult)
		}

		if resp.RawResponse != test1 {
			t.Errorf("expect RawResponse persist, but got `%v`", resp.RawResponse)
		}
	})
}
