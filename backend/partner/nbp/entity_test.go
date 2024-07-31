package nbp

import (
	"encoding/json"
	"testing"
)

func TestEntityMarshal(t *testing.T) {
	t.Run("LoadRemittanceRequest should marshal correctly", func(t *testing.T) {
		test1 := `{"Token":"dummy","Agency_Code":"dummy","Remitter_Name":"xuerong","remitter_DOB":"1989-06-04"}`

		l := &LoadRemittanceRequest{
			requestCommon: requestCommon{
				Token:      "dummy",
				AgencyCode: "dummy",
			},
			RemitterName: "xuerong",
			RemitterDOB:  "1989-06-04",
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
