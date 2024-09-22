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
	sQLInteracCITxInsert = `
        INSERT INTO interact_ci_tx
        (
            status, cash_in_acc_id, amount, currency,
			end_to_end_id, parent_tx_id, owner_id) VALUES(?,?,?,?,?,?,?)
    `
	sQLInteracCITxUpdateById = `
		UPDATE interact_ci_tx SET 
			status = ?, scotia_payment_id = ?, scotia_status = ?,
			scotia_clearing_reference = ?, payment_url = ?
		WHERE id = ?
	`
	sQLInteracCITxGetForUpdateById = `
		SELECT 
			t.id, t.status, t.cash_in_acc_id,
			t.amount, t.currency, t.scotia_payment_id, 
			t.scotia_status, t.scotia_clearing_reference, t.payment_url, t.end_to_end_id,
			t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
		FROM interact_ci_tx t
		where t.id = ?
		FOR UPDATE
	`
	sQLInteracCITxGetUniqueById = `
        SELECT 
            t.id, t.status, t.cash_in_acc_id,
            t.amount, t.currency, t.scotia_payment_id, 
			t.scotia_status, t.scotia_clearing_reference, t.payment_url, t.end_to_end_id,
            t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
        FROM interact_ci_tx t
        where t.id = ?

    `
	sQLInteracCITxGetUniqueByParentTxId = `
        SELECT 
            t.id, t.status, t.cash_in_acc_id,
            t.amount, t.currency, t.scotia_payment_id, 
			t.scotia_status, t.scotia_clearing_reference, t.payment_url, t.end_to_end_id,
            t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
        FROM interact_ci_tx t
        where t.parent_tx_id = ?
    `
	sQLInteracCITxGetUniqueByScotiaPaymentId = `
		SELECT 
			t.id, t.status, t.cash_in_acc_id,
			t.amount, t.currency, t.scotia_payment_id, 
			t.scotia_status, t.scotia_clearing_reference, t.payment_url, t.end_to_end_id,
			t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
		FROM interact_ci_tx t
		where t.scotia_payment_id = ?
	`
	sQLInteracCITxGetAllByStatus = `
		SELECT 
			t.id, t.status, t.cash_in_acc_id,
			t.amount, t.currency, t.scotia_payment_id, 
			t.scotia_status, t.scotia_clearing_reference, t.payment_url, t.end_to_end_id,
			t.parent_tx_id, t.owner_id, t.created_at, t.updated_at
		FROM interact_ci_tx t
		where t.status = ?
	`
)

type InteracCITx struct {
	ID                      int64                   `json:"id"`
	Status                  TxStatus                `json:"status"`
	ScotiaPaymentId         string                  `json:"scotiaPaymentId"`
	ScotiaStatus            string                  `json:"scotiaStatus"`
	ScotiaClearingReference string                  `json:"scotiaClearingReference"`
	PaymentUrl              string                  `json:"paymentUrl"`
	EndToEndId              string                  `json:"endToEndId"`
	CashInAccId             int64                   `json:"CashInAccId"`
	CashInAcc               *account.InteracAccount `json:"CashInAcc"`
	Amt                     types.AmountData        `json:"Amt"`
	ParentTxId              int64                   `json:"parentTxId"`
	OwnerId                 int64                   `json:"OwnerId"`
	CreatedAt               time.Time               `json:"createdAt"`
	UpdatedAt               time.Time               `json:"updatedAt"`
}

func NewInteracCITxRepo(db *sql.DB) *InteracCITxRepo {
	return &InteracCITxRepo{db: db}
}

type InteracCITxRepo struct {
	db *sql.DB
}

func (repo *InteracCITxRepo) InsertInteracCITx(ctx context.Context, tx InteracCITx) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLInteracCITxInsert,
			tx.Status,
			tx.CashInAccId,
			tx.Amt.Amount,
			tx.Amt.Currency,
			tx.EndToEndId,
			tx.ParentTxId,
			tx.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLInteracCITxInsert,
			tx.Status,
			tx.CashInAccId,
			tx.Amt.Amount,
			tx.Amt.Currency,
			tx.EndToEndId,
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

func (repo *InteracCITxRepo) UpdateInteracCITxById(ctx context.Context, tx InteracCITx) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(sQLInteracCITxUpdateById, tx.Status, tx.ScotiaPaymentId, tx.ScotiaStatus, tx.ScotiaClearingReference, tx.PaymentUrl, tx.ID)

	} else {
		_, err = repo.db.Exec(sQLInteracCITxUpdateById, tx.Status, tx.ScotiaPaymentId, tx.ScotiaStatus, tx.ScotiaClearingReference, tx.PaymentUrl, tx.ID)

	}
	if err != nil {
		return err
	}
	return nil
}

func (repo *InteracCITxRepo) GetUniqueInteracCITxByParentTxId(ctx context.Context, parentTxId int64) (*InteracCITx, error) {
	rows, err := repo.db.Query(sQLInteracCITxGetUniqueByParentTxId, parentTxId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *InteracCITx

	for rows.Next() {
		f, err = scanRowIntoInteracCITx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *InteracCITxRepo) GetUniqueInteracCITxByScotiaPaymentId(ctx context.Context, scotiaPaymentId string) (*InteracCITx, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var rows *sql.Rows

	if ok {
		rows, err = dTx.Query(sQLInteracCITxGetUniqueByScotiaPaymentId, scotiaPaymentId)
	} else {
		rows, err = repo.db.Query(sQLInteracCITxGetUniqueByScotiaPaymentId, scotiaPaymentId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *InteracCITx

	for rows.Next() {
		f, err = scanRowIntoInteracCITx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *InteracCITxRepo) GetUniqueInteracCITxById(ctx context.Context, id int64) (*InteracCITx, error) {
	rows, err := repo.db.Query(sQLInteracCITxGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *InteracCITx

	for rows.Next() {
		f, err = scanRowIntoInteracCITx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *InteracCITxRepo) GetAllInteracCITxByStatus(ctx context.Context, status TxStatus) ([]*InteracCITx, error) {
	rows, err := repo.db.Query(sQLInteracCITxGetAllByStatus, status)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]*InteracCITx, 0)
	for rows.Next() {
		p, err := scanRowIntoInteracCITx(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return txs, nil
}

func scanRowIntoInteracCITx(rows *sql.Rows) (*InteracCITx, error) {
	tx := new(InteracCITx)
	err := rows.Scan(
		&tx.ID,
		&tx.Status,
		&tx.CashInAccId,
		&tx.Amt.Amount,
		&tx.Amt.Currency,
		&tx.ScotiaPaymentId,
		&tx.ScotiaStatus,
		&tx.ScotiaClearingReference,
		&tx.PaymentUrl,
		&tx.EndToEndId,
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
