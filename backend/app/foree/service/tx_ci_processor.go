package service

import (
	"context"
	"database/sql"
	"fmt"

	"xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/partner/scotia"
)

type InteracCreditAccount struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type CITxProcessor struct {
	scotiaClient *scotia.ScotiaClient
	foreeTxRepo  *transaction.ForeeTxRepo
	db           *sql.DB
}

func (p *CITxProcessor) requestPayment(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	// Lock transaction and safety check.
	nForeeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, tx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	if nForeeTx.CurStage != transaction.TxStageInteracCI && nForeeTx.CurStageStatus != transaction.TxStatusInitial {
		return nil, fmt.Errorf("transaction `%v` is in status `%s` at stage `%s`", nForeeTx.ID, nForeeTx.CurStageStatus, nForeeTx.Status)
	}

	// API call
	//
	return nil, nil
}

func (p *CITxProcessor) createRequestPaymentReq(ctx context.Context, tx transaction.ForeeTx) {

}

func (p *CITxProcessor) waitPaymentReceive(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *CITxProcessor) doWait(ctx context.Context, tx transaction.ForeeTx) {
	return
}
