package service

import (
	"context"
	"database/sql"

	"xue.io/go-pay/app/foree/transaction"
)

type InteracCreditAccount struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type CITxProcessor struct {
	db *sql.DB
}

func (p *CITxProcessor) requestPayment(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *CITxProcessor) waitPaymentReceive(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *CITxProcessor) doWait(ctx context.Context, tx transaction.ForeeTx) {
	return
}
