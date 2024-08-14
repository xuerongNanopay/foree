package service

import "xue.io/go-pay/app/foree/transaction"

// It is the internal service for transaction process.

type TxProcessor struct {
	interacTxRepo    *transaction.InteracCITxRepo
	npbTxRepo        *transaction.NBPCOTxRepo
	idmTxRepo        *transaction.IdmTxRepo
	txHistoryRepo    *transaction.TxHistoryRepo
	txSummaryRepo    *transaction.TxSummaryRepo
	txLimitRepo      *transaction.TxLimitRepo
	txLimitCacheRepo *transaction.TxLimitCacheRepo
	foreeTxRepo      *transaction.ForeeTxRepo
}
