package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"xue.io/go-pay/app/foree/account"
	foree_constant "xue.io/go-pay/app/foree/constant"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/partner/nbp"
	"xue.io/go-pay/partner/scotia"
	"xue.io/go-pay/server/transport"
)

func NewTransactionService(
	db *sql.DB,
	authService *AuthService,
	userGroupRepo *auth.UserGroupRepo,
	foreeTxRepo *transaction.ForeeTxRepo,
	txSummaryRepo *transaction.TxSummaryRepo,
	txQuoteRepo *transaction.TxQuoteRepo,
	rewardRepo *transaction.RewardRepo,
	dailyTxLimiteRepo *transaction.DailyTxLimitRepo,
	contactAccountRepo *account.ContactAccountRepo,
	interacAccountRepo *account.InteracAccountRepo,
	feeJointRepo *transaction.FeeJointRepo,
	rateService *RateService,
	feeService *FeeService,
	txProcessor *TxProcessor,
	scotiaClient scotia.ScotiaClient,
	nbpClient nbp.NBPClient,
) *TransactionService {
	return &TransactionService{
		db:                 db,
		authService:        authService,
		userGroupRepo:      userGroupRepo,
		foreeTxRepo:        foreeTxRepo,
		txSummaryRepo:      txSummaryRepo,
		txQuoteRepo:        txQuoteRepo,
		rewardRepo:         rewardRepo,
		dailyTxLimiteRepo:  dailyTxLimiteRepo,
		contactAccountRepo: contactAccountRepo,
		interacAccountRepo: interacAccountRepo,
		feeJointRepo:       feeJointRepo,
		rateService:        rateService,
		feeService:         feeService,
		txProcessor:        txProcessor,
		scotiaClient:       scotiaClient,
		nbpClient:          nbpClient,
	}
}

type TransactionService struct {
	db                 *sql.DB
	authService        *AuthService
	userGroupRepo      *auth.UserGroupRepo
	foreeTxRepo        *transaction.ForeeTxRepo
	txSummaryRepo      *transaction.TxSummaryRepo
	txQuoteRepo        *transaction.TxQuoteRepo
	rewardRepo         *transaction.RewardRepo
	dailyTxLimiteRepo  *transaction.DailyTxLimitRepo
	contactAccountRepo *account.ContactAccountRepo
	interacAccountRepo *account.InteracAccountRepo
	feeJointRepo       *transaction.FeeJointRepo
	rateService        *RateService
	feeService         *FeeService
	txProcessor        *TxProcessor
	scotiaClient       scotia.ScotiaClient
	nbpClient          nbp.NBPClient
}

