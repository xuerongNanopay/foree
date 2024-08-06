package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type InteracCITransaction struct {
	ID                  int64
	Status              TxStatus
	SrcInteracAccId     int64
	SrcInteracAcc       *InteracCITransaction
	DescInteracAccId    int64
	DescInteracAcc      *InteracCITransaction
	Amt                 types.AmountData
	ParentTransactionId int64
	OwnerId             int64
	CreateAt            time.Time `json:"createAt"`
	UpdateAt            time.Time `json:"updateAt"`
}
