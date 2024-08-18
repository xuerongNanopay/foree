package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transaction"
	time_util "xue.io/go-pay/util/time"
)

type TxProcessorConfig struct {
}

// It is the internal service for transaction process.

type TxProcessor struct {
	interacTxRepo  *transaction.InteracCITxRepo
	npbTxRepo      *transaction.NBPCOTxRepo
	idmTxRepo      *transaction.IdmTxRepo
	txHistoryRepo  *transaction.TxHistoryRepo
	txSummaryRepo  *transaction.TxSummaryRepo
	foreeTxRepo    *transaction.ForeeTxRepo
	contactRepo    *account.ContactAccountRepo
	interacRepo    *account.InteracAccountRepo
	processingMap  []map[int64]*transaction.ForeeTx // Avoid duplicate process
	processingLock sync.RWMutex
}

func (p *TxProcessor) createTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *TxProcessor) InsertTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *TxProcessor) loadTx(id int64) (*transaction.ForeeTx, error) {
	ctx := context.Background()
	foree, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load CI
	ci, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, foree.ID)
	if err != nil {
		return nil, err
	}
	if ci == nil {
		return nil, fmt.Errorf("InteracCITx no found for ForeeTx `%v`", foree.ID)
	}

	srcInteracAcc, err := p.interacRepo.GetUniqueInteracAccountById(ctx, ci.SrcInteracAccId)
	if err != nil {
		return nil, err
	}
	if srcInteracAcc == nil {
		return nil, fmt.Errorf("SrcInteracAcc no found for InteracCITx `%v`", ci.SrcInteracAccId)
	}
	ci.SrcInteracAcc = srcInteracAcc

	destInteracAcc, err := p.interacRepo.GetUniqueInteracAccountById(ctx, ci.DestInteracAccId)
	if err != nil {
		return nil, err
	}
	if destInteracAcc == nil {
		return nil, fmt.Errorf("DestInteracAcc no found for InteracCITx `%v`", ci.DestInteracAccId)
	}
	ci.DestInteracAcc = destInteracAcc
	foree.CI = ci

	// Load IDM
	idm, err := p.idmTxRepo.GetUniqueIDMTxByParentTxId(ctx, foree.ID)
	if err != nil {
		return nil, err
	}
	if idm == nil {
		return nil, fmt.Errorf("IDMTx no found for ForeeTx `%v`", foree.ID)
	}
	foree.IDM = idm

	// Load COUT
	cout, err := p.npbTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, foree.ID)
	if err != nil {
		return nil, err
	}
	if cout == nil {
		return nil, fmt.Errorf("NBPCOTx no found for ForeeTx `%v`", foree.ID)
	}

	destContactAcc, err := p.contactRepo.GetUniqueContactAccountById(ctx, cout.DestContactAccId)
	if err != nil {
		return nil, err
	}
	if destContactAcc == nil {
		return nil, fmt.Errorf("DestContactAcc no found for NBPCOTx `%v`", cout.DestContactAccId)
	}
	cout.DestContactAcc = destContactAcc
	foree.COUT = cout

	// TODO: fees?, rewards?

	return foree, nil
}

// TODO: change argument to int64
func (p *TxProcessor) processTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if tx.Type != transaction.TxTypeInteracToNBP {
		return nil, fmt.Errorf("unknow ForeeTx type `%s`", tx.Type)
	}
	var err error
	var nTx *transaction.ForeeTx
	maxLoop := 16
	i := 0
	ctx := context.Background()
	for {
		nTx, err = p.doProcessTx(ctx, tx)
		if err != nil {
			return nil, err
		}
		if tx.CurStage == nTx.CurStage && nTx.CurStageStatus == tx.CurStageStatus {
			return nTx, nil
		}
		// Record the history.
		go p.recordTxHistory(*transaction.NewTxHistory(nTx, ""))
		tx = *nTx

		if i > maxLoop {
			return nil, fmt.Errorf("unexpect looping for ForeeTx `%v`", nTx.ID)
		}
		i += 1
	}

}

func (p *TxProcessor) recordTxHistory(h transaction.TxHistory) {
	if _, err := p.txHistoryRepo.InserTxHistory(h); err != nil {
		fmt.Println(err.Error())
	}

}

