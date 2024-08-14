package service

import (
	"fmt"

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
	var err error
	var nTx *transaction.ForeeTx
	for {
		nTx, err = p.doProcessTx(tx)
		if err != nil {
			return nil, err
		}
		if tx.CurStage == nTx.CurStage && nTx.CurStageStatus == tx.CurStageStatus {
			return nTx, nil
		}
		// Record the history.
		go p.recordTxHistory(transaction.NewTxHistory(nTx, ""))
		tx = *nTx
	}

}

func (p *TxProcessor) recordTxHistory(h *transaction.TxHistory) {
	if _, err := p.txHistoryRepo.InserTxHistory(*h); err != nil {
		fmt.Println(err.Error())
	}

}

func (p *TxProcessor) doProcessTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if tx.Status == transaction.TxStatusInitial {
		tx.Status = transaction.TxStatusProcessing
		tx.CurStage = transaction.TxStageInteracCI
		tx.CurStageStatus = transaction.TxStatusInitial
		return &tx, nil
	} else {
		//TODO:
		return nil, nil
	}
	return nil, nil
}
