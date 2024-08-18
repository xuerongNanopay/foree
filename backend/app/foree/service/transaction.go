package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/app/foree/types"
)

var rateCacheTimeout time.Duration = 15 * time.Minute

type RateCacheItem struct {
	rate   transaction.Rate
	expire time.Time
}

type TransactionService struct {
	txSummaryRepo   *transaction.TxSummaryRepo
	txQuoteRepo     *transaction.TxQuoteRepo
	rateRepo        *transaction.RateRepo
	txProcessor     *TxProcessor
	rateCache       map[string]RateCacheItem
	rateCacheRWLock sync.RWMutex
}

// Can be cache for 5 minutes.
func (t *TransactionService) GetRate(ctx context.Context, req GetRateReq) (*RateDTO, transport.ForeeError) {
	rate, err := t.getRate(ctx, req.SrcCurrency, req.DestCurrency)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if rate == nil {
		return nil, transport.NewFormError(
			"Invalid rate request",
			"srcCurrency",
			fmt.Sprintf("unsupport srcCurrency %s", req.SrcCurrency),
			"destCurrency",
			fmt.Sprintf("unsupport destCurrency %s", req.DestCurrency),
		)
	}
	return NewRateDTO(rate), nil
}

// Only case this cache won't work is that volume density of request is high.
func (t *TransactionService) getRate(ctx context.Context, src, dest string) (*transaction.Rate, error) {
	rateId := transaction.GenerateRateId(src, dest)

	t.rateCacheRWLock.RLock()
	rateCache, ok := t.rateCache[rateId]
	t.rateCacheRWLock.RUnlock()

	if ok && rateCache.expire.After(time.Now()) {
		return &rateCache.rate, nil
	}

	rate, err := t.rateRepo.GetUniqueRateById(ctx, rateId)
	if err != nil {
		return nil, err
	}

	if !t.rateCacheRWLock.TryLock() {
		return rate, nil
	}
	defer t.rateCacheRWLock.Unlock()

	t.rateCache[rateId] = RateCacheItem{
		rate:   *rate,
		expire: time.Now().Add(rateCacheTimeout),
	}
	return rate, nil
}

// Can be use same cache as above.
// Do we want it? Or we can calculate at frontend.
func (t *TransactionService) FreeQuote(ctx context.Context, req FreeQuoteReq) (*TxSummaryDetailDTO, transport.ForeeError) {
	rate, err := t.getRate(ctx, req.SrcCurrency, req.DestCurrency)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	if rate == nil {
		return nil, transport.NewFormError(
			"Invalid rate request",
			"srcCurrency",
			fmt.Sprintf("unsupport srcCurrency %s", req.SrcCurrency),
			"destCurrency",
			fmt.Sprintf("unsupport destCurrency %s", req.DestCurrency),
		)
	}

	//TODO: calculate fee.
	sumTx := &TxSummaryDetailDTO{
		Summary:      "Free qupte",
		SrcAmount:    types.Amount(req.SrcAmount),
		SrcCurrency:  req.SrcCurrency,
		DestAmount:   types.Amount(rate.CalculateForwardAmount(req.SrcAmount)),
		DestCurrency: req.DestCurrency,
	}
	return sumTx, nil
}

func (t *TransactionService) QuoteTx(ctx context.Context, req QuoteTransactionReq) (*TxSummaryDetailDTO, transport.ForeeError) {
	return nil, nil
}

func (t *TransactionService) ConfirmQuote(ctx context.Context, req ConfirmQuoteReq) (*TxSummaryDetailDTO, transport.ForeeError) {
	return nil, nil
}

func (t *TransactionService) GetTransaction(ctx context.Context, req GetTransactionReq) (*TxSummaryDetailDTO, transport.ForeeError) {
	return nil, nil
}

func (t *TransactionService) GetAllTransactions(ctx context.Context, req GetAllTransactionReq) ([]*TxSummaryDTO, transport.ForeeError) {
	return nil, nil
}

func (t *TransactionService) QueryTransactions(ctx context.Context, req QueryTransactionReq) ([]*TxSummaryDTO, transport.ForeeError) {
	return nil, nil
}