func (t *TransactionService) GetRate(ctx context.Context, req GetRateReq) (*RateDTO, transport.HError) {
	rate, err := t.rateService.GetRate(req.SrcCurrency, req.DestCurrency)
	if err != nil {
		foree_logger.Logger.Error("GetRate_Fail",
			"ip", loadRealIp(ctx),
			"rateId", transaction.GenerateRateId(req.SrcCurrency, req.DestCurrency),
			"cause", err.Error(),
		)
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
	foree_logger.Logger.Debug("GetRate_Success", "ip", loadRealIp(ctx), "rate", fmt.Sprintf("`%v`", rate))
	return NewRateDTO(rate), nil
}

func (t *TransactionService) FreeQuote(ctx context.Context, req FreeQuoteReq) (*QuoteTransactionDTO, transport.HError) {
	rate, err := t.rateService.GetRate(req.SrcCurrency, req.DestCurrency)
	if err != nil {
		foree_logger.Logger.Error("FreeQuote_Fail",
			"ip", loadRealIp(ctx),
			"rateId", transaction.GenerateRateId(req.SrcCurrency, req.DestCurrency),
			"cause", err.Error(),
		)
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
	feeJoints, err := t.feeService.applyFee(foree_constant.DefaultFeeGroup, types.AmountData{Amount: types.Amount(req.SrcAmount), Currency: req.SrcCurrency})
	if err != nil {
		foree_logger.Logger.Error("FreeQuote_Fail",
			"defaultFeeGroup", foree_constant.DefaultFeeGroup,
			"ip", loadRealIp(ctx),
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	//Total = req.srcAmount + fees - rewards
	totalAmt := types.AmountData{
		Amount:   types.Amount(req.SrcAmount),
		Currency: req.SrcCurrency,
	}

	totalFee := types.AmountData{}
	for _, joint := range feeJoints {
		totalAmt.Amount += joint.Amt.Amount
		totalFee.Amount += joint.Amt.Amount
		totalFee.Currency = joint.Amt.Currency
	}

	//TODO: calculate fee.
	txSum := TxSummaryDetailDTO{
		Summary:       "Free qupte",
		SrcAmount:     types.Amount(req.SrcAmount),
		SrcCurrency:   req.SrcCurrency,
		DestAmount:    types.Amount(rate.CalculateForwardAmount(req.SrcAmount)),
		DestCurrency:  req.DestCurrency,
		FeeAmount:     totalFee.Amount,
		FeeCurrency:   totalFee.Currency,
		TotalAmount:   totalAmt.Amount,
		TotalCurrency: totalAmt.Currency,
	}
	return &QuoteTransactionDTO{
		TxSum: txSum,
	}, nil
}

func (t *TransactionService) QuoteTx(ctx context.Context, req QuoteTransactionReq) (*QuoteTransactionDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("QuoteTx_Fail",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	user := *session.User
	rate, err := t.rateService.GetRate(req.SrcCurrency, req.DestCurrency)
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"rateId", transaction.GenerateRateId(req.SrcCurrency, req.DestCurrency),
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	if rate == nil {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"rateId", transaction.GenerateRateId(req.SrcCurrency, req.DestCurrency),
			"cause", "missing rate",
		)
		return nil, transport.NewInteralServerError("user `%v` try to quote transaction with unkown rate `%s`", user.ID, transaction.GenerateRateId(req.SrcCurrency, req.DestCurrency))
	}

	// Get CI account.
	ciAcc, err := t.interacAccountRepo.GetUniqueActiveInteracAccountByOwnerAndId(ctx, user.ID, req.CinAccId)
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	if ciAcc == nil {
		foree_logger.Logger.Warn("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"ciAccId", req.CinAccId,
			"cause", "Cash in account no found",
		)
		return nil, transport.NewInteralServerError("user `%v` try to use unkown ci account `%v`", user.ID, req.CinAccId)
	}

	// Get Cout account.
	coutAcc, err := t.contactAccountRepo.GetUniqueActiveContactAccountByOwnerAndId(ctx, user.ID, req.CoutAccId)
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	if coutAcc == nil {
		foree_logger.Logger.Warn("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"coutAccId", req.CoutAccId,
			"cause", "Cash out account no found",
		)
		return nil, transport.NewInteralServerError("user `%v` try to use unkown cout account `%v`", user.ID, req.CoutAccId)
	}

	// Get reward
	var reward *transaction.Reward
	if len(req.RewardIds) == 1 {
		rewardId := req.RewardIds[1]
		r, err := t.rewardRepo.GetUniqueRewardById(ctx, rewardId)
		if err != nil {
			foree_logger.Logger.Error("QuoteTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"cause", err.Error(),
			)
			return nil, transport.WrapInteralServerError(err)
		}
		if r == nil {
			foree_logger.Logger.Warn("QuoteTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"rewardId", rewardId,
				"cause", "reward no found",
			)
			return nil, transport.NewInteralServerError("user `%v` try to redeem unknown reward `%v`", user.ID, rewardId)
		}
		if r.OwnerId != user.ID {
			foree_logger.Logger.Warn("QuoteTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"rewardId", rewardId,
				"cause", "user try to use other account's reward",
			)
			return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` that is belong to `%v`", user.ID, rewardId, r.OwnerId)
		}
		if r.Status != transaction.RewardStatusActive {
			foree_logger.Logger.Warn("QuoteTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"rewardId", rewardId,
				"rewardStatus", r.Status,
				"cause", "user try to use non-active reward",
			)
			return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` that is currently in status `%v`", user.ID, rewardId, r.Status)
		}
		if r.Amt.Currency != req.SrcCurrency {
			foree_logger.Logger.Warn("QuoteTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"rewardId", rewardId,
				"requestCurrency", req.SrcCurrency,
				"cause", "user try to use reward that has different currency",
			)
			return nil, transport.NewInteralServerError("user `%v` try to redeem reward `%v` that apply currency `%v` to currency `%v`", user.ID, rewardId, r.Amt.Currency, req.SrcCurrency)
		}
		reward = r
	}

	feeJoints, err := t.feeService.applyFee(session.UserGroup.FeeGroup, types.AmountData{Amount: types.Amount(req.SrcAmount), Currency: req.SrcCurrency})
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	txLimit, ok := foree_constant.TxLimits[foree_constant.TransactionLimitGroup(session.UserGroup.TransactionLimitGroup)]
	if !ok {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"transactionLimitGroup", session.UserGroup.TransactionLimitGroup,
			"cause", "unknown transaction group",
		)
		return nil, transport.NewInteralServerError("transaction limit no found for group `%v`", session.UserGroup.TransactionLimitGroup)
	}

	dailyLimit, err := t.getDailyTxLimit(ctx, *session)
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	//Total = req.srcAmount + fees - rewards
	totalAmt := types.AmountData{
		Amount:   types.Amount(req.SrcAmount),
		Currency: req.SrcCurrency,
	}

	if totalAmt.Amount+dailyLimit.UsedAmt.Amount > txLimit.MaxAmt.Amount {
		foree_logger.Logger.Warn("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"requstAmount", totalAmt.Amount,
			"requstCurrency", totalAmt.Currency,
			"remainingAmount", txLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount,
			"maxAmount", txLimit.MaxAmt.Amount,
			"cause", "overlimit",
		)
		return nil, transport.NewFormError("Invalid req transaction request", "srcAmount", fmt.Sprintf("available amount is %v", txLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount))
	}

	if reward != nil {
		totalAmt.Amount -= reward.Amt.Amount
	}

	if totalAmt.Amount < txLimit.MinAmt.Amount {
		foree_logger.Logger.Warn("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"amoutnAfterReward", totalAmt.Amount,
			"requstCurrency", totalAmt.Currency,
			"minAmount", txLimit.MinAmt.Amount,
			"cause", "underlimit",
		)
		return nil, transport.NewFormError("Invalid req transaction request", "srcAmount", fmt.Sprintf("amount should at lease %v %s without rewards", txLimit.MinAmt.Amount, txLimit.MinAmt.Currency))
	}

	totalFee := types.AmountData{}
	if len(feeJoints) > 0 {
		for _, joint := range feeJoints {
			totalAmt.Amount += joint.Amt.Amount
			totalFee.Amount += joint.Amt.Amount
			totalFee.Currency = joint.Amt.Currency
		}
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
		FeeJoints:          feeJoints,
		TotalFeeAmt:        totalFee,
		TransactionPurpose: req.TransactionPurpose,
		CinAccId:           req.CinAccId,
		CoutAccId:          req.CoutAccId,
		InteracAcc:         ciAcc,
		ContactAcc:         coutAcc,
		OwnerId:            user.ID,
		Owner:              &user,
	}

	if reward != nil {
		foreeTx.RewardIds = req.RewardIds
		foreeTx.Rewards = []*transaction.Reward{reward}
		foreeTx.TotalRewardAmt = reward.Amt
	}

	foreeTx.TotalAmt = totalAmt

	txSummary := transaction.TxSummary{
		Summary: fmt.Sprintf(
			"$%.2f%s -> %.2f%s | %s -> %s",
			foreeTx.SrcAmt.Amount, foreeTx.SrcAmt.Currency,
			foreeTx.DestAmt.Amount, foreeTx.DestAmt.Currency,
			foreeTx.InteracAcc.GetLegalName(),
			foreeTx.ContactAcc.GetLegalName(),
		),
		Type:           string(coutAcc.Type),
		Status:         transaction.TxSummaryStatusInitial,
		Rate:           rate.ToSummary(),
		SrcAccSummary:  foreeTx.InteracAcc.GetLegalName(),
		SrcAmount:      foreeTx.SrcAmt.Amount,
		SrcCurrency:    foreeTx.SrcAmt.Currency,
		DestAccSummary: foreeTx.ContactAcc.GetLegalName(),
		DestAmount:     foreeTx.DestAmt.Amount,
		DestCurrency:   foreeTx.DestAmt.Currency,
		TotalAmount:    foreeTx.TotalAmt.Amount,
		TotalCurrency:  foreeTx.TotalAmt.Currency,
		OwnerId:        user.ID,
	}

	if len(feeJoints) > 0 {
		txSummary.FeeAmount = foreeTx.TotalFeeAmt.Amount
		txSummary.FeeCurrency = foreeTx.TotalFeeAmt.Currency
	}

	if reward != nil {
		txSummary.RewardAmount = foreeTx.TotalRewardAmt.Amount
		txSummary.RewardCurrency = foreeTx.TotalRewardAmt.Currency
	}

	foreeTx.Summary = &txSummary

	quoteId, err := t.txQuoteRepo.InsertTxQuote(ctx, transaction.TxQuote{
		Tx:     foreeTx,
		OwerId: user.ID,
	})

	if err != nil {
		foree_logger.Logger.Error("QuoteTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("QuoteTx_Success", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId)

	return &QuoteTransactionDTO{
		QuoteId: quoteId,
		TxSum:   *NewTxSummaryDetailDTO(&txSummary),
	}, nil
}

func (t *TransactionService) CreateTx(ctx context.Context, req CreateTransactionReq) (*TxSummaryDetailDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("CreateTx_Fail",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	user := session.User
	quote := t.txQuoteRepo.GetUniqueById(ctx, req.QuoteId)

	if quote == nil {
		foree_logger.Logger.Warn("CreateTx_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"quoteId", req.QuoteId,
			"cause", "quote not found",
		)
		return nil, transport.NewInteralServerError("user `%v` try to get null quote", user.ID)
	}

	foreeTx := quote.Tx

	// Start database transaction.
	dTx, err := t.db.Begin()
	if err != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)
	//Lock ci account.
	_, err = t.interacAccountRepo.GetUniqueActiveInteracAccountForUpdateByOwnerAndId(ctx, foreeTx.OwnerId, foreeTx.CinAccId)
	if err != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}

	var wg sync.WaitGroup

	// Recheck and update rewards
	var rewardErr transport.HError
	rewardChecker := func() {
		defer wg.Done()
		if len(foreeTx.Rewards) == 1 {
			reward := foreeTx.Rewards[1]
			reward, err := t.rewardRepo.GetUniqueRewardById(ctx, reward.ID)
			if err != nil {
				foree_logger.Logger.Error("CreateTx_Fail",
					"ip", loadRealIp(ctx),
					"userId", session.UserId,
					"sessionId", req.SessionId,
					"quoteId", req.QuoteId,
					"cause", err.Error(),
				)
				rewardErr = transport.WrapInteralServerError(err)
			}
			if reward.Status != transaction.RewardStatusActive {
				foree_logger.Logger.Warn("CreateTx_Fail",
					"ip", loadRealIp(ctx),
					"userId", session.UserId,
					"sessionId", req.SessionId,
					"quoteId", req.QuoteId,
					"rewardStatus", reward.Status,
					"cause", "reward is not in active",
				)
				rewardErr = transport.NewInteralServerError("user `%v` try to create a transaction with reward `%v` in status `%s`", user.ID, reward.ID, reward.Status)
			}
		}
	}
	wg.Add(1)
	go rewardChecker()

	// Recheck and update limit
	var limitErr transport.HError
	limitChecker := func() {
		defer wg.Done()
		txLimit := foree_constant.TxLimits[foree_constant.TransactionLimitGroup(session.UserGroup.TransactionLimitGroup)]

		dailyLimit, err := t.getDailyTxLimit(ctx, *session)
		if err != nil {
			foree_logger.Logger.Error("CreateTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"quoteId", req.QuoteId,
				"cause", err.Error(),
			)
			limitErr = transport.WrapInteralServerError(err)
		}

		if foreeTx.SrcAmt.Amount+dailyLimit.UsedAmt.Amount > txLimit.MaxAmt.Amount {
			foree_logger.Logger.Warn("QuoteTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"requstAmount", foreeTx.SrcAmt.Amount,
				"requstCurrency", foreeTx.SrcAmt.Currency,
				"remainingAmount", txLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount,
				"maxAmount", txLimit.MaxAmt.Amount,
				"cause", "overlimit",
			)
			limitErr = transport.NewInteralServerError("user `%v` try to create a transaction with `%v` but the remaining limit is `%v`", user.ID, foreeTx.SrcAmt.Amount, txLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount)
		}

		if _, err := t.addDailyTxLimit(ctx, *session, foreeTx.SrcAmt); err != nil {
			dTx.Rollback()
			foree_logger.Logger.Error("CreateTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"quoteId", req.QuoteId,
				"cause", err.Error(),
			)
			limitErr = transport.WrapInteralServerError(err)
		}
	}

	wg.Add(1)
	go limitChecker()

	// Create foree transaction.
	var foreeTxErr transport.HError
	var foreeTxID int64
	createForeeTx := func() {
		defer wg.Done()
		id, err := t.foreeTxRepo.InsertForeeTx(ctx, *foreeTx)
		if err != nil {
			dTx.Rollback()
			foree_logger.Logger.Error("CreateTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"quoteId", req.QuoteId,
				"cause", err.Error(),
			)
			foreeTxErr = transport.WrapInteralServerError(err)
		}
		foreeTxID = id
	}

	wg.Add(1)
	go createForeeTx()

	wg.Wait()
	if limitErr != nil {
		dTx.Rollback()
		return nil, limitErr
	}
	if rewardErr != nil {
		dTx.Rollback()
		return nil, rewardErr
	}
	if foreeTxErr != nil {
		dTx.Rollback()
		return nil, foreeTxErr
	}

	// Below code can do in other coroutine.
	foreeTx.ID = foreeTxID
	// Create TxSummary, feeJoin, update reward, update limit.
	wg = sync.WaitGroup{}

	// Create TxSummary
	var txSummaryErr transport.HError
	createSummaryTx := func() {
		defer wg.Done()
		foreeTx.Summary.ParentTxId = foreeTxID
		foreeTx.Summary.NBPReference = transaction.GenerateNbpId(foree_constant.DefaultNBPIdPrefix, foreeTxID)
		id, err := t.txSummaryRepo.InsertTxSummary(ctx, *foreeTx.Summary)
		if err != nil {
			foree_logger.Logger.Error("CreateTx_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"quoteId", req.QuoteId,
				"cause", err.Error(),
			)
			txSummaryErr = transport.WrapInteralServerError(err)
		}
		foreeTx.Summary.ID = id
	}
	wg.Add(1)
	go createSummaryTx()

	//TODO: update to patch insert.
	// Create feeJoint
	var feeJointError transport.HError
	createFeeJoint := func() {
		defer wg.Done()
		for _, feeJoin := range foreeTx.FeeJoints {
			feeJoin.ParentTxId = foreeTxID
			_, err := t.feeJointRepo.InsertFeeJoint(ctx, *feeJoin)
			if err != nil {
				foree_logger.Logger.Error("CreateTx_Fail",
					"ip", loadRealIp(ctx),
					"userId", session.UserId,
					"sessionId", req.SessionId,
					"quoteId", req.QuoteId,
					"cause", err.Error(),
				)
				feeJointError = transport.WrapInteralServerError(err)
				return
			}
		}
	}
	wg.Add(1)
	go createFeeJoint()

	// Update reward
	var rewardError transport.HError
	updateReward := func() {
		defer wg.Done()
		if len(foreeTx.Rewards) == 1 {
			r := foreeTx.Rewards[1]
			r.AppliedTransactionId = foreeTxID
			r.Status = transaction.RewardStatusPending
			err := t.rewardRepo.UpdateRewardTxById(ctx, *r)
			if err != nil {
				foree_logger.Logger.Error("CreateTx_Fail",
					"ip", loadRealIp(ctx),
					"userId", session.UserId,
					"sessionId", req.SessionId,
					"quoteId", req.QuoteId,
					"cause", err.Error(),
				)
				rewardError = transport.WrapInteralServerError(err)
			}
		}
	}

	wg.Add(1)
	go updateReward()

	wg.Wait()
	if txSummaryErr != nil {
		dTx.Rollback()
		return nil, txSummaryErr
	}
	if feeJointError != nil {
		dTx.Rollback()
		return nil, feeJointError
	}
	if rewardError != nil {
		dTx.Rollback()
		return nil, rewardError
	}

	if err = dTx.Commit(); err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	go t.txProcessor.createAndProcessTx(*foreeTx)

	foree_logger.Logger.Info("CreateTx_Success", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "foreeTxId", foreeTxID)
	return NewTxSummaryDetailDTO(foreeTx.Summary), nil
}

// func (t *TransactionService) GetTxLimit(user auth.User) (*transaction.TxLimit, error) {
// 	txLimit, ok := foree_constant.TxLimits[userGroup]
// 	if !ok {
// 		return nil, transport.NewInteralServerError("transaction limit no found for group `%v`", user.Group)
// 	}
// 	return &txLimit, nil
// }

func (t *TransactionService) GetDailyTxLimit(ctx context.Context, req transport.SessionReq) (*DailyTxLimitDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("CreateTx_Fail",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	limit, err := t.getDailyTxLimit(ctx, *session)
	if err != nil {
		foree_logger.Logger.Error("GetDailyTxLimit_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	foree_logger.Logger.Debug("GetDailyTxLimit_Success",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
	)
	return NewDailyTxLimitDTO(limit), nil
}

func (t *TransactionService) addDailyTxLimit(ctx context.Context, session auth.Session, amt types.AmountData) (*transaction.DailyTxLimit, error) {
	dailyLimit, err := t.getDailyTxLimit(ctx, session)
	if err != nil {
		return nil, err
	}

	dailyLimit.UsedAmt.Amount += amt.Amount

	if err := t.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return nil, err
	}

	return dailyLimit, nil
}

// I don't case race condition here, cause create transaction will rescure it.
func (t *TransactionService) getDailyTxLimit(ctx context.Context, session auth.Session) (*transaction.DailyTxLimit, error) {
	reference := transaction.GenerateDailyTxLimitReference(session.UserId)
	dailyLimit, err := t.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
	if err != nil {
		return nil, err
	}

	// If not create one.
	if dailyLimit == nil {
		txLimit, ok := foree_constant.TxLimits[foree_constant.TransactionLimitGroup(session.UserGroup.TransactionLimitGroup)]
		if !ok {
			return nil, fmt.Errorf("transaction limit no found for group `%v`", session.UserGroup.TransactionLimitGroup)
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
			OwnerId: session.UserId,
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

func (t *TransactionService) GetTxSummary(ctx context.Context, req GetTransactionReq) (*TxSummaryDetailDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("GetTxSummary_Fail",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	summaryTx, err := t.txSummaryRepo.GetUniqueTxSummaryByOwnerAndId(ctx, session.UserId, req.TransactionId)
	if err != nil {
		foree_logger.Logger.Error("GetTxSummary_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	if summaryTx == nil {
		return nil, nil
	}

	foree_logger.Logger.Debug("GetTxSummary_Success",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
	)
	return NewTxSummaryDetailDTO(summaryTx), nil
}

func (t *TransactionService) QuerySummaryTxs(ctx context.Context, req QueryTransactionReq) ([]*TxSummaryDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("QuerySummaryTxs_Fail",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}
	//TODO: limit, offset pruning
	var summaryTxs []*transaction.TxSummary
	var err error

	if req.Status == "" {
		summaryTxs, err = t.txSummaryRepo.GetAllTxSummaryByOwnerIdWithPagination(ctx, session.UserId, req.Limit, req.Offset)
	} else {
		summaryTxs, err = t.txSummaryRepo.QueryTxSummaryByOwnerIdAndStatusWithPagination(ctx, session.UserId, req.Status, req.Limit, req.Offset)
	}

	if err != nil {
		foree_logger.Logger.Error("QuerySummaryTxs_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	rets := make([]*TxSummaryDTO, len(summaryTxs))

	for i, v := range summaryTxs {
		rets[i] = NewTxSummaryDTO(v)
	}

	foree_logger.Logger.Debug("QuerySummaryTxs_Success",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
	)
	return rets, nil
}

// Check transaction status, see if is able to cancel.
func (t *TransactionService) CancelTransaction(ctx context.Context, req CancelTransactionReq) (*TxCancelDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("CancelTransaction_Fail",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	summaryTx, err := t.txSummaryRepo.GetUniqueTxSummaryByOwnerAndId(ctx, session.UserId, req.TransactionId)
	if err != nil {
		foree_logger.Logger.Error("CancelTransaction_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	if summaryTx == nil {
		return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "no found")
	}

	fTx, err := t.txProcessor.LoadTx(summaryTx.ParentTxId)
	if err != nil {
		foree_logger.Logger.Error("CancelTransaction_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	if fTx.CurStage == transaction.TxStageInteracCI && fTx.CurStageStatus == transaction.TxStatusSent {
		resp, err := t.scotiaClient.CancelPayment(scotia.CancelPaymentRequest{
			PaymentId:    fTx.CI.ScotiaPaymentId,
			CancelReason: req.CancelReason,
		})
		//TODO: log
		if err != nil {
			foree_logger.Logger.Error("CancelTransaction_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"cause", err.Error(),
			)
			return nil, transport.WrapInteralServerError(err)
		} else if resp.StatusCode/100 != 2 {
			return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "transaction can not cancel")
		}
	} else if fTx.CurStage == transaction.TxStageNBPCO && fTx.CurStageStatus == transaction.TxStatusSent && fTx.COUT.CashOutAcc.Type == foree_constant.ContactAccountTypeCash {
		resp, err := t.nbpClient.CancelTransaction(nbp.CancelTransactionRequest{
			GlobalId:           fTx.COUT.NBPReference,
			CancellationReason: req.CancelReason,
		})
		//TODO: log
		if err != nil {
			foree_logger.Logger.Error("CancelTransaction_Fail",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"cause", err.Error(),
			)
			return nil, transport.WrapInteralServerError(err)
		} else if resp.StatusCode/100 != 2 {
			return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "transaction can not cancel")
		}
	} else {
		foree_logger.Logger.Warn("CancelTransaction_Fail",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"foreeTxId", fTx.ID,
			"CurStage", fTx.CurStage,
			"CurStageStatus", fTx.CurStageStatus,
		)
		return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "transaction can not cancel")
	}

	foree_logger.Logger.Info("CancelTransaction_Success",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
		"foreeTxId", fTx.ID,
	)
	return &TxCancelDTO{
		TransactionId: req.TransactionId,
		Message:       "cancel successfully",
	}, nil
}
