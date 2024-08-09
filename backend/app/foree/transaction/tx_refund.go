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
	ParentTxId         int64
	OwnerId            int64
	RefundType         string
	Status             string
	IsRefunded         bool
	RefundAt           time.Time
	RefundInteracAccId int64
	RefundInteracAcc   *account.InteracAccount
	CreateAt           time.Time `json:"createAt"`
	UpdateAt           time.Time `json:"updateAt"`
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
// 		&tx.Amt.Curreny,
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
