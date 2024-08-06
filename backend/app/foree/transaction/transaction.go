package foree_transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type TxStatus string

const (
	TxStatusInitial TxStatus = "INITIAL"
	// TxStatusInitial TxStatus = "INITIAL"
)

type TxType string

const (
	TxTypeInteracToNBP TxType = "INTERAC-NBP"
)

type ForeeTransaction struct {
	ID      int64
	Type    string
	SrcAmt  types.AmountData
	DestAmt types.AmountData
	Rate    types.RateDate
	Status  string
	Total   types.AmountData

	FeeIDs          []int64
	Fees            []FeeJoint
	PromotionIds    []int64
	Promotions      []Promotion
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
	OwnerId         int64

	CI   *InteracCITransaction
	IDM  *IDMTransaction
	COUT *NBPCOTransaction
}

type InteracCITransaction struct {
	ForeeTransactionId int64
	OwnerId            int64
}

type IDMTransaction struct {
	ForeeTransactionId int64
	OwnerId            int64
}

type NBPCOTransaction struct {
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
