package service

import (
	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
)

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
	userRepo         *auth.UserRepo
	contactRepo      *account.ContactAccountRepo
	interacRepo      *account.InteracAccountRepo
}

func (p *TxProcessor) CreateTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *TxProcessor) LoadTx(id int64) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *TxProcessor) ProcessTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return p.doProcessTx(tx)
}

func (p *TxProcessor) doProcessTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if tx.Status == transaction.TxStatusInitial {
		tx.CurStage = transaction.TxStageInteracCI
		tx.CurStageStatus = transaction.TxStatusInitial
		return &tx, nil
	} else {

	}
	return nil, nil
}
