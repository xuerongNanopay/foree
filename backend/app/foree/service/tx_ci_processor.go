package service

import (
	"context"

	"xue.io/go-pay/app/foree/transaction"
)

type InteracCreditAccount struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type CITxProcessor struct {
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
