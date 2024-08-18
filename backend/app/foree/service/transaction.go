package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
)

var rateCacheTimeout time.Duration = 15 * time.Minute

type RateCacheItem struct {
	rate   transaction.Rate
	expire time.Time
}

const (
	FeeName string = "FOREE_TX_CAD_FEE"
)

// Group level transaction limit.
var txLimits = map[string]transaction.TxLimit{
	"foree_personal": {
		Name: "foree_personal-group-tx-limit",
		MinAmt: types.AmountData{
			Amount:   types.Amount(10.0),
			Currency: "CAD",
		},
		MaxAmt: types.AmountData{
			Amount:   types.Amount(1000.0),
			Currency: "CAD",
		},
		IsEnable: true,
	},
	"foree_bo": {
		Name: "foree_bo-group-tx-limit",
		MinAmt: types.AmountData{
			Amount:   types.Amount(2.0),
			Currency: "CAD",
		},
		MaxAmt: types.AmountData{
			Amount:   types.Amount(1000.0),
			Currency: "CAD",
		},
		IsEnable: true,
	},
}

type TransactionService struct {
	txSummaryRepo     *transaction.TxSummaryRepo
	txQuoteRepo       *transaction.TxQuoteRepo
	rateRepo          *transaction.RateRepo
	rewardRepo        *transaction.RewardRepo
	dailyTxLimiteRepo *transaction.DailyTxLimitRepo
	feeRepo           *transaction.FeeRepo
	feeJointRepo      *transaction.FeeJointRepo
	txProcessor       *TxProcessor
	rateCache         map[string]RateCacheItem
	rateCacheRWLock   sync.RWMutex
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

func (p *TransactionService) quoteTx(ctx context.Context, user auth.User, quote QuoteTransactionReq) (*transaction.ForeeTx, transport.ForeeError) {
	rate, err := p.rateRepo.GetUniqueRateById(ctx, transaction.GenerateRateId(quote.SrcCurrency, quote.DestCurrency))
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if rate == nil {
		return nil, transport.NewInteralServerError("user `%v` try to create transaction with unkown rate `%s`", user.ID, transaction.GenerateRateId(quote.SrcCurrency, quote.DestCurrency))
	}

	//Reward
	var reward *transaction.Reward
	if len(quote.RewardIds) == 1 {
		rewardId := quote.RewardIds[1]
		r, err := p.rewardRepo.GetUniqueRewardById(ctx, rewardId)
		if err != nil {
			return nil, transport.WrapInteralServerError(err)
		}
		if r == nil {
			return nil, transport.NewInteralServerError("user `%v` try to redeem unknown reward `%v`", user.ID, rewardId)
		}
		if r.OwnerId != user.ID {
			return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` that is belong to `%v`", user.ID, rewardId, r.OwnerId)
		}
		if r.Status != transaction.RewardStatusActive {
			return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` that is currently in status `%v`", user.ID, rewardId, r.Status)
		}
		if r.Amt.Currency != quote.SrcCurrency {
			return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` that apply currency `%v` to currency `%v`", user.ID, rewardId, r.Amt.Currency, quote.SrcCurrency)
		}
		// if (quote.SrcAmount - float64(r.Amt.Amount)) < 10 {
		// 	return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` with srcAmount `%v`", user.ID, rewardId, quote.SrcCurrency))
		// }
		reward = r
	}

	//TODO: PromoCode
	//Don't return err. Just ignore the promocode reward.
	// 	if quote.PromoCode != "" {
	// 		promoCode, err := p.promoCodeRepo.GetUniquePromoCodeByCode(ctx, quote.PromoCode)
	// 		if err != nil {
	// 			//TODO: log
	// 			goto existpromo
	// 		}
	// 		if promoCode == nil {
	// 			//TODO: log
	// 			goto existpromo
	// 		}
	// 		if !promoCode.IsValid() {
	// 			//TODO: log
	// 			goto existpromo
	// 		}
	// 		if quote.SrcCurrency != promoCode.MinAmt.Currency {
	// 			//TODO: log
	// 			goto existpromo
	// 		}
	// 		if quote.SrcAmount < float64(promoCode.MinAmt.Amount) {
	// 			//TODO: log
	// 			goto existpromo
	// 		}
	// 		//TODO: check account limit.
	// 	}
	// existpromo:

	//Fee
	fee, err := p.feeRepo.GetUniqueFeeByName(FeeName)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if fee == nil {
		return nil, transport.NewInteralServerError("fee `%v` not found", FeeName)
	}

	joint, err := fee.MaybeApplyFee(types.AmountData{Amount: types.Amount(quote.SrcAmount), Currency: quote.SrcCurrency})
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if joint != nil {
		joint.Description = fee.Description
		joint.OwnerId = user.ID
	}

	txLimit, ok := txLimits[user.Group]
	if !ok {
		return nil, transport.NewInteralServerError("transaction limit no found for group `%v`", user.Group)
	}

	//TODO: check srcAmount/limit.
	dailyLimit, err := p.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	//Total = quote.srcAmount + fees - rewards
	totalAmt := types.AmountData{
		Amount:   types.Amount(quote.SrcAmount),
		Currency: quote.SrcCurrency,
	}

	if totalAmt.Amount+dailyLimit.UsedAmt.Amount > txLimit.MaxAmt.Amount {
		return nil, transport.NewFormError("Invalid quote transaction request", "srcAmount", fmt.Sprintf("available amount is %v", txLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount))
	}

	if reward != nil {
		totalAmt.Amount -= reward.Amt.Amount
	}

	if totalAmt.Amount < txLimit.MinAmt.Amount {
		return nil, transport.NewFormError("Invalid quote transaction request", "srcAmount", fmt.Sprintf("amount should at lease %v %s without rewards", txLimit.MinAmt.Amount, txLimit.MinAmt.Currency))
	}

	if joint != nil {
		totalAmt.Amount += joint.Amt.Amount
	}

	foreeTx := &transaction.ForeeTx{
		Type:   transaction.TxTypeInteracToNBP,
		Status: transaction.TxStatusInitial,
		Rate:   types.Amount(rate.CalculateForwardAmount(quote.SrcAmount)),
		SrcAmt: types.AmountData{
			Amount:   types.Amount(quote.SrcAmount),
			Currency: quote.SrcCurrency,
		},
		DestAmt: types.AmountData{
			Amount:   types.Amount(rate.CalculateForwardAmount(quote.SrcAmount)),
			Currency: quote.DestCurrency,
		},
		TransactionPurpose: quote.TransactionPurpose,
		CinAccId:           quote.CinAccId,
		CoutAccId:          quote.CoutAccId,
	}

	if joint != nil {
		foreeTx.Fees = []*transaction.FeeJoint{joint}
		foreeTx.TotalFeeAmt = joint.Amt
	}

	if reward != nil {
		foreeTx.RewardIds = quote.RewardIds
		foreeTx.Rewards = []*transaction.Reward{reward}
		foreeTx.TotalRewardAmt = reward.Amt
	}

	foreeTx.TotalAmt = totalAmt

	return foreeTx, nil
}

func (p *TransactionService) rollBackTx(tx transaction.ForeeTx) {
	// log error.
}

func (p *TransactionService) GetTxLimit(user auth.User) (*transaction.TxLimit, error) {
	txLimit, ok := txLimits[user.Group]
	if !ok {
		return nil, transport.NewInteralServerError("transaction limit no found for group `%v`", user.Group)
	}
	return &txLimit, nil
}

func (p *TransactionService) GetDailyTxLimit(user auth.User) (*transaction.DailyTxLimit, error) {
	ctx := context.Background()
	return p.getDailyTxLimit(ctx, user)
}

func (p *TransactionService) addDailyTxLimit(ctx context.Context, user auth.User, amt types.AmountData) (*transaction.DailyTxLimit, error) {
	dailyLimit, err := p.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}

	dailyLimit.UsedAmt.Amount += amt.Amount

	if err := p.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return nil, err
	}
	newDailyLimit, err := p.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}
	return newDailyLimit, nil
}

