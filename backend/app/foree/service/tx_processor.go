package service

import (
	"fmt"
	"sync"

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
	processingMap    []map[int64]*transaction.ForeeTx // Avoid duplicate process
	processingLock   sync.RWMutex
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
	maxLoop := 16
	i := 0
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

		if i > maxLoop {
			return nil, fmt.Errorf("unexpect looping for ForeeTx `%v`", nTx.ID)
		}
		i += 1
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
	}
	if tx.Status == transaction.TxStatusComplete || tx.Status == transaction.TxStatusCancel || tx.Status == transaction.TxStatusReject {
		//TODO: log warn.
		return &tx, nil
	}

	switch tx.CurStage {
	case transaction.TxStageInteracCI:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send scotia API
			//Set to Send
		case transaction.TxStatusSent:
			//Check status from scotia API.
		case transaction.TxStatusComplete:
			//Move to next stage
		case transaction.TxStatusReject:
			//Set to reject
		case transaction.TxStatusCancel:
			// set to cancel
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	case transaction.TxStageIDM:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send IDMAPI
			//Set to Send
		case transaction.TxStatusComplete:
			//Move to next stage
		case transaction.TxStatusReject:
			//Set to reject
		case transaction.TxStatusSuspend:
			//Wait to approve
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	case transaction.TxStageNBPCI:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send NBP API
		case transaction.TxStatusSent:
			//Check status from scotia API.
		case transaction.TxStatusComplete:
			// set tx sum to complete
		case transaction.TxStatusReject:
			//TODO: refund
			// set tx sum to cancel
		case transaction.TxStatusCancel:
			// tx.Status =
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	default:
		return nil, fmt.Errorf("transaction `%v` in unknown stage `%s`", tx.ID, tx.CurStage)
	}
	return nil, nil
}
