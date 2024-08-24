package service

import (
	"database/sql"

	"xue.io/go-pay/app/foree/transaction"
)

type TransactionSuperService struct {
	db          *sql.DB
	foreeTxRepo *transaction.ForeeTxRepo
	idmTxRepo   *transaction.IdmTxRepo
}

//TODO: ForceCI
//TODO: ApproveIDM
//TODO: RejectIDM
//TDOD: ForceNBP
