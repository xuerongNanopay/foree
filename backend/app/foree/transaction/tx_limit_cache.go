package transaction

import (
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/types"
)

type TxLimitCache struct {
	ID       int64
	Identity string
	UsedAmt  types.Amount
	MaxAmt   types.Amount
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

func generateIdentity(referenceId int64) string {
	now := time.Now()
	loc, err := time.LoadLocation("America/Toronto")
	if err == nil {
		now = now.In(loc)
	}
	return fmt.Sprintf("%v_%s", referenceId, now.Format(time.DateOnly))
}
