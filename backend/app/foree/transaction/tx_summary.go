package transaction

import (
	"database/sql"
	"time"
)

const (
	sQLTxSummaryInsert = `
        INSERT INTO tx_summary
        (
            summary, type, status, rate, 
            src_acc_summary, src_amount, src_currency, 
            dest_acc_summary, dest_amount, dest_currency,
            total_amount, total_currency,
            fee_amount, fee_currency, 
            reward_amount, reward_currency, 
            is_cancel_allowed, parent_tx_id, owner_id
        ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
    `
	sQLTxSummaryGetUniqueById = `
        SELECT 
            t.id, t.summary, t.type, t.status, t.rate
            t.src_acc_summary, t.src_amount, t.src_currency, 
            t.dest_acc_summary, t.dest_amount, t.dest_currency,
            t.total_amount, t.total_currency,
            t.fee_amount, t.fee_currency, 
            t.reward_amount, t.reward_currency, 
            t.is_cancel_allowed, t.parent_tx_id, t.owner_id, 
            t.create_at, t.update_at
        FROM tx_summary t
        where t.id = ?
    `
)

type TxSummary struct {
	ID              int64     `json:"id"`
	Summary         string    `json:"sumary"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Rate            string    `json:"rate"`
	SrcAccSummary   string    `json:"srcAccSummary"`
	SrcAmount       string    `json:"srcAmount"`
	SrcCurrency     string    `json:"srcCurrency"`
	DestAccSummary  string    `json:"destAccSummary"`
	DestAmount      string    `json:"destAmount"`
	DestCurrency    string    `json:"destCurrency"`
	TotalAmount     string    `json:"totalAmount"`
	TotalCurrency   string    `json:"totalCurrency"`
	FeeAmount       string    `json:"feeAmount"`
	FeeCurrency     string    `json:"feeCurrency"`
	RewardAmount    string    `json:"rewardAmount"`
	RewardCurrency  string    `json:"rewardCurrency"`
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	ParentTxId      int64     `json:"parentTxd"`
	OwnerId         int64     `json:"owerId"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
}

func NewTxSummaryRepo(db *sql.DB) *TxSummaryRepo {
	return &TxSummaryRepo{db: db}
}

type TxSummaryRepo struct {
	db *sql.DB
}

func scanRowIntoTxSummary(rows *sql.Rows) (*TxSummary, error) {
	tx := new(TxSummary)
	err := rows.Scan(
		&tx.ID,
		&tx.Summary,
		&tx.Type,
		&tx.Status,
		&tx.Rate,
		&tx.SrcAccSummary,
		&tx.SrcAmount,
		&tx.SrcCurrency,
		&tx.DestAccSummary,
		&tx.DestAmount,
		&tx.DestCurrency,
		&tx.TotalAmount,
		&tx.TotalCurrency,
		&tx.FeeAmount,
		&tx.FeeCurrency,
		&tx.RewardAmount,
		&tx.RewardCurrency,
		&tx.IsCancelAllowed,
		&tx.ParentTxId,
		&tx.OwnerId,
		&tx.CreateAt,
		&tx.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
