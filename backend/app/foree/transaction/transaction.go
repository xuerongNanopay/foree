package foree_transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type TxStage string

const (
	TxStageCI TxStage = "CASH_IN"
)

type ForeeTransaction struct {
	ID      int64
	SrcAmt  types.AmountData
	DestAmt types.AmountData
	Rate    types.RateDate
	Status  string

	FeeIDs       []int64
	Fees         []FeeJoint
	PromotionIds []int64
	// Promotions   []types.Promotion
	Total           types.AmountData
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
	OwnerId         int64

	CI   *CITransaction
	IDM  *IDMTransaction
	COUT *NBPTransaction
}

type CITransaction struct {
	ForeeTransactionId int64
	OwnerId            int64
}

type IDMTransaction struct {
	ForeeTransactionId int64
	OwnerId            int64
}

type NBPTransaction struct {
	ForeeTransactionId int64
	OwnerId            int64
}

// type
// Src
// Dest
type TransactionSummary struct {
	ID          string
	Summary     string
	SrcAmt      types.AmountData
	DestAmt     types.AmountData
	CiIsoStatus string
	Type        string
	Status      string
	FXRate      string

	Fees []Fee
	// Promotions []types.Promotion
	Total   types.AmountData
	Created time.Time
}
