package nbp

import (
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
		test1 := `{"Token":"dummy","Agency_Code":"dummy","Amount":64.64,"Remitter_Name":"xuerong","remitter_DOB":"1989-06-04"}`
		d := NBPDate(time.Date(1989, time.June, 4, 0, 0, 0, 0, time.UTC))

		l := &LoadRemittanceRequest{
			requestCommon: requestCommon{
				Token:      "dummy",
				AgencyCode: "dummy",
			},
			RemitterName: "xuerong",
			RemitterDOB:  &d,
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
}

func TestEntityUnMarshal(t *testing.T) {

}
