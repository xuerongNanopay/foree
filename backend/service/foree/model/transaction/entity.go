package transaction

import (
	"fmt"
	"time"
)

type TransactionStage string
type TransactionStatus string

const (
	INTERAC_CI TransactionStage = "INTERAC_CI"
	IDM        TransactionStage = "IDM"
	NBP_CO     TransactionStage = "NBP_CO"
)

const (
	PENDING TransactionStatus = "PENDING"
	SENT    TransactionStatus = "SENT"
	SUSPEND TransactionStatus = "SUSPEND"
	DECLINE TransactionStatus = "DECLINE"
	CANCEL  TransactionStatus = "CANCEL"
)

type Transaction struct {
	id     uint64
	userId uint64

	srcCur  string
	srcAmt  uint64
	destCur string
	destAmt uint64

	stage  TransactionStage
	status TransactionStatus

	clientIp    string
	clientAgent string

	createAt time.Timer
	updateAt time.Timer
}

type IdmPayload struct {
	transactionId uint64

	payload string

	createAt time.Timer
	updateAt time.Timer
}

type TransactionStatusHistory struct {
	transactionId uint64

	stage  TransactionStage
	status TransactionStatus

	createAt time.Timer
	updateAt time.Timer
}

func (t *Transaction) processTx() (*Transaction, error) {
	switch t.stage {
	case INTERAC_CI:
		return t.processIntracCi()
	case IDM:
		return t.processIDM()
	case NBP_CO:
		return t.processNbpCo()
	default:
		return t, fmt.Errorf("Unknown transaction stage `%v`", t.stage)
	}
}

func (t *Transaction) processIntracCi() (*Transaction, error) {
	switch t.status {
	default:
		return t, fmt.Errorf("Unknown transaction status `%v` in stage `%v`", t.status, t.stage)
	}
}

func (t *Transaction) processIDM() (*Transaction, error) {
	return t, nil
}

func (t *Transaction) processNbpCo() (*Transaction, error) {
	return t, nil
}
