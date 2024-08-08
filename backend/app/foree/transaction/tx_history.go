package transaction

import (
	"database/sql"
	"time"
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
			h.parent_tx_id, h.owner_id, h.create_at
		FROM tx_history h
		where h.parent_tx_id = ?
	`
)

type TxHistory struct {
	ID         int64
	Stage      TxStage
	Status     TxStatus
	ExtraInfo  string
	ParentTxId int64
	OwnerId    int64
	CreateAt   time.Time `json:"createAt"`
}

func NewTxHistoryRepo(db *sql.DB) *TxHistoryRepo {
	return &TxHistoryRepo{db: db}
}

type TxHistoryRepo struct {
	db *sql.DB
}

func (repo *TxHistoryRepo) InserTxHistory(h TxHistory) (int64, error) {
	result, err := repo.db.Exec(
		sQLTxHistoryInsert,
		h.Stage,
		h.Stage,
		h.ExtraInfo,
		h.ParentTxId,
		h.OwnerId,
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
		&tx.CreateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
