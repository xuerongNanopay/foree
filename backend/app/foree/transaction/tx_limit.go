package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLTxLimitGet = `
		SELECT
			h.id, h.stage, h.status, h.extra_info,
			h.parent_tx_id, h.owner_id
		FROM tx_history h
		where h.parent_tx_id = ?
	`
)

type TxLimit struct {
	ID         string
	Amt        types.Amount
	IsMinLimit bool
	IsEnable   bool
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
}

func NewTxLimitRepo(db *sql.DB) *TxLimitRepo {
	return &TxLimitRepo{db: db}
}

type TxLimitRepo struct {
	db *sql.DB
}

func scanRowIntoTxLimit(rows *sql.Rows) (*TxHistory, error) {
	tx := new(TxHistory)
	err := rows.Scan(
		&tx.ID,
		&tx.Stage,
		&tx.Status,
		&tx.ExtraInfo,
		&tx.ParentTxId,
		&tx.OwnerId,
		&tx.CreateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
