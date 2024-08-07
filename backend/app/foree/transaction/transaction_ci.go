package transaction

import (
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
			amount, currency, parent_tx_id, owner_id
		) VALUES()
	`
	sQLInteracCITxGetUniqueByParentTxId = `
		SELECT 
			t.id, t.status, t.src_interac_acc_id, t.dest_interac_acc_id,
			t.amount, t.currency, t.api_reference, t.url
			t.parent_tx_id, t.owner_id, t.create_at, t.update_at
		FROM interact_ci_tx t
		where t.parent_tx_id = ?
	`
	sQLInteracCITxGetUniqueById = `
		SELECT 
			t.id, t.status, t.src_interac_acc_id, t.dest_interac_acc_id,
			t.amount, t.currency, t.api_reference, t.url
			t.parent_tx_id, t.owner_id, t.create_at, t.update_at
		FROM interact_ci_tx t
		where t.id = ?
	
	`
	sQLInteracCITxUpdateById = `
		UPDATE interact_ci_tx SET 
			status = ?, api_reference = ?, url = ?
		WHERE id = ?
	`
)

type InteracCITx struct {
	ID               int64
	Status           TxStatus
	APIReference     string
	Url              string
	SrcInteracAccId  int64
	SrcInteracAcc    *account.InteracAccount
	DestInteracAccId int64
	DestInteracAcc   *account.InteracAccount
	Amt              types.AmountData
	ParentTxId       int64
	OwnerId          int64
	CreateAt         time.Time `json:"createAt"`
	UpdateAt         time.Time `json:"updateAt"`
}

func NewInteracCIRepo(db *sql.DB) *InteracCIRepo {
	return &InteracCIRepo{db: db}
}

type InteracCIRepo struct {
	db *sql.DB
}

func (repo *InteracCIRepo) InsertInteracCITx(tx InteracCITx) (int64, error) {
	result, err := repo.db.Exec(
		sQLInteracCITxInsert,
		tx.Status,
		tx.SrcInteracAccId,
		tx.DestInteracAccId,
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

func (repo *InteracCIRepo) GetUniqueInteracCITxByParentTxId(parentTxId int64) (*InteracCITx, error) {
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

func (repo *InteracCIRepo) GetUniqueInteracCITxById(id int64) (*InteracCITx, error) {
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

func (repo *InteracCIRepo) UpdateInteracCITxById(tx InteracCITx) error {
	_, err := repo.db.Exec(sQLInteracCITxUpdateById, tx.Status, tx.APIReference, tx.Url, tx.ID)
	if err != nil {
		return err
	}
	return nil
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
