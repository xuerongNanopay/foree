package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/account"
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
	RewardIds       []int64
	Rewards         []Reward
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
	OwnerId         int64

	CI   *InteracCITransaction
	IDM  *IDMTransaction
	COUT *NBPCOTransaction
}

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

type IDMTransaction struct {
	ID                   int64
	Status               TxStatus
	SrcInteracAccId      int64
	SrcInteracAcc        *InteracCITransaction
	DescContactAccountId int64
	DescContactAccount   *account.ForeeContactAccount
	Amt                  types.AmountData
	ParentTransactionId  int64
	OwnerId              int64
	CreateAt             time.Time `json:"createAt"`
	UpdateAt             time.Time `json:"updateAt"`
}

type NBPCOTransaction struct {
	ID                   int64
	Status               TxStatus
	Amt                  types.AmountData
	DescContactAccountId int64
	DescContactAccount   *account.ForeeContactAccount
	ParentTransactionId  int64
	OwnerId              int64
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
