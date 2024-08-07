package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type ScotiaInteracCITransaction struct {
	ID                  int64
	Status              TxStatus
	ScotialId           string
	Url                 string
	SrcInteracAccId     int64
	SrcInteracAcc       *ScotiaInteracCITransaction
	DescInteracAccId    int64
	DescInteracAcc      *ScotiaInteracCITransaction
	Amt                 types.AmountData
	ParentTransactionId int64
	OwnerId             int64
	CreateAt            time.Time `json:"createAt"`
	UpdateAt            time.Time `json:"updateAt"`
}