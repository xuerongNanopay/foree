package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/types"
)

const (
	sQLNBPCOTxInsert = `
		INSERT INTO nbp_co_tx
		(
			status, amount, currency,  dest_contact_acc_id
			parent_tx_id, owner_id
		) VALUES(?,?,?,?,?,?)
	`
	sQLNBPCOTxGetUniqueById = `
	SELECT 
		t.id, t.status, t.amount, t.currency,
		t.dest_contact_acc_id, t.parent_tx_id, t.owner_id,
		t.create_at, t.update_at
	FROM nbp_co_tx t
	where t.id = ?

`
	sQLNBPCOTxGetUniqueByParentTxId = `
		SELECT 
			t.id, t.status, t.amount, t.currency,
			t.dest_contact_acc_id, t.parent_tx_id, t.owner_id,
			t.create_at, t.update_at
		FROM nbp_co_tx t
		where t.parent_tx_id = ?
	`
	sQLNBPCOTxUpdateById = `
		UPDATE nbp_co_tx SET 
			status = ?
		WHERE id = ?
	`
)

type NBPCOTx struct {
	ID               int64
	Status           TxStatus
	Amt              types.AmountData
	DestContactAccId int64
	DestContactAcc   *account.ContactAccount
	ParentTxId       int64
	OwnerId          int64
	CreateAt         time.Time `json:"createAt"`
	UpdateAt         time.Time `json:"updateAt"`
}

func NewNBPCORepo(db *sql.DB) *NBPCORepo {
	return &NBPCORepo{db: db}
}

type NBPCORepo struct {
	db *sql.DB
}
