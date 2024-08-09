package transaction

import (
	"database/sql"
	"time"
)

const (
	sQLTxSummaryGetUniqueById = `
        SELECT 
            t.id, t.summary, t.type, t.status,
            t.src_acc_summary, t.src_amount, t.src_currency, 
            t.dest_acc_summary, t.dest_amount, t.dest_currency,
            t.is_cancel_allowed, t.parent_tx_id, t.owner_id, t.create_at, t.update_at
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
	DestCurrency    string    `json:"desCurrency"`
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
		&tx.SrcAccSummary,
		&tx.SrcAmount,
		&tx.SrcCurrency,
		&tx.DestAccSummary,
		&tx.DestAmount,
		&tx.DestCurrency,
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