func (p *TxProcessor) doProcessTx(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if tx.Status == transaction.TxStatusInitial {
		tx.Status = transaction.TxStatusProcessing
		tx.CurStage = transaction.TxStageInteracCI
		tx.CurStageStatus = transaction.TxStatusInitial
		return &tx, nil
	}
	if tx.Status == transaction.TxStatusCompleted || tx.Status == transaction.TxStatusCancelled || tx.Status == transaction.TxStatusRejected {
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
			//Check status from scotia API.(Webhook, or cron)
			//Just do noting waiting for cron
		case transaction.TxStatusCompleted:
			tx.CurStage = transaction.TxStageInteracCI
			tx.CurStageStatus = transaction.TxStatusInitial
			return &tx, nil
		case transaction.TxStatusRejected:
			//Set to reject
		case transaction.TxStatusCancelled:
			// set to cancel
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	case transaction.TxStageIDM:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send IDMAPI
			//Set to Send
		case transaction.TxStatusCompleted:
			//Move to next stage
			tx.CurStage = transaction.TxStageNBPCO
			tx.CurStageStatus = transaction.TxStatusInitial
			return &tx, nil
		case transaction.TxStatusRejected:
			// Set to ForeeTx to terminal status.
			tx.Status = transaction.TxStatusRejected
			tx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", tx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			// Close remaing non-terminated transactions.
			nT, err := p.closeRemainingTx(ctx, tx)
			if err != nil {
				return nil, err
			}
			go p.maybeRefund(*nT)
			return nT, nil
		case transaction.TxStatusSuspend:
			//Wait to approve
			//Log warn?
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	case transaction.TxStageNBPCO:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send NBP API
		case transaction.TxStatusSent:
			//Check status from NBP API.
			//Or just wait for clone
		case transaction.TxStatusCompleted:
			tx.Status = transaction.TxStatusCompleted
			tx.Conclusion = fmt.Sprintf("Complete at %s.", time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			return &tx, nil
			// set tx sum to complete
		case transaction.TxStatusRejected:
			tx.Status = transaction.TxStatusRejected
			tx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", tx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			go p.maybeRefund(tx)
			return &tx, nil
		case transaction.TxStatusCancelled:
			tx.Status = transaction.TxStatusCancelled
			tx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", tx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			go p.maybeRefund(tx)
			return &tx, nil
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	default:
		return nil, fmt.Errorf("transaction `%v` in unknown stage `%s`", tx.ID, tx.CurStage)
	}
	return nil, nil
}

func (p *TxProcessor) closeRemainingTx(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	switch tx.CurStage {
	case transaction.TxStageInteracCI:
		idm := tx.IDM
		co := tx.COUT
		idm.Status = transaction.TxStatusClosed
		co.Status = transaction.TxStatusClosed
		if err := p.idmTxRepo.UpdateIDMTxById(ctx, *idm); err != nil {
			return nil, err
		}
		if err := p.npbTxRepo.UpdateNBPCOTxById(ctx, *co); err != nil {
			return nil, err
		}
		return &tx, nil
	case transaction.TxStageIDM:
		co := tx.COUT
		co.Status = transaction.TxStatusClosed
		if err := p.npbTxRepo.UpdateNBPCOTxById(ctx, *co); err != nil {
			return nil, err
		}
		return &tx, nil
	default:
		//TODO: Log warn
		return &tx, nil
	}
}

func (p *TxProcessor) maybeRefund(tx transaction.ForeeTx) {
	//TODO: implement
}

// TODO: change argement to id.
func (p *TxProcessor) approveIDM(ctx context.Context, tx transaction.ForeeTx) {
	if tx.CurStage == transaction.TxStageIDM && tx.CurStageStatus == transaction.TxStatusSuspend {

	}
	//TODO: implement
}

// TODO: change argement to id.
func (p *TxProcessor) rejectIDM(ctx context.Context, tx transaction.ForeeTx) {
	if tx.CurStage == transaction.TxStageIDM && tx.CurStageStatus == transaction.TxStatusSuspend {

	}
}
