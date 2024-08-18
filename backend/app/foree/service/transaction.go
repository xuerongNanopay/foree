package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
)

var (
	rateCacheTimeout time.Duration = 15 * time.Minute
	feeCacheTimeout  time.Duration = time.Hour
)

type CacheItem[T any] struct {
	item     T
	createAt time.Time
}

const (
	FeeName           string = "FOREE_TX_CAD_FEE"
	DefaultForeeGroup string = "FOREE_PERSONAL"
)

// Group level transaction limit.
var txLimits = map[string]transaction.TxLimit{
	"FOREE_PERSONAL": {
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
	"FOREE_BO": {
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
	authService       *AuthService
	txSummaryRepo     *transaction.TxSummaryRepo
	txQuoteRepo       *transaction.TxQuoteRepo
	rateRepo          *transaction.RateRepo
	rewardRepo        *transaction.RewardRepo
	dailyTxLimiteRepo *transaction.DailyTxLimitRepo
	feeRepo           *transaction.FeeRepo
	contactRepo       *account.ContactAccountRepo
	interacRepo       *account.InteracAccountRepo
	feeJointRepo      *transaction.FeeJointRepo
	txProcessor       *TxProcessor
	rateCache         map[string]CacheItem[transaction.Rate]
	rateCacheRWLock   sync.RWMutex
	feeCache          map[string]CacheItem[transaction.Fee]
	feeCacheRWLock    sync.RWMutex
}

// Can be cache for 5 minutes.
func (t *TransactionService) GetRate(ctx context.Context, req GetRateReq) (*RateDTO, transport.ForeeError) {
	rate, err := t.getRate(ctx, req.SrcCurrency, req.DestCurrency, 30*time.Minute)
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

func (t *TransactionService) getRate(ctx context.Context, src, dest string, validIn time.Duration) (*transaction.Rate, error) {
	rateId := transaction.GenerateRateId(src, dest)

	t.rateCacheRWLock.RLock()
	rateCache, ok := t.rateCache[rateId]
	t.rateCacheRWLock.RUnlock()

	if ok && rateCache.createAt.Add(validIn).After(time.Now()) {
		return &rateCache.item, nil
	}

	rate, err := t.rateRepo.GetUniqueRateById(ctx, rateId)
	if err != nil {
		return nil, err
	}

	//There is a change that write lock never work.
	//But if this case happen, we already a big company.
	if !t.rateCacheRWLock.TryLock() {
		return rate, nil
	}
	defer t.rateCacheRWLock.Unlock()

	t.rateCache[rateId] = CacheItem[transaction.Rate]{
		item:     *rate,
		createAt: time.Now(),
	}
	return rate, nil
}

func (t *TransactionService) FreeQuote(ctx context.Context, req FreeQuoteReq) (*QuoteTransactionDTO, transport.ForeeError) {
	rate, err := t.getRate(ctx, req.SrcCurrency, req.DestCurrency, 30*time.Minute)
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

	//fee
	fee, err := t.getFee(ctx, FeeName, 2*time.Hour)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if fee == nil {
		return nil, transport.NewInteralServerError("fee `%v` not found", FeeName)
	}

	joint, err := fee.MaybeApplyFee(types.AmountData{Amount: types.Amount(req.SrcAmount), Currency: req.SrcCurrency})
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	//Total = req.srcAmount + fees - rewards
	totalAmt := types.AmountData{
		Amount:   types.Amount(req.SrcAmount),
		Currency: req.SrcCurrency,
	}

	if joint != nil {
		totalAmt.Amount += joint.Amt.Amount
	}

	//TODO: calculate fee.
	txSum := TxSummaryDetailDTO{
		Summary:       "Free qupte",
		SrcAmount:     types.Amount(req.SrcAmount),
		SrcCurrency:   req.SrcCurrency,
		DestAmount:    types.Amount(rate.CalculateForwardAmount(req.SrcAmount)),
		DestCurrency:  req.DestCurrency,
		FeeAmount:     joint.Amt.Amount,
		FeeCurrency:   joint.Amt.Currency,
		TotalAmount:   totalAmt.Amount,
		TotalCurrency: totalAmt.Currency,
	}
	return &QuoteTransactionDTO{
		TxSum: txSum,
	}, nil
}

func (t *TransactionService) getFee(ctx context.Context, feeName string, validIn time.Duration) (*transaction.Fee, error) {
	t.feeCacheRWLock.RLock()
	feeCache, ok := t.feeCache[feeName]
	t.feeCacheRWLock.RUnlock()

	if ok && feeCache.createAt.Add(validIn).After(time.Now()) {
		return &feeCache.item, nil
	}

	fee, err := t.feeRepo.GetUniqueFeeByName(ctx, feeName)
	if err != nil {
		return nil, err
	}

	//There is a change that write lock never work.
	//But if this case happen, we already a big company.
	if !t.feeCacheRWLock.TryLock() {
		return fee, nil
	}
	defer t.feeCacheRWLock.Unlock()

	t.feeCache[feeName] = CacheItem[transaction.Fee]{
		item:     *fee,
		createAt: time.Now(),
	}
	return fee, nil
}

func (t *TransactionService) QuoteTx(ctx context.Context, req QuoteTransactionReq) (*QuoteTransactionDTO, transport.ForeeError) {
	session, serr := t.authService.VerifySession(ctx, req.SessionId)
	if serr != nil {
		return nil, serr
	}

	user := *session.User
	rate, err := t.getRate(ctx, req.SrcCurrency, req.DestCurrency, 5*time.Minute)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if rate == nil {
		return nil, transport.NewInteralServerError("user `%v` try to quote transaction with unkown rate `%s`", user.ID, transaction.GenerateRateId(req.SrcCurrency, req.DestCurrency))
	}

	// Get CI account.
	ciAcc, err := t.interacRepo.GetUniqueActiveInteracAccountByOwnerAndId(ctx, user.ID, req.CinAccId)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if ciAcc == nil {
		return nil, transport.NewInteralServerError("user `%v` try to use unkown ci account `%v`", user.ID, req.CinAccId)
	}

	// Get Cout account.
	coutAcc, err := t.contactRepo.GetUniqueActiveContactAccountByOwnerAndId()

	// Get reward
	var reward *transaction.Reward
	if len(req.RewardIds) == 1 {
		rewardId := req.RewardIds[1]
		r, err := t.rewardRepo.GetUniqueRewardById(ctx, rewardId)
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
		if r.Amt.Currency != req.SrcCurrency {
			return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` that apply currency `%v` to currency `%v`", user.ID, rewardId, r.Amt.Currency, req.SrcCurrency)
		}
		// if (req.SrcAmount - float64(r.Amt.Amount)) < 10 {
		// 	return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` with srcAmount `%v`", user.ID, rewardId, req.SrcCurrency))
		// }
		reward = r
	}

	//TODO: PromoCode
	//Don't return err. Just ignore the promocode reward.
	// 	if req.PromoCode != "" {
	// 		promoCode, err := t.promoCodeRepo.GetUniquePromoCodeByCode(ctx, req.PromoCode)
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
	// 		if req.SrcCurrency != promoCode.MinAmt.Currency {
	// 			//TODO: log
	// 			goto existpromo
	// 		}
	// 		if req.SrcAmount < float64(promoCode.MinAmt.Amount) {
	// 			//TODO: log
	// 			goto existpromo
	// 		}
	// 		//TODO: check account limit.
	// 	}
	// existpromo:

	//Fee
	fee, err := t.getFee(ctx, FeeName, time.Hour)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if fee == nil {
		return nil, transport.NewInteralServerError("fee `%v` not found", FeeName)
	}

	joint, err := fee.MaybeApplyFee(types.AmountData{Amount: types.Amount(req.SrcAmount), Currency: req.SrcCurrency})
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
	dailyLimit, err := t.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	//Total = req.srcAmount + fees - rewards
	totalAmt := types.AmountData{
		Amount:   types.Amount(req.SrcAmount),
		Currency: req.SrcCurrency,
	}

	if totalAmt.Amount+dailyLimit.UsedAmt.Amount > txLimit.MaxAmt.Amount {
		return nil, transport.NewFormError("Invalid req transaction request", "srcAmount", fmt.Sprintf("available amount is %v", txLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount))
	}

	if reward != nil {
		totalAmt.Amount -= reward.Amt.Amount
	}

	if totalAmt.Amount < txLimit.MinAmt.Amount {
		return nil, transport.NewFormError("Invalid req transaction request", "srcAmount", fmt.Sprintf("amount should at lease %v %s without rewards", txLimit.MinAmt.Amount, txLimit.MinAmt.Currency))
	}

	if joint != nil {
		totalAmt.Amount += joint.Amt.Amount
	}

	foreeTx := &transaction.ForeeTx{
		Type:   transaction.TxTypeInteracToNBP,
		Status: transaction.TxStatusInitial,
		Rate:   types.Amount(rate.CalculateForwardAmount(req.SrcAmount)),
		SrcAmt: types.AmountData{
			Amount:   types.Amount(req.SrcAmount),
			Currency: req.SrcCurrency,
		},
		DestAmt: types.AmountData{
			Amount:   types.Amount(rate.CalculateForwardAmount(req.SrcAmount)),
			Currency: req.DestCurrency,
		},
		TransactionPurpose: req.TransactionPurpose,
		CinAccId:           req.CinAccId,
		CoutAccId:          req.CoutAccId,
	}

	if joint != nil {
		foreeTx.Fees = []*transaction.FeeJoint{joint}
		foreeTx.TotalFeeAmt = joint.Amt
	}

	if reward != nil {
		foreeTx.RewardIds = req.RewardIds
		foreeTx.Rewards = []*transaction.Reward{reward}
		foreeTx.TotalRewardAmt = reward.Amt
	}

	foreeTx.TotalAmt = totalAmt

	txSum := &TxSummaryDetailDTO{
		Summary:       "Free qupte",
		SrcAmount:     types.Amount(req.SrcAmount),
		SrcCurrency:   req.SrcCurrency,
		DestAmount:    types.Amount(rate.CalculateForwardAmount(req.SrcAmount)),
		DestCurrency:  req.DestCurrency,
		FeeAmount:     joint.Amt.Amount,
		FeeCurrency:   joint.Amt.Currency,
		TotalAmount:   totalAmt.Amount,
		TotalCurrency: totalAmt.Currency,
	}
	return &QuoteTransactionDTO{
		QuoteId: "TODO",
		TxSum:   *txSum,
	}, nil
}

func (t *TransactionService) rollBackTx(tx transaction.ForeeTx) {
	// log error.
}

func (t *TransactionService) GetTxLimit(user auth.User) (*transaction.TxLimit, error) {
	txLimit, ok := txLimits[user.Group]
	if !ok {
		return nil, transport.NewInteralServerError("transaction limit no found for group `%v`", user.Group)
	}
	return &txLimit, nil
}

func (t *TransactionService) GetDailyTxLimit(user auth.User) (*transaction.DailyTxLimit, error) {
	ctx := context.Background()
	return t.getDailyTxLimit(ctx, user)
}

func (t *TransactionService) addDailyTxLimit(ctx context.Context, user auth.User, amt types.AmountData) (*transaction.DailyTxLimit, error) {
	dailyLimit, err := t.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}

	dailyLimit.UsedAmt.Amount += amt.Amount

	if err := t.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return nil, err
	}
	newDailyLimit, err := t.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}
	return newDailyLimit, nil
}

func (t *TransactionService) minusDailyTxLimit(ctx context.Context, user auth.User, amt types.AmountData) (*transaction.DailyTxLimit, error) {
	dailyLimit, err := t.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}

	dailyLimit.UsedAmt.Amount -= amt.Amount

	if err := t.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return nil, err
	}
	newDailyLimit, err := t.getDailyTxLimit(ctx, user)
	if err != nil {
		return nil, err
	}
	return newDailyLimit, nil
}

// I don't case race condition here, cause create transaction will save it.
func (t *TransactionService) getDailyTxLimit(ctx context.Context, user auth.User) (*transaction.DailyTxLimit, error) {
	reference := transaction.GenerateDailyTxLimitReference(user.ID)
	dailyLimit, err := t.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
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
		_, err := t.dailyTxLimiteRepo.InsertDailyTxLimit(ctx, *dailyLimit)
		if err != nil {
			return nil, err
		}
		dl, err := t.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
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
