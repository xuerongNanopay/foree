package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type TxLimit struct {
	ID         string
	Amt        types.Amount
	IsMinLimit bool
	IsEnable   bool
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
}
