package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

type TxStatus string

const (
	TxStatusInitial    TxStatus = "INITIAL"
	TxStatusProcessing TxStatus = "PROCESSING"
	TxStatusSuspend    TxStatus = "Suspend"
	TxStatusSent       TxStatus = "Sent"
	TxStatusReject     TxStatus = "REJECT"
	TxStatusCancel     TxStatus = "CANCEL"
	TxStatusComplete   TxStatus = "COMPLETE"
)

type TxType string

const (
	TxTypeInteracToNBP TxType = "INTERAC-NBP"
)

type TxStage string

const (
	TxStageInteracCI TxStage = "INTERAC-CI"
	TxStageIDM       TxStage = "Compliance-IDM"
	TxStageNBPCI     TxStage = "INTERAC-CO"
)

// Only Support PROCESSING, CANCEL, COMPLETE.
type ForeeTx struct {
	ID             int64
	Type           string
	SrcAmt         types.AmountData
	DestAmt        types.AmountData
	Rate           types.RateDate
	Status         TxStage
	Total          types.AmountData
	CurStage       TxStage
	CurStageStatus string
	Conclusion     string

	FeeIDs          []int64
	Fees            []FeeJoint
	RewardIds       []int64
	Rewards         []Reward
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
	OwnerId         int64

	CI       *InteracCITx
	IDM      *IDMTx
	COUT     *NBPCOTx
	Summary  *TxSummary
	RefundTx *ForeeRefundTx
	History  []*TxHistory
}

func NewForeeTxRepo(db *sql.DB) *ForeeTxRepo {
	return &ForeeTxRepo{db: db}
}

type ForeeTxRepo struct {
	db *sql.DB
}
