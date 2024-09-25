package foree_service

import (
	"database/sql"

	"xue.io/go-pay/app/foree/transaction"
)

type TransactionSuperService struct {
	db                 *sql.DB
	foreeTxRepo        *transaction.ForeeTxRepo
	ciTxRepo           *transaction.InteracCITxRepo
	interacTxRepo      *transaction.InteracCITxRepo
	idmTxRepo          *transaction.IdmTxRepo
	npbTxRepo          *transaction.NBPCOTxRepo
	interacTxProcessor *InteracTxProcessor
	txProcessor        *TxProcessor
}

// func (t *TransactionSuperService) ForceCIStatusUpdate(ctx context.Context, fTxId int64, newStatus transaction.TxStage) (*transaction.ForeeTx, error) {
// 	if newStatus != transaction.TxStage(transaction.TxStatusRejected) &&
// 		newStatus != transaction.TxStage(transaction.TxStatusCompleted) {
// 		return nil, fmt.Errorf("ForceCIStatusUpdate -- unacceptable new transaction status %s", newStatus)
// 	}
// 	fTx, err := t.txProcessor.loadTx(fTxId, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Double check to avoid create DB transaction if not necessary.
// 	if fTx.Stage != transaction.TxStageInteracCI && fTx.StageStatus != transaction.TxStatusSent {
// 		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, fTx.StageStatus, fTx.Stage)
// 	}

// 	dTx, err := t.db.Begin()
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: log err
// 		return nil, err
// 	}

// 	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

// 	cfTx, err := t.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: log err
// 		return nil, err
// 	}

// 	// Recheck with DB transaction.
// 	if cfTx.Stage != transaction.TxStageInteracCI && cfTx.StageStatus != transaction.TxStatusSent {
// 		dTx.Rollback()
// 		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, cfTx.StageStatus, cfTx.Stage)
// 	}

// 	fTx.CI.Status = transaction.TxStatus(newStatus)
// 	fTx.StageStatus = transaction.TxStatus(newStatus)

// 	err = t.interacTxRepo.UpdateInteracCITxById(ctx, *fTx.CI)
// 	if err != nil {
// 		dTx.Rollback()
// 		return nil, err
// 	}

// 	err = t.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
// 	if err != nil {
// 		dTx.Rollback()
// 		return nil, err
// 	}

// 	if err = dTx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	go func() {
// 		_, err := t.txProcessor.processTx(*fTx)
// 		if err != nil {
// 			//log: error
// 		}
// 	}()
// 	return fTx, nil
// }

// func (t *TransactionSuperService) IdmStatusUpdate(ctx context.Context, fTxId int64, newStatus transaction.TxStage) (*transaction.ForeeTx, error) {
// 	if newStatus != transaction.TxStage(transaction.TxStatusRejected) &&
// 		newStatus != transaction.TxStage(transaction.TxStatusCompleted) {
// 		return nil, fmt.Errorf("IdmStatusUpdate -- unacceptable new transaction status %s", newStatus)
// 	}

// 	fTx, err := t.txProcessor.loadTx(fTxId, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Double check to avoid create DB transaction if not necessary.
// 	if fTx.Stage != transaction.TxStageIDM && fTx.StageStatus != transaction.TxStatusSuspend {
// 		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, fTx.StageStatus, fTx.Stage)
// 	}

// 	dTx, err := t.db.Begin()
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: log err
// 		return nil, err
// 	}

// 	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

// 	cfTx, err := t.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: log err
// 		return nil, err
// 	}

// 	// Recheck with DB transaction.
// 	if cfTx.Stage != transaction.TxStageIDM && cfTx.StageStatus != transaction.TxStatusSuspend {
// 		dTx.Rollback()
// 		return nil, fmt.Errorf("forceCIStatusUpdate -- transaction `%v` is currently in status `%s` at stage `%s`", fTxId, cfTx.StageStatus, cfTx.Stage)
// 	}

// 	fTx.CI.Status = transaction.TxStatus(newStatus)
// 	fTx.StageStatus = transaction.TxStatus(newStatus)

// 	err = t.interacTxRepo.UpdateInteracCITxById(ctx, *fTx.CI)
// 	if err != nil {
// 		dTx.Rollback()
// 		return nil, err
// 	}

// 	err = t.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
// 	if err != nil {
// 		dTx.Rollback()
// 		return nil, err
// 	}

// 	if err = dTx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	go func() {
// 		_, err := t.txProcessor.processTx(*fTx)
// 		if err != nil {
// 			//log: error
// 		}
// 	}()

// 	return fTx, nil
// }

//TDOD: ForceNBP
