package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLForeeTxInsert = `
		INSERT INTO foree_tx
		(
			type, status, rate,
			src_amount, src_currency, dest_amount, dest_currency
			total_amount, total_currency, cur_stage, cur_stage_status,
			conclusion, is_cancel_allowed, owner_id
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
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

// Only Support INITIAL, PROCESSING, CANCEL, COMPLETE.
type ForeeTx struct {
	ID              int64
	Type            string
	Status          TxStage
	Rate            float64
	SrcAmt          types.AmountData
	DestAmt         types.AmountData
	Total           types.AmountData
	CurStage        TxStage
	CurStageStatus  string
	Conclusion      string
	IsCancelAllowed bool `json:"isCancelAllowed"`
	OwnerId         int64
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`

	FeeIDs    []int64
	Fees      []FeeJoint
	RewardIds []int64
	Rewards   []Reward
	CI        *InteracCITx
	IDM       *IDMTx
	COUT      *NBPCOTx
	Summary   *TxSummary
	RefundTx  *ForeeRefundTx
	History   []*TxHistory
}

func NewForeeTxRepo(db *sql.DB) *ForeeTxRepo {
	return &ForeeTxRepo{db: db}
}

type ForeeTxRepo struct {
	db *sql.DB
}
