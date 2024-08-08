package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

// TODO
type TxLimitCache struct {
	ID       int64
	Identity string
	UsdAmt   types.Amount
	MaxAmt   types.Amount
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}
