package transaction

import (
	"database/sql"
	"time"
)

const ()

type TxHistory struct {
	ID         int64
	ParentTxId int64
	Stage      TxStage
	Status     TxStatus
	ExtraInfo  string
	OwnerId    int64
	CreateAt   time.Time `json:"createAt"`
}

func NewTxHistoryRepo(db *sql.DB) *TxHistoryRepo {
	return &TxHistoryRepo{db: db}
}

type TxHistoryRepo struct {
	db *sql.DB
}

func scanRowIntoTxHistory(rows *sql.Rows) (*TxHistory, error) {
	tx := new(TxHistory)
	err := rows.Scan(
		&tx.ID,
		&tx.ParentTxId,
		&tx.Stage,
		&tx.Status,
		&tx.ExtraInfo,
		&tx.OwnerId,
		&tx.CreateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
