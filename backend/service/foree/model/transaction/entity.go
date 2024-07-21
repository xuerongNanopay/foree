package transaction

import "time"

type TransactionStage string
type TransactionStatus string

type Transaction struct {
	id     uint64
	userId uint64

	srcCur  string
	srcAmt  uint64
	destCur string
	destAmt uint64

	stage  TransactionStage
	status TransactionStatus

	createAt time.Timer
	updateAt time.Timer
}

// type TransactionStatusHistory struct {
// 	transactionId uint64

// 	stage  TransactionStage
// 	status TransactionStatus
// }
