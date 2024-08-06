package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

type IDMTransaction struct {
	ID                   int64
	Status               TxStatus
	SrcInteracAccId      int64
	SrcInteracAcc        *ScotiaInteracCITransaction
	DescContactAccountId int64
	DescContactAccount   *account.ForeeContactAccount
	Amt                  types.AmountData
	ParentTransactionId  int64
	OwnerId              int64
	CreateAt             time.Time `json:"createAt"`
	UpdateAt             time.Time `json:"updateAt"`
}
