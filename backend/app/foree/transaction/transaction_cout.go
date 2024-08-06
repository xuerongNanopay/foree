package transaction

import (
	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

type NBPCOTransaction struct {
	ID                   int64
	Status               TxStatus
	Amt                  types.AmountData
	DescContactAccountId int64
	DescContactAccount   *account.ForeeContactAccount
	ParentTransactionId  int64
	OwnerId              int64
}
