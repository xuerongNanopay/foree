package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/partner/nbp"
)

const (
	sQLNBPCOTxInsert = `
        INSERT INTO nbp_co_tx
        (
            status, mode, amount, currency, nbp_reference,
			cash_out_acc_id, parent_tx_id, owner_id
        ) VALUES(?,?,?,?,?,?,?,?)
    `
	sQLNBPCOTxUpdateById = `
        UPDATE nbp_co_tx SET 
            status = ?
        WHERE id = ?
    `
	sQLNBPCOTxGetUniqueById = `
        SELECT 
            t.id, t.status, t.mode, t.amount, t.currency, t.nbp_reference,
            t.cash_out_acc_id, t.parent_tx_id, t.owner_id,
            t.created_at, t.updated_at
        FROM nbp_co_tx t
        where t.id = ?

    `
	sQLNBPCOTxGetUniqueByNBPReference = `
        SELECT 
            t.id, t.status, t.mode, t.amount, t.currency, t.nbp_reference,
            t.cash_out_acc_id, t.parent_tx_id, t.owner_id,
            t.created_at, t.updated_at
        FROM nbp_co_tx t
        where t.nbp_reference = ?

    `
	sQLNBPCOTxGetUniqueByParentTxId = `
        SELECT 
            t.id, t.status, t.mode, t.amount, t.currency, t.nbp_reference,
            t.cash_out_acc_id, t.parent_tx_id, t.owner_id,
            t.created_at, t.updated_at
        FROM nbp_co_tx t
        where t.parent_tx_id = ?
    `
)

type NBPCOTx struct {
	ID           int64                   `json:"id"`
	Status       TxStatus                `json:"status"`
	Mode         nbp.PMTMode             `json:"mode"`
	Amt          types.AmountData        `json:"amt"`
	NBPReference string                  `json:"nbpReference"`
	CashOutAccId int64                   `json:"CashOutAccId"`
	CashOutAcc   *account.ContactAccount `json:"CashOutAcc"`
	ParentTxId   int64                   `json:"parentTxId"`
	OwnerId      int64                   `json:"OwnerId"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
}

func NewNBPCOTxRepo(db *sql.DB) *NBPCOTxRepo {
	return &NBPCOTxRepo{db: db}
}

type NBPCOTxRepo struct {
	db *sql.DB
}

func (repo *NBPCOTxRepo) InsertNBPCOTx(ctx context.Context, tx NBPCOTx) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLNBPCOTxInsert,
			tx.Status,
			tx.Mode,
			tx.Amt.Amount,
			tx.Amt.Currency,
			tx.NBPReference,
			tx.CashOutAccId,
			tx.ParentTxId,
			tx.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLNBPCOTxInsert,
			tx.Status,
			tx.Mode,
			tx.Amt.Amount,
			tx.Amt.Currency,
			tx.NBPReference,
			tx.CashOutAccId,
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

func (repo *NBPCOTxRepo) UpdateNBPCOTxById(ctx context.Context, tx NBPCOTx) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(sQLNBPCOTxUpdateById, tx.Status, tx.ID)
	} else {
		_, err = repo.db.Exec(sQLNBPCOTxUpdateById, tx.Status, tx.ID)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *NBPCOTxRepo) GetUniqueNBPCOTxById(ctx context.Context, id int64) (*NBPCOTx, error) {
	rows, err := repo.db.Query(sQLNBPCOTxGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *NBPCOTx

	for rows.Next() {
		f, err = scanRowIntoNBPCOTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *NBPCOTxRepo) GetUniqueNBPCOTxByParentTxId(ctx context.Context, id int64) (*NBPCOTx, error) {
	rows, err := repo.db.Query(sQLNBPCOTxGetUniqueByParentTxId, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *NBPCOTx

	for rows.Next() {
		f, err = scanRowIntoNBPCOTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *NBPCOTxRepo) GetUniqueNBPCOTxByNBPReference(ctx context.Context, nbpReference string) (*NBPCOTx, error) {
	rows, err := repo.db.Query(sQLNBPCOTxGetUniqueByNBPReference, nbpReference)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *NBPCOTx

	for rows.Next() {
		f, err = scanRowIntoNBPCOTx(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoNBPCOTx(rows *sql.Rows) (*NBPCOTx, error) {
	tx := new(NBPCOTx)
	err := rows.Scan(
		&tx.ID,
		&tx.Status,
		&tx.Mode,
		&tx.Amt.Amount,
		&tx.Amt.Currency,
		&tx.NBPReference,
		&tx.CashOutAccId,
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
