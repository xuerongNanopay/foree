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
			transaction_purpose, conclusion, is_cancel_allowed, owner_id
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
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
	ID                 int64            `json:"id"`
	Type               string           `json:"type"`
	Status             TxStage          `json:"status"`
	Rate               types.Amount     `json:"Rate"`
	SrcAmt             types.AmountData `json:"srcAmt"`
	DestAmt            types.AmountData `json:"destAmt"`
	Total              types.AmountData `json:"total"`
	CurStage           TxStage          `json:"curStage"`
	CurStageStatus     string           `json:"curStageStatus"`
	TransactionPurpose string           `json:"transactionPurpose"`
	Conclusion         string           `json:"conclusion"`
	IsCancelAllowed    bool             `json:"isCancelAllowed"`
	OwnerId            int64            `json:"ownerId"`
	CreateAt           time.Time        `json:"createAt"`
	UpdateAt           time.Time        `json:"updateAt"`

	FeeIDs    []int64        `json:"feeIds"`
	Fees      []FeeJoint     `json:"fees"`
	RewardIds []int64        `json:"rewardIds"`
	Rewards   []Reward       `json:"rewards"`
	CI        *InteracCITx   `json:"ci"`
	IDM       *IDMTx         `json:"idm"`
	COUT      *NBPCOTx       `json:"cout"`
	Summary   *TxSummary     `json:"summary"`
	RefundTx  *ForeeRefundTx `json:"refundTx"`
	History   []*TxHistory   `json:"history"`
}

func NewForeeTxRepo(db *sql.DB) *ForeeTxRepo {
	return &ForeeTxRepo{db: db}
}

type ForeeTxRepo struct {
	db *sql.DB
}
