package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/account"
)

type ForeeRefundTx struct {
	ParentTxId         int64
	OwnerId            int64
	RefundType         string
	IsRefund           bool
	RefundAt           time.Time
	RefundInteracAccId int64
	RefundInteracAcc   *account.InteracAccount
	CreateAt           time.Time `json:"createAt"`
	UpdateAt           time.Time `json:"updateAt"`
}
