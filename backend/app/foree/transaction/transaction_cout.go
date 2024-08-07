package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

type NBPCOTransaction struct {
	ID               int64
	Status           TxStatus
	Amt              types.AmountData
	DescContactAccId int64
	DescContactAcc   *account.ContactAccount
	ParentTxId       int64
	OwnerId          int64
	CreateAt         time.Time `json:"createAt"`
	UpdateAt         time.Time `json:"updateAt"`
}
