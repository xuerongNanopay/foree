package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/constant"
)

const (
	sQLForeeRefundTxInsert = `
		INSERT INTO foree_refund_tx
		(
			status, refund_amount, refund_currency, parent_tx_id, owner_id
		) VALUES(?,?,?,?,?)
	`
	sQLForeeRefundTxUpdateById = `
		UPDATE foree_refund_tx SET
			status = ?
		where id = ?
	`
	sQLForeeRefundTxGetUniqueById = `
		SELECT
			t.id, t.status, refund_amount, refund_currency,
			t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
		FROM foree_refund_tx as t
		WHERE t.id = ?
	`
	sQLForeeRefundTxGetUniqueByParentTxId = `
		SELECT
			t.id, t.status, refund_amount, refund_currency,
			t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
		FROM foree_refund_tx as t
		WHERE t.parent_tx_id = ?
	`
)

type RefundTxStatus string

const (
	RefundTxStatusInitial  RefundTxStatus = "INITIAL"
	RefundTxStatusRefunded RefundTxStatus = "REFUNDED"
)

type ForeeRefundTx struct {
	ID         int64            `json:"id"`
	Status     RefundTxStatus   `json:"status"`
	RefundAmt  types.AmountData `json:"refundAmt"`
	ParentTxId int64            `json:"parentTxId"`
	OwnerId    int64            `json:"ownerId"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`

	RefundInteracAcc *account.InteracAccount `json:"refundInteracAcc"`
}

func NewForeeRefundTxRepo(db *sql.DB) *ForeeRefundTxRepo {
	return &ForeeRefundTxRepo{db: db}
}

type ForeeRefundTxRepo struct {
	db *sql.DB
}

func (repo *ForeeRefundTxRepo) InsertForeeRefundTx(ctx context.Context, tx ForeeRefundTx) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLForeeRefundTxInsert,
			tx.Status,
			tx.RefundAmt.Amount,
			tx.RefundAmt.Currency,
			tx.ParentTxId,
			tx.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLForeeRefundTxInsert,
			tx.Status,
			tx.RefundAmt.Amount,
			tx.RefundAmt.Currency,
			tx.ParentTxId,
			tx.OwnerId,
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

func (repo *ForeeRefundTxRepo) UpdateForeeRefundTxById(ctx context.Context, tx ForeeRefundTx) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(sQLForeeRefundTxUpdateById, tx.Status, tx.ID)

	} else {
		_, err = repo.db.Exec(sQLForeeRefundTxUpdateById, tx.Status, tx.ID)

	}
	if err != nil {
		return err
	}
	return nil
}

func (repo *ForeeRefundTxRepo) GetUniqueForeeRefundTxByParentTxId(ctx context.Context, parentTxId int64) (*ForeeRefundTx, error) {
	rows, err := repo.db.Query(sQLForeeRefundTxGetUniqueByParentTxId, parentTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *ForeeRefundTx

	for rows.Next() {
		f, err = scanRowForeeRefundTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *ForeeRefundTxRepo) GetUniqueForeeRefundTxById(ctx context.Context, id int64) (*ForeeRefundTx, error) {
	rows, err := repo.db.Query(sQLForeeRefundTxGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *ForeeRefundTx

	for rows.Next() {
		f, err = scanRowForeeRefundTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowForeeRefundTx(rows *sql.Rows) (*ForeeRefundTx, error) {
	tx := new(ForeeRefundTx)
	err := rows.Scan(
		&tx.ID,
		&tx.Status,
		&tx.RefundAmt.Amount,
		&tx.RefundAmt.Currency,
		&tx.ParentTxId,
		&tx.OwnerId,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
