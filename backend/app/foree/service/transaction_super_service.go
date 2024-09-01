package service

import (
	"context"
	"database/sql"
	"fmt"

	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/constant"
)

type TransactionSuperService struct {
	db            *sql.DB
	foreeTxRepo   *transaction.ForeeTxRepo
	ciTxRepo      *transaction.InteracCITxRepo
	interacTxRepo *transaction.InteracCITxRepo
	idmTxRepo     *transaction.IdmTxRepo
	npbTxRepo     *transaction.NBPCOTxRepo
	ciTxProcessor *CITxProcessor
	txProcessor   *TxProcessor
}

func (t *TransactionSuperService) ForceCIStatusUpdate(ctx context.Context, fTxId int64, newStatus transaction.TxStage) (*transaction.ForeeTx, error) {
	if newStatus != transaction.TxStage(transaction.TxStatusRejected) &&
		newStatus != transaction.TxStage(transaction.TxStatusCompleted) {
		return nil, fmt.Errorf("ForceCIStatusUpdate -- unacceptable new transaction status %s", newStatus)
	}
	fTx, err := t.txProcessor.loadTx(fTxId, true)
	if err != nil {
		return nil, err
	}

	// Double check to avoid create DB transaction if not necessary.
	if fTx.CurStage != transaction.TxStageInteracCI && fTx.CurStageStatus != transaction.TxStatusSent {
		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, fTx.CurStageStatus, fTx.CurStage)
	}

	dTx, err := t.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	cfTx, err := t.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	// Recheck with DB transaction.
	if cfTx.CurStage != transaction.TxStageInteracCI && cfTx.CurStageStatus != transaction.TxStatusSent {
		dTx.Rollback()
		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, cfTx.CurStageStatus, cfTx.CurStage)
	}

	fTx.CI.Status = transaction.TxStatus(newStatus)
	fTx.CurStageStatus = transaction.TxStatus(newStatus)

	err = t.interacTxRepo.UpdateInteracCITxById(ctx, *fTx.CI)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	err = t.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if err = dTx.Commit(); err != nil {
		return nil, err
	}

	go func() {
		_, err := t.txProcessor.processTx(*fTx)
		if err != nil {
			//log: error
		}
	}()
	return fTx, nil
}

func (t *TransactionSuperService) IdmStatusUpdate(ctx context.Context, fTxId int64, newStatus transaction.TxStage) (*transaction.ForeeTx, error) {
	if newStatus != transaction.TxStage(transaction.TxStatusRejected) &&
		newStatus != transaction.TxStage(transaction.TxStatusCompleted) {
		return nil, fmt.Errorf("IdmStatusUpdate -- unacceptable new transaction status %s", newStatus)
	}

	fTx, err := t.txProcessor.loadTx(fTxId, true)
	if err != nil {
		return nil, err
	}

	// Double check to avoid create DB transaction if not necessary.
	if fTx.CurStage != transaction.TxStageIDM && fTx.CurStageStatus != transaction.TxStatusSuspend {
		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, fTx.CurStageStatus, fTx.CurStage)
	}

	dTx, err := t.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	cfTx, err := t.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	// Recheck with DB transaction.
	if cfTx.CurStage != transaction.TxStageIDM && cfTx.CurStageStatus != transaction.TxStatusSuspend {
		dTx.Rollback()
		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, cfTx.CurStageStatus, cfTx.CurStage)
	}

	fTx.CI.Status = transaction.TxStatus(newStatus)
	fTx.CurStageStatus = transaction.TxStatus(newStatus)

	err = t.interacTxRepo.UpdateInteracCITxById(ctx, *fTx.CI)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	err = t.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if err = dTx.Commit(); err != nil {
		return nil, err
	}

	go func() {
		_, err := t.txProcessor.processTx(*fTx)
		if err != nil {
			//log: error
		}
	}()

	return fTx, nil
}

//TDOD: ForceNBP