func (p *TransactionService) minusDailyTxLimit(ctx context.Context, user auth.User, amt types.AmountData) (*transaction.DailyTxLimit, error) {
	dailyLimit, err := p.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}

	dailyLimit.UsedAmt.Amount -= amt.Amount

	if err := p.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return nil, err
	}
	newDailyLimit, err := p.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}
	return newDailyLimit, nil
}

// I don't case race condition here, cause create transaction will save it.
func (p *TransactionService) getDailyTxLimit(ctx context.Context, user auth.User) (*transaction.DailyTxLimit, error) {
	reference := transaction.GenerateDailyTxLimitReference(user.ID)
	dailyLimit, err := p.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
	if err != nil {
		return nil, err
	}

	// If not create one.
	if dailyLimit == nil {
		txLimit, ok := txLimits[user.Group]
		if !ok {
			return nil, fmt.Errorf("transaction limit no found for group `%v`", user.Group)
		}
		dailyLimit = &transaction.DailyTxLimit{
			Reference: reference,
			UsedAmt: types.AmountData{
				Amount:   0.0,
				Currency: txLimit.MaxAmt.Currency,
			},
			MaxAmt: types.AmountData{
				Amount:   txLimit.MaxAmt.Amount,
				Currency: txLimit.MaxAmt.Currency,
			},
		}
		_, err := p.dailyTxLimiteRepo.InsertDailyTxLimit(ctx, *dailyLimit)
		if err != nil {
			return nil, err
		}
		dl, err := p.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
		if err != nil {
			return nil, err
		}
		dailyLimit = dl
	}
	return dailyLimit, nil
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
