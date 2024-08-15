package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

const (
	sQLInteracCITxInsert = `
        INSERT INTO interact_ci_tx
        (
            status, src_interac_acc_id, dest_interac_acc_id,
            api_reference, amount, currency, parent_tx_id, owner_id
        ) VALUES(?,?,?,?,?,?,?)
    `
	sQLInteracCITxGetUniqueById = `
        SELECT 
            t.id, t.status, t.src_interac_acc_id, t.dest_interac_acc_id,
            t.amount, t.currency, t.api_reference, t.url
            t.parent_tx_id, t.owner_id, t.create_at, t.update_at
        FROM interact_ci_tx t
        where t.id = ?

    `
	sQLInteracCITxGetUniqueByParentTxId = `
        SELECT 
            t.id, t.status, t.src_interac_acc_id, t.dest_interac_acc_id,
            t.amount, t.currency, t.api_reference, t.url
            t.parent_tx_id, t.owner_id, t.create_at, t.update_at
        FROM interact_ci_tx t
        where t.parent_tx_id = ?
    `
	sQLInteracCITxUpdateById = `
        UPDATE interact_ci_tx SET 
            status = ?, api_reference = ?, url = ?
        WHERE id = ?
    `
)

type InteracCITx struct {
	ID               int64                   `json:"id"`
	Status           TxStatus                `json:"status"`
	APIReference     string                  `json:"apiReference"`
	Url              string                  `json:"url"`
	SrcInteracAccId  int64                   `json:"srcInteracAccId"`
	SrcInteracAcc    *account.InteracAccount `json:"srcInteracAcc"`
	DestInteracAccId int64                   `json:"destInteracAccId"`
	DestInteracAcc   *account.InteracAccount `json:"DestInteracAcc"`
	Amt              types.AmountData        `json:"Amt"`
	ParentTxId       int64                   `json:"parentTxId"`
	OwnerId          int64                   `json:"OwnerId"`
	CreateAt         time.Time               `json:"createAt"`
	UpdateAt         time.Time               `json:"updateAt"`
}

func NewInteracCITxRepo(db *sql.DB) *InteracCITxRepo {
	return &InteracCITxRepo{db: db}
}

type InteracCITxRepo struct {
	db *sql.DB
}

func (repo *InteracCITxRepo) InsertInteracCITx(ctx context.Context, tx InteracCITx) (int64, error) {
	result, err := repo.db.Exec(
		sQLInteracCITxInsert,
		tx.Status,
		tx.SrcInteracAccId,
		tx.DestInteracAccId,
		tx.APIReference,
		tx.Amt.Amount,
		tx.Amt.Curreny,
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

func (repo *InteracCITxRepo) UpdateInteracCITxById(ctx context.Context, tx InteracCITx) error {
	_, err := repo.db.Exec(sQLInteracCITxUpdateById, tx.Status, tx.APIReference, tx.Url, tx.ID)
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

	if f.ID == 0 {
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

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoInteracCITx(rows *sql.Rows) (*InteracCITx, error) {
	tx := new(InteracCITx)
	err := rows.Scan(
		&tx.ID,
		&tx.Status,
		&tx.SrcInteracAccId,
		&tx.DestInteracAccId,
		&tx.Amt.Amount,
		&tx.Amt.Curreny,
		&tx.APIReference,
		&tx.Url,
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
