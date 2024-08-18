package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

const (
	sQLForeeTxInsert = `
		INSERT INTO foree_tx
		(
			type, status, rate, cin_acc_id, cout_acc_id,
			src_amount, src_currency, dest_amount, dest_currency
			total_fee_amount, total_fee_currency, total_reward_amount, total_reward_currency,
			total_amount, total_currency, cur_stage, cur_stage_status,
			transaction_purpose, conclusion, owner_id
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	sQLForeeTxUpdateById = `
	    UPDATE foree_tx SET 
			status = ?, cur_stage = ?, cur_stage_status = ?, conclusion = ?
        WHERE id = ?
	`
	sQLForeeTxGetById = `
	    SELECT 
            t.id, t.type, t.status, t.rate
			t.cin_acc_id, t.cout_acc_id,
            t.src_amount, t.src_currency, 
            t.dest_amount, t.dest_currency,
			t.total_fee_amount, t.total_fee_currency, 
            t.total_reward_amount, t.total_reward_currency, 
            t.total_amount, t.total_currency,
			t.cur_stage, t.cur_stage_status, t.transaction_purpose, t.conclusion,
            t.owner_id, t.create_at, t.update_at
        FROM foree_tx t
        where t.id = ?
	`
	//TODO: support get alls?
)

type TxStatus string

const (
	TxStatusInitial    TxStatus = "INITIAL"
	TxStatusProcessing TxStatus = "PROCESSING"
	TxStatusSuspend    TxStatus = "SUSPEND"
	TxStatusSent       TxStatus = "SENT"
	TxStatusRejected   TxStatus = "REJECTED"
	TxStatusCancelled  TxStatus = "CANCELLED"
	TxStatusCompleted  TxStatus = "COMPLETED"
	TxStatusClosed     TxStatus = "CLOSED"
)

type TxType string

const (
	TxTypeInteracToNBP TxType = "INTERAC-NBP"
)

type TxStage string

const (
	TxStageInteracCI TxStage = "INTERAC-CI"
	TxStageIDM       TxStage = "Compliance-IDM"
	TxStageNBPCO     TxStage = "NBP-CO"
)

// Only Support INITIAL, PROCESSING, CANCEL, COMPLETE.
type ForeeTx struct {
	ID                 int64            `json:"id"`
	Type               TxType           `json:"type"`
	Status             TxStatus         `json:"status"`
	Rate               types.Amount     `json:"Rate"`
	CinAccId           int64            `json:"cinAccId"`
	CoutAccId          int64            `json:"coutAccId"`
	SrcAmt             types.AmountData `json:"srcAmt"`
	DestAmt            types.AmountData `json:"destAmt"`
	TotalFeeAmt        types.AmountData `json:"totalFeeAmt"`
	TotalRewardAmt     types.AmountData `json:"totalRewardAmt"`
	TotalAmt           types.AmountData `json:"totalAmt"`
	CurStage           TxStage          `json:"curStage"`
	CurStageStatus     TxStatus         `json:"curStageStatus"`
	TransactionPurpose string           `json:"transactionPurpose"`
	Conclusion         string           `json:"conclusion"`
	OwnerId            int64            `json:"ownerId"`
	CreateAt           time.Time        `json:"createAt"`
	UpdateAt           time.Time        `json:"updateAt"`

	InteracAcc  *account.InteracAccount `json:"interacAcc"`
	ContactAcc  *account.ContactAccount `json:"contactAcc"`
	FeeJointIds []int64                 `json:"feeJointIds"`
	Fees        []*FeeJoint             `json:"fees"`
	RewardIds   []int64                 `json:"rewardIds"`
	Rewards     []*Reward               `json:"rewards"`
	CI          *InteracCITx            `json:"ci"`
	IDM         *IDMTx                  `json:"idm"`
	COUT        *NBPCOTx                `json:"cout"`
	Summary     *TxSummary              `json:"summary"`
	RefundTx    *ForeeRefundTx          `json:"refundTx"`
	History     []*TxHistory            `json:"history"`
}

func NewForeeTxRepo(db *sql.DB) *ForeeTxRepo {
	return &ForeeTxRepo{db: db}
}

type ForeeTxRepo struct {
	db *sql.DB
}

func (repo *ForeeTxRepo) InsertForeeTx(ctx context.Context, tx ForeeTx) (int64, error) {
	result, err := repo.db.Exec(
		sQLForeeTxInsert,
		tx.Type,
		tx.Status,
		tx.Rate,
		tx.CinAccId,
		tx.CoutAccId,
		tx.SrcAmt.Amount,
		tx.SrcAmt.Currency,
		tx.DestAmt.Amount,
		tx.DestAmt.Currency,
		tx.TotalFeeAmt.Amount,
		tx.TotalFeeAmt.Currency,
		tx.TotalRewardAmt.Amount,
		tx.TotalRewardAmt.Currency,
		tx.TotalAmt.Amount,
		tx.TotalAmt.Currency,
		tx.CurStage,
		tx.CurStageStatus,
		tx.TransactionPurpose,
		tx.Conclusion,
		tx.OwnerId,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *ForeeTxRepo) UpdateForeeTxById(ctx context.Context, tx ForeeTx) error {
	_, err := repo.db.Exec(sQLForeeTxUpdateById, tx.Status, tx.CurStage, tx.CurStageStatus, tx.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ForeeTxRepo) GetUniqueForeeTxById(ctx context.Context, id int64) (*ForeeTx, error) {
	rows, err := repo.db.Query(sQLForeeTxGetById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *ForeeTx

	for rows.Next() {
		f, err = scanRowIntoForeeTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoForeeTx(rows *sql.Rows) (*ForeeTx, error) {
	tx := new(ForeeTx)
	err := rows.Scan(
		&tx.ID,
		&tx.Type,
		&tx.Status,
		&tx.Rate,
		&tx.CinAccId,
		&tx.CoutAccId,
		&tx.SrcAmt.Amount,
		&tx.SrcAmt.Currency,
		&tx.DestAmt.Amount,
		&tx.DestAmt.Currency,
		&tx.TotalFeeAmt.Amount,
		&tx.TotalFeeAmt.Currency,
		&tx.TotalRewardAmt.Amount,
		&tx.TotalRewardAmt.Currency,
		&tx.TotalAmt.Amount,
		&tx.TotalAmt.Currency,
		&tx.CurStage,
		&tx.CurStageStatus,
		&tx.TransactionPurpose,
		&tx.Conclusion,
		&tx.OwnerId,
		&tx.CreateAt,
		&tx.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
