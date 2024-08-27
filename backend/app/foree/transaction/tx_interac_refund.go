package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
)

const (
	sQLInteracRefundTxInsert = `
		INSERT INTO interac_refund_tx
		(
			status, refund_interac_acc_id, parent_tx_id, owner_id
		) VALUES(?,?,?,?)
	`
	sQLInteracRefundTxUpdate = `
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
	ID                 int64                   `json:"id"`
	Status             RefundTxStatus          `json:"status"`
	RefundInteracAccId int64                   `json:"refundInteracAccId"`
	RefundInteracAcc   *account.InteracAccount `json:"refundInteracAcc"`
	ParentTxId         int64                   `json:"parentTxId"`
	OwnerId            int64                   `json:"ownerId"`
	CreateAt           time.Time               `json:"createAt"`
	UpdateAt           time.Time               `json:"updateAt"`
}

func NewForeeRefundTxRepo(db *sql.DB) *ForeeRefundTxRepo {
	return &ForeeRefundTxRepo{db: db}
}

type ForeeRefundTxRepo struct {
	db *sql.DB
}

// func scanRowIntoInteracCITx(rows *sql.Rows) (*InteracCITx, error) {
// 	tx := new(InteracCITx)
// 	err := rows.Scan(
// 		&tx.ID,
// 		&tx.Status,
// 		&tx.CashInAccId,
// 		&tx.DestInteracAccId,
// 		&tx.Amt.Amount,
// 		&tx.Amt.Currency,
// 		&tx.APIReference,
// 		&tx.Url,
// 		&tx.ParentTxId,
// 		&tx.OwnerId,
// 		&tx.CreateAt,
// 		&tx.UpdateAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tx, nil
// }
