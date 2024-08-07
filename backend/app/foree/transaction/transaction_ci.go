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
			t.amount, t.currency, t.parent_tx_id, t.owner_id
		FROM interact_ci_tx t
		where t.parent_tx_id = ?
	`
	sQLInteracCITxGetUniqueById = `
		SELECT 
			t.id, t.status, t.src_interac_acc_id, t.dest_interac_acc_id,
			t.amount, t.currency, t.parent_tx_id, t.owner_id
		FROM interact_ci_tx t
		where t.id = ?
	
	`
	sQLInteracCITxUpdateById = `
		UPDATE interact_ci_tx SET 
			status = ?, scotia_id = ?, url = ?
		WHERE id = ?
	`
)

type ScotiaInteracCITx struct {
	ID               int64
	Status           TxStatus
	ScotialId        string
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

func (repo *InteracCIRepo) InsertReferral(referal Referral) (int64, error) {
	result, err := repo.db.Exec(
		sQLReferralInsert,
		referal.Code,
		referal.ReferralType,
		referal.ReferralValue,
		referal.Status,
		referal.ReferrerId,
		referal.IsRedeemed,
		referal.ExpireAt,
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
