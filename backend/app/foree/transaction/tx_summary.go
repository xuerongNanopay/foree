package transaction

import "time"

type TxSummary struct {
	ID              int64     `json:"id"`
	Summary         string    `json:"sumary"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Rate            string    `json:"rate"`
	SrcAmount       string    `json:"srcAmount"`
	SrcCurrency     string    `json:"srcCurrency"`
	DestAmount      string    `json:"destAmount"`
	DestCurrency    string    `json:"desCurrency"`
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	ParentTxId      int64     `json:"parentTxd"`
	OwnerId         int64     `json:"owerId"`
	CreateAt        time.Time `json:"createAt"`
	UpdateAt        time.Time `json:"updateAt"`
}
