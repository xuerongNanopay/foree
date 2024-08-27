package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/constant"
)

const (
	sQLInteracRefundTxInsert = `
		INSERT INTO interac_refund_tx
		(
			status, refund_interac_acc_id, parent_tx_id, owner_id
		) VALUES(?,?,?,?)
	`
	sQLInteracRefundTxUpdateById = `
		UPDATE interac_refund_tx SET
			status = ?
		where id = ?
	`
	sQLInteracRefundTxGetUniqueById = `
		SELECT
			t.id, t.status, t.refund_interac_acc_id, t.parent_tx_id,
			t.owner_id, t.create_at, t.update_at
		FROM interac_refund_tx as t
		WHERE t.id = ?
	`
	sQLInteracRefundTxGetUniqueByParentTxId = `
		SELECT
			t.id, t.status, t.refund_interac_acc_id, t.parent_tx_id,
			t.owner_id, t.create_at, t.update_at
		FROM interac_refund_tx as t
		WHERE t.parent_tx_id = ?
	`
)

type RefundTxStatus string

const (
	RefundTxStatusInitial   RefundTxStatus = "INITIAL"
	RefundTxStatusRefunding RefundTxStatus = "REFUNDING"
	RefundTxStatusRefunded  RefundTxStatus = "REFUNDED"
)

type InteracRefundTx struct {
	ID                 int64          `json:"id"`
	Status             RefundTxStatus `json:"status"`
	RefundInteracAccId int64          `json:"refundInteracAccId"`
	ParentTxId         int64          `json:"parentTxId"`
	OwnerId            int64          `json:"ownerId"`
	CreateAt           time.Time      `json:"createAt"`
	UpdateAt           time.Time      `json:"updateAt"`

	RefundInteracAcc *account.InteracAccount `json:"refundInteracAcc"`
}

func NewInteracRefundTxRepo(db *sql.DB) *InteracRefundTxRepo {
	return &InteracRefundTxRepo{db: db}
}

type InteracRefundTxRepo struct {
	db *sql.DB
}

func (repo *InteracRefundTxRepo) InsertInteracRefundTx(ctx context.Context, tx InteracRefundTx) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLInteracRefundTxInsert,
			tx.Status,
			tx.RefundInteracAccId,
			tx.ParentTxId,
			tx.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLInteracRefundTxInsert,
			tx.Status,
			tx.RefundInteracAccId,
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

func (repo *InteracRefundTxRepo) UpdateInteracRefundTxById(ctx context.Context, tx InteracRefundTx) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(sQLInteracRefundTxUpdateById, tx.Status, tx.ID)

	} else {
		_, err = repo.db.Exec(sQLInteracRefundTxUpdateById, tx.Status, tx.ID)

	}
	if err != nil {
		return err
	}
	return nil
}

func (repo *InteracRefundTxRepo) GetUniqueInteracRefundTxByParentTxId(ctx context.Context, parentTxId int64) (*InteracRefundTx, error) {
	rows, err := repo.db.Query(sQLInteracRefundTxGetUniqueByParentTxId, parentTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *InteracRefundTx

	for rows.Next() {
		f, err = scanRowInteracRefundTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *InteracRefundTxRepo) GetUniqueInteracRefundTxById(ctx context.Context, id int64) (*InteracRefundTx, error) {
	rows, err := repo.db.Query(sQLInteracRefundTxGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *InteracRefundTx

	for rows.Next() {
		f, err = scanRowInteracRefundTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowInteracRefundTx(rows *sql.Rows) (*InteracRefundTx, error) {
	tx := new(InteracRefundTx)
	err := rows.Scan(
		&tx.ID,
		&tx.Status,
		&tx.RefundInteracAccId,
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
