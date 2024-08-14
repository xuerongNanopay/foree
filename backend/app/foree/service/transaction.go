package service

import (
	"context"

	"xue.io/go-pay/app/foree/transport"
)

type TransactionService struct {
}

func (t *TransactionService) GetRate(ctx context.Context, req GetRateReq) (*RateDTO, transport.ForeeError) {
	return nil, nil
}
