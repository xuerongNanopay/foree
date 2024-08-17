package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/account"
)

//DO this later.
// const (
// 	RefundType = "MANUAL_REFUND"
// )

type ForeeRefundTx struct {
	ID                 int64                   `json:"id"`
	ParentTxId         int64                   `json:"parentTxId"`
	OwnerId            int64                   `json:"ownerId"`
	RefundType         string                  `json:"refundType"`
	Status             string                  `json:"status"`
	IsRefunded         bool                    `json:"isRefunded"`
	RefundAt           time.Time               `json:"refundAt"`
	RefundInteracAccId int64                   `json:"refundInteracAccId"`
	RefundInteracAcc   *account.InteracAccount `json:"refundInteracAcc"`
	CreateAt           time.Time               `json:"createAt"`
	UpdateAt           time.Time               `json:"updateAt"`
}

// func NewForeeRefundTxRepo(db *sql.DB) *ForeeRefundTxRepo {
// 	return &ForeeRefundTxRepo{db: db}
// }

// type ForeeRefundTxRepo struct {
// 	db *sql.DB
// }

// func scanRowIntoInteracCITx(rows *sql.Rows) (*InteracCITx, error) {
// 	tx := new(InteracCITx)
// 	err := rows.Scan(
// 		&tx.ID,
// 		&tx.Status,
// 		&tx.SrcInteracAccId,
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
