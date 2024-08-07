package transaction

import "time"

type TransactionSummary struct {
	ID              int64     `json:"id"`
	Summary         string    `json:"sumary"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Rate            string    `json:"rate"`
	SrcId           int64     `json:"srcId"`
	DescId          int64     `json:"descId"`
	SrcAmt          string    `json:"srcAmt"`
	DescAmt         string    `json:"descAmt"`
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	CreateAt        time.Time `json:"createAt"`
	ParentTxId      int64
	OwnerId         int64
}
