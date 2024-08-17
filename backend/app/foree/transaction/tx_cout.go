package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

const (
	sQLNBPCOTxInsert = `
        INSERT INTO nbp_co_tx
        (
            status, amount, currency, api_reference,  dest_contact_acc_id
            parent_tx_id, owner_id
        ) VALUES(?,?,?,?,?,?,?)
    `
	sQLNBPCOTxUpdateById = `
        UPDATE nbp_co_tx SET 
            status = ?
        WHERE id = ?
    `
	sQLNBPCOTxGetUniqueById = `
        SELECT 
            t.id, t.status, t.amount, t.currency, t.api_reference,
            t.dest_contact_acc_id, t.parent_tx_id, t.owner_id,
            t.create_at, t.update_at
        FROM nbp_co_tx t
        where t.id = ?

    `
	sQLNBPCOTxGetUniqueByParentTxId = `
        SELECT 
            t.id, t.status, t.amount, t.currency, t.api_reference,
            t.dest_contact_acc_id, t.parent_tx_id, t.owner_id,
            t.create_at, t.update_at
        FROM nbp_co_tx t
        where t.parent_tx_id = ?
    `
)

type NBPCOTx struct {
	ID               int64                   `json:"id"`
	Status           TxStatus                `json:"status"`
	Amt              types.AmountData        `json:"amt"`
	APIReference     string                  `json:"apiReference"`
	DestContactAccId int64                   `json:"destContactAccId"`
	DestContactAcc   *account.ContactAccount `json:"destContactAcc"`
	ParentTxId       int64                   `json:"parentTxId"`
	OwnerId          int64                   `json:"OwnerId"`
	CreateAt         time.Time               `json:"createAt"`
	UpdateAt         time.Time               `json:"updateAt"`
}

func NewNBPCOTxRepo(db *sql.DB) *NBPCOTxRepo {
	return &NBPCOTxRepo{db: db}
}

type NBPCOTxRepo struct {
	db *sql.DB
}

func (repo *NBPCOTxRepo) InsertNBPCOTx(ctx context.Context, tx NBPCOTx) (int64, error) {
	result, err := repo.db.Exec(
		sQLNBPCOTxInsert,
		tx.Status,
		tx.Amt.Amount,
		tx.Amt.Currency,
		tx.APIReference,
		tx.DestContactAccId,
		tx.ParentTxId,
		tx.OwnerId,
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

func (repo *NBPCOTxRepo) UpdateNBPCOTxById(ctx context.Context, tx NBPCOTx) error {
	_, err := repo.db.Exec(sQLNBPCOTxUpdateById, tx.Status, tx.ID)
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

	if f.ID == 0 {
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

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoNBPCOTx(rows *sql.Rows) (*NBPCOTx, error) {
	tx := new(NBPCOTx)
	err := rows.Scan(
		&tx.ID,
		&tx.Status,
		&tx.Amt.Amount,
		&tx.Amt.Currency,
		&tx.APIReference,
		&tx.DestContactAccId,
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
