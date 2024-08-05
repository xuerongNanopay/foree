package foree_transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type TxStage string

const (
	TxStageCI TxStage = "CASH_IN"
)

type TxAmountData struct {
	Amount  types.AmountData `json:"amount,omitempty"`
	Curreny string           `json:"currency,omitempty"`
}

type ForeeTransaction struct {
}

type CITransaction struct {
}

type IDMTransaction struct {
}

type NBPTransaction struct {
}

type Fee struct {
	ID            string
	TransactionId int64
}

// type
// Src
// Dest
type TransactionSummary struct {
	ID          string
	Summary     string
	SrcAmt      TxAmountData
	DestAmt     TxAmountData
	CiIsoStatus string
	Type        string
	Status      string
	FXRate      string

	Total   TxAmountData
	Created time.Time
}
