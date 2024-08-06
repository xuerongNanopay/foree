package transaction

import "time"

type TransactionStatusHistory struct {
	ID                  int64
	ParentTransactionId int64
	Stage               string
	Status              string
	CreateAt            time.Time `json:"createAt"`
	OwnerId             int64
}
