package transaction

import "time"

type TransactionStatusHistory struct {
	ID                  int64
	ParentTransactionId int64
	Stage               TxStage
	Status              TxStatus
	ExtraInfo           string
	CreateAt            time.Time `json:"createAt"`
	OwnerId             int64
}
