package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

type IDMTransaction struct {
	ID               int64
	Status           TxStatus
	Ip               string `json:"ip"`
	UserAgent        string `json:"userAgent"`
	SrcInteracAccId  int64
	SrcInteracAcc    *account.InteracAccount
	DescContactAccId int64
	DescContactAcc   *account.ContactAccount
	Amt              types.AmountData
	ParentTxId       int64
	OwnerId          int64
	CreateAt         time.Time `json:"createAt"`
	UpdateAt         time.Time `json:"updateAt"`
}

// Large object.
type IDMCompliance struct {
	ID            int64
	IDMTxId       int64
	IDMStatusCode int
	IDMResult     string
	RequestJson   string
	ResponseJson  string
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
}
