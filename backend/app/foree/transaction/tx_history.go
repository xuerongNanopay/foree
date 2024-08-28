package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLTxHistoryInsert = `
        INSERT INTO tx_history
        (
            stage, status, extra_info, parent_tx_id, owner_id
        ) VALUES(?,?,?,?,?)
    `
	sQLTxHistoryGetAllByParentTxId = `
        SELECT
            h.id, h.stage, h.status, h.extra_info,
            h.parent_tx_id, h.owner_id, h.created_at
        FROM tx_history h
        where h.parent_tx_id = ?
    `
)

type TxHistory struct {
	ID         int64     `json:"id"`
	Stage      TxStage   `json:"stage"`
	Status     TxStatus  `json:"status"`
	ExtraInfo  string    `json:"extraInfo"`
	ParentTxId int64     `json:"parentTxId"`
	OwnerId    int64     `json:"ownerId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func NewTxHistory(tx *ForeeTx, extraInfo string) *TxHistory {
	return &TxHistory{
		Stage:      tx.CurStage,
		Status:     tx.CurStageStatus,
		ExtraInfo:  extraInfo,
		ParentTxId: tx.ID,
		OwnerId:    tx.OwnerId,
		CreatedAt:  time.Now(),
	}
}

func NewTxHistoryRepo(db *sql.DB) *TxHistoryRepo {
	return &TxHistoryRepo{db: db}
}

type TxHistoryRepo struct {
	db *sql.DB
}

func (repo *TxHistoryRepo) InserTxHistory(ctx context.Context, h TxHistory) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLTxHistoryInsert,
			h.Stage,
			h.Stage,
			h.ExtraInfo,
			h.ParentTxId,
			h.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLTxHistoryInsert,
			h.Stage,
			h.Stage,
			h.ExtraInfo,
			h.ParentTxId,
			h.OwnerId,
		)
	}

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *TxHistoryRepo) GetAllTxHistoryByTransactionId(parentTxId int64) ([]*TxHistory, error) {
	rows, err := repo.db.Query(sQLTxHistoryGetAllByParentTxId, parentTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hist := make([]*TxHistory, 16)
	for rows.Next() {
		p, err := scanRowIntoTxHistory(rows)
		if err != nil {
			return nil, err
		}
		hist = append(hist, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hist, nil
}

func scanRowIntoTxHistory(rows *sql.Rows) (*TxHistory, error) {
	tx := new(TxHistory)
	err := rows.Scan(
		&tx.ID,
		&tx.Stage,
		&tx.Status,
		&tx.ExtraInfo,
		&tx.ParentTxId,
		&tx.OwnerId,
		&tx.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
