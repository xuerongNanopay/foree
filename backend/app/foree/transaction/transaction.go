package transaction

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
	RewardIds       []int64
	Rewards         []Reward
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
	OwnerId         int64

	CI      *ScotiaInteracCITransaction
	IDM     *IDMTransaction
	COUT    *NBPCOTransaction
	Summary *TransactionSummary
	History *TransactionStatusHistory
}
