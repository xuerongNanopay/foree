package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

type NBPCOTransaction struct {
	ID                   int64
	Status               TxStatus
	Amt                  types.AmountData
	DescContactAccountId int64
	DescContactAccount   *account.ContactAccount
	ParentTransactionId  int64
	OwnerId              int64
	CreateAt             time.Time `json:"createAt"`
	UpdateAt             time.Time `json:"updateAt"`
}
