package nbp

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestEntityMarshal(t *testing.T) {
	t.Run("LoadRemittanceRequest should marshal correctly", func(t *testing.T) {
		test1 := `{"Token":"dummy","Agency_Code":"dummy","Remitter_Name":"xuerong"}`

		l := &LoadRemittanceRequest{
			requestCommon: requestCommon{
				Token:      "dummy",
				AgencyCode: "dummy",
			},
			RemitterName: "xuerong",
		}

		s, err := json.Marshal(l)
		if err != nil {
			t.Errorf("LoadReimttanceRequest marshal should not get err, but got `%v`", err.Error())
		}
		if test1 != string(s) {
			t.Errorf("LoadReimttanceRequest marshal incorrect, got `%v`", string(s))
		}
	})

	t.Run("NBPAmount should marshal correctly", func(t *testing.T) {
		test1 := `{"Token":"dummy","Agency_Code":"dummy","Amount":64.64,"Remitter_Name":"xuerong"}`

		l := &LoadRemittanceRequest{
			requestCommon: requestCommon{
				Token:      "dummy",
				AgencyCode: "dummy",
			},
			RemitterName: "xuerong",
			Amount:       64.644,
		}

		s, err := json.Marshal(l)
		if err != nil {
			t.Errorf("LoadReimttanceRequest marshal should not get err, but got `%v`", err.Error())
		}
		if test1 != string(s) {
			t.Errorf("LoadReimttanceRequest marshal incorrect, got `%v`", string(s))
		}
	})

	t.Run("NBPDate should marshal correctly", func(t *testing.T) {
		test1 := `{"Token":"dummy","Agency_Code":"dummy","Remitter_Name":"xuerong","remitter_DOB":"1989-06-04"}`
		d := NBPDate(time.Date(1989, time.June, 4, 0, 0, 0, 0, time.UTC))

		l := &LoadRemittanceRequest{
			requestCommon: requestCommon{
				Token:      "dummy",
				AgencyCode: "dummy",
			},
			RemitterName: "xuerong",
			RemitterDOB:  &d,
		}

		s, err := json.Marshal(l)
		if err != nil {
			t.Errorf("LoadReimttanceRequest marshal should not get err, but got `%v`", err.Error())
		}
		if test1 != string(s) {
			t.Errorf("LoadReimttanceRequest marshal incorrect, got `%v`", string(s))
		}
	})
}

func TestEntityUnMarshal(t *testing.T) {
	t.Run("TransactionStatusByIdsResponse should unmarshal correctly", func(t *testing.T) {
		test1 := `{"ResponseCode":"200","ResponseMessage":"Success","transactionStatuses":[{"Global_Id": "111111"}]}`

		res := &TransactionStatusByIdsResponse{
			ResponseCommon: ResponseCommon{
				RawResponse: test1,
			},
		}

		err := json.NewDecoder(bytes.NewBuffer([]byte(test1))).Decode(res)
		if err != nil {
			t.Errorf("TransactionStatusByIdsResponse unmarshal should not get err, but got `%v`", err.Error())
		}

		if res.ResponseCode != "200" {
			t.Errorf("expect `200` but got `%v`", res.ResponseCode)
		}

		if res.ResponseMessage != "Success" {
			t.Errorf("expect `Success` but got `%v`", res.ResponseMessage)
		}

		if res.RawResponse != test1 {
			t.Errorf("expect RawResponse persist, but got `%v`", res.RawResponse)
		}

		if len(res.TransactionStatuses) != 1 {
			t.Errorf("expect transactionStatuses length is `1`, but got `%v`", len(res.TransactionStatuses))
		}

		if res.TransactionStatuses[0].GlobalId != "111111" {
			t.Errorf("expect GlobalId of index 0 is `111111`, but got `%v`", res.TransactionStatuses[0].GlobalId)
		}
	})
}
