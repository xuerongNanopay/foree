package transaction

import "time"

type TransactionStatusHistory struct {
	ID         int64
	ParentTxId int64
	Stage      TxStage
	Status     TxStatus
	ExtraInfo  string
	CreateAt   time.Time `json:"createAt"`
	OwnerId    int64
}
