package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"xue.io/go-pay/app/foree/account"
	foree_constant "xue.io/go-pay/app/foree/constant"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
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
	rewardRepo *promotion.RewardRepo,
	contactAccountRepo *account.ContactAccountRepo,
	interacAccountRepo *account.InteracAccountRepo,
	feeJointRepo *transaction.FeeJointRepo,
	rateService *RateService,
	feeService *FeeService,
	txLimitService *TxLimitService,
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
		contactAccountRepo: contactAccountRepo,
		interacAccountRepo: interacAccountRepo,
		feeJointRepo:       feeJointRepo,
		rateService:        rateService,
		feeService:         feeService,
		txLimitService:     txLimitService,
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
	rewardRepo         *promotion.RewardRepo
	contactAccountRepo *account.ContactAccountRepo
	interacAccountRepo *account.InteracAccountRepo
	feeJointRepo       *transaction.FeeJointRepo
	rateService        *RateService
	feeService         *FeeService
	txLimitService     *TxLimitService
	txProcessor        *TxProcessor
	scotiaClient       scotia.ScotiaClient
	nbpClient          nbp.NBPClient
}

func (t *TransactionService) GetRate(ctx context.Context, req GetRateReq) (*RateDTO, transport.HError) {
	rate, err := t.rateService.GetRate(req.SrcCurrency, req.DestCurrency)
	if err != nil {
		foree_logger.Logger.Error("GetRate_FAIL",
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
	foree_logger.Logger.Debug("GetRate_SUCCESS", "ip", loadRealIp(ctx), "rate", fmt.Sprintf("`%v`", rate))
	return NewRateDTO(rate), nil
}

func (t *TransactionService) FreeQuote(ctx context.Context, req FreeQuoteReq) (*QuoteTransactionDTO, transport.HError) {
	rate, err := t.rateService.GetRate(req.SrcCurrency, req.DestCurrency)
	if err != nil {
		foree_logger.Logger.Error("FreeQuote_FAIL",
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
		foree_logger.Logger.Error("FreeQuote_FAIL",
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
		foree_logger.Logger.Info("QuoteTx_FAIL",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	user := *session.User
	rate, err := t.rateService.GetRate(req.SrcCurrency, req.DestCurrency)
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"rateId", transaction.GenerateRateId(req.SrcCurrency, req.DestCurrency),
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	if rate == nil {
		foree_logger.Logger.Error("QuoteTx_FAIL",
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
		foree_logger.Logger.Error("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	if ciAcc == nil {
		foree_logger.Logger.Warn("QuoteTx_FAIL",
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
		foree_logger.Logger.Error("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	if coutAcc == nil {
		foree_logger.Logger.Warn("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"coutAccId", req.CoutAccId,
			"cause", "Cash out account no found",
		)
		return nil, transport.NewInteralServerError("user `%v` try to use unkown cout account `%v`", user.ID, req.CoutAccId)
	}

	// Get reward
	var rewards = make([]*promotion.Reward, 0)

	if len(req.RewardIds) > 0 {
		rs, err := t.rewardRepo.GetAllRewardByOwnerIdAndIds(ctx, session.UserId, req.RewardIds)
		if err != nil {
			foree_logger.Logger.Error("QuoteTx_FAIL",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"cause", err.Error(),
			)
			return nil, transport.WrapInteralServerError(err)
		}

		for _, reward := range rs {
			if reward.Status != promotion.RewardStatusActive {
				foree_logger.Logger.Warn("QuoteTx_FAIL",
					"ip", loadRealIp(ctx),
					"userId", session.UserId,
					"sessionId", req.SessionId,
					"rewardId", reward.ID,
					"rewardStatus", reward.Status,
					"cause", "invalid reward status",
				)
				continue
			}
			if reward.Amt.Currency != req.SrcCurrency {
				foree_logger.Logger.Warn("QuoteTx_FAIL",
					"ip", loadRealIp(ctx),
					"userId", session.UserId,
					"sessionId", req.SessionId,
					"rewardId", reward.ID,
					"requestCurrency", req.SrcCurrency,
					"cause", "user try to use reward that has different currency",
				)
				continue
			}
			rewards = append(rewards, reward)
		}
	}

	feeJoints, err := t.feeService.applyFee(session.UserGroup.FeeGroup, types.AmountData{Amount: types.Amount(req.SrcAmount), Currency: req.SrcCurrency})
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	txLimit, err := t.txLimitService.getTxLimit(ctx, session.UserGroup.TransactionLimitGroup)
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	if txLimit == nil {
		foree_logger.Logger.Error("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"transactionLimitGroup", session.UserGroup.TransactionLimitGroup,
			"cause", "unknown transaction group",
		)
		return nil, transport.NewInteralServerError("transaction limit no found for group `%v`", session.UserGroup.TransactionLimitGroup)
	}

	dailyLimit, err := t.txLimitService.getDailyTxLimit(ctx, *session)
	if err != nil {
		foree_logger.Logger.Error("QuoteTx_FAIL",
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

	if totalAmt.Amount+dailyLimit.UsedAmt.Amount > dailyLimit.MaxAmt.Amount {
		foree_logger.Logger.Warn("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"requstAmount", totalAmt.Amount,
			"requstCurrency", totalAmt.Currency,
			"remainingAmount", dailyLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount,
			"maxAmount", dailyLimit.MaxAmt.Amount,
			"cause", "overlimit",
		)
		return nil, transport.NewFormError("Invalid req transaction request", "srcAmount", fmt.Sprintf("available amount is %v", dailyLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount))
	}

	for _, reward := range rewards {
		totalAmt.Amount -= reward.Amt.Amount
	}

	if totalAmt.Amount < txLimit.MinAmt.Amount {
		foree_logger.Logger.Warn("QuoteTx_FAIL",
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
		Type:           transaction.TxTypeInteracToNBP,
		Stage:          transaction.TxStageBegin,
		LimitReference: dailyLimit.Reference,
		Rate:           types.Amount(rate.CalculateForwardAmount(req.SrcAmount)),
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

	if len(rewards) > 0 {
		rIds := make([]int64, len(rewards))
		rs := make([]*promotion.Reward, len(rewards))
		tRewardAmt := types.AmountData{}

		for i, reward := range rewards {
			rIds[i] = reward.ID
			rs[i] = reward
			tRewardAmt = types.AmountData{
				Amount:   tRewardAmt.Amount + reward.Amt.Amount,
				Currency: reward.Amt.Currency,
			}
		}
		foreeTx.RewardIds = req.RewardIds
		foreeTx.Rewards = rs
		foreeTx.TotalRewardAmt = tRewardAmt
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
		SrcAccId:       ciAcc.ID,
		SrcAccSummary:  foreeTx.InteracAcc.GetLegalName(),
		SrcAmount:      foreeTx.SrcAmt.Amount,
		SrcCurrency:    foreeTx.SrcAmt.Currency,
		DestAccId:      coutAcc.ID,
		DestAccSummary: foreeTx.ContactAcc.GetLegalName(),
		DestAmount:     foreeTx.DestAmt.Amount,
		DestCurrency:   foreeTx.DestAmt.Currency,
		TotalAmount:    foreeTx.TotalAmt.Amount,
		TotalCurrency:  foreeTx.TotalAmt.Currency,
		OwnerId:        user.ID,
		SrcAccount:     ciAcc,
		DestAccount:    coutAcc,
	}

	if len(feeJoints) > 0 {
		txSummary.FeeAmount = foreeTx.TotalFeeAmt.Amount
		txSummary.FeeCurrency = foreeTx.TotalFeeAmt.Currency
	}

	if len(rewards) > 0 {
		txSummary.RewardAmount = foreeTx.TotalRewardAmt.Amount
		txSummary.RewardCurrency = foreeTx.TotalRewardAmt.Currency
	}

	foreeTx.Summary = &txSummary

	quoteId, err := t.txQuoteRepo.InsertTxQuote(ctx, transaction.TxQuote{
		Tx:     foreeTx,
		OwerId: user.ID,
	})

	if err != nil {
		foree_logger.Logger.Error("QuoteTx_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("QuoteTx_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId)

	return &QuoteTransactionDTO{
		QuoteId: quoteId,
		TxSum:   *NewTxSummaryDetailDTO(&txSummary),
	}, nil
}

// Investigation: Can we run parallel query inside single mysql transaction.
func (t *TransactionService) CreateTx(ctx context.Context, req CreateTransactionReq) (*TxSummaryDetailDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("CreateTx_FAIL",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	user := session.User
	quote := t.txQuoteRepo.GetUniqueById(ctx, req.QuoteId)

	if quote == nil {
		foree_logger.Logger.Warn("CreateTx_FAIL",
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
		if len(foreeTx.Rewards) > 0 {
			rewards, err := t.rewardRepo.GetAllRewardByOwnerIdAndIds(ctx, session.UserId, foreeTx.RewardIds)
			if err != nil {
				foree_logger.Logger.Error("CreateTx_FAIL",
					"ip", loadRealIp(ctx),
					"userId", session.UserId,
					"sessionId", req.SessionId,
					"quoteId", req.QuoteId,
					"cause", err.Error(),
				)
				rewardErr = transport.WrapInteralServerError(err)
				return
			}
			for _, reward := range rewards {
				if reward.Status != promotion.RewardStatusActive {
					foree_logger.Logger.Warn("CreateTx_FAIL",
						"ip", loadRealIp(ctx),
						"userId", session.UserId,
						"sessionId", req.SessionId,
						"quoteId", req.QuoteId,
						"rewardStatus", reward.Status,
						"cause", "reward is not in active",
					)
					rewardErr = transport.NewInteralServerError("user `%v` try to create a transaction with reward `%v` in status `%s`", user.ID, reward.ID, reward.Status)
					return
				}
			}
		}
	}
	wg.Add(1)
	go rewardChecker()

	// Recheck and update limit
	var limitErr transport.HError
	limitChecker := func() {
		defer wg.Done()

		dailyLimit, err := t.txLimitService.getDailyTxLimit(ctx, *session)
		if err != nil {
			foree_logger.Logger.Error("CreateTx_FAIL",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"quoteId", req.QuoteId,
				"cause", err.Error(),
			)
			limitErr = transport.WrapInteralServerError(err)
			return
		}

		if foreeTx.SrcAmt.Amount+dailyLimit.UsedAmt.Amount > dailyLimit.MaxAmt.Amount {
			foree_logger.Logger.Warn("QuoteTx_FAIL",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"requstAmount", foreeTx.SrcAmt.Amount,
				"requstCurrency", foreeTx.SrcAmt.Currency,
				"remainingAmount", dailyLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount,
				"maxAmount", dailyLimit.MaxAmt.Amount,
				"cause", "overlimit",
			)
			limitErr = transport.NewInteralServerError("user `%v` try to create a transaction with `%v` but the remaining limit is `%v`", user.ID, foreeTx.SrcAmt.Amount, dailyLimit.MaxAmt.Amount-dailyLimit.UsedAmt.Amount)
			return
		}

		if _, err := t.txLimitService.addDailyTxLimit(ctx, *session, foreeTx.SrcAmt); err != nil {
			foree_logger.Logger.Error("CreateTx_FAIL",
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
			foree_logger.Logger.Error("CreateTx_FAIL",
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
	var sumId int64
	createSummaryTx := func() {
		defer wg.Done()
		foreeTx.Summary.ParentTxId = foreeTxID
		foreeTx.Summary.NBPReference = transaction.GenerateNbpId(foree_constant.DefaultNBPIdPrefix, foreeTxID)
		sumId, err = t.txSummaryRepo.InsertTxSummary(ctx, *foreeTx.Summary)
		if err != nil {
			foree_logger.Logger.Error("CreateTx_FAIL",
				"ip", loadRealIp(ctx),
				"userId", session.UserId,
				"sessionId", req.SessionId,
				"quoteId", req.QuoteId,
				"cause", err.Error(),
			)
			txSummaryErr = transport.WrapInteralServerError(err)
		}
		foreeTx.Summary.ID = sumId
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
			feeJoin.OwnerId = session.UserId
			_, err := t.feeJointRepo.InsertFeeJoint(ctx, *feeJoin)
			if err != nil {
				foree_logger.Logger.Error("CreateTx_FAIL",
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
			r.Status = promotion.RewardStatusPending
			err := t.rewardRepo.UpdateRewardTxById(ctx, *r)
			if err != nil {
				foree_logger.Logger.Error("CreateTx_FAIL",
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

	// go t.txProcessor.createAndProcessTx(*foreeTx)
	t.txProcessor.createAndProcessTx(*foreeTx)

	foree_logger.Logger.Info("CreateTx_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "foreeTxId", foreeTxID)

	return t.GetTxSummary(ctx, GetTransactionReq{
		SessionReq:    req.SessionReq,
		TransactionId: sumId,
	})
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
		foree_logger.Logger.Info("GetDailyTxLimit_FAIL",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	limit, err := t.txLimitService.getDailyTxLimit(ctx, *session)
	if err != nil {
		foree_logger.Logger.Error("GetDailyTxLimit_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}
	foree_logger.Logger.Debug("GetDailyTxLimit_SUCCESS",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
	)
	return NewDailyTxLimitDTO(limit), nil
}

func (t *TransactionService) GetReward(ctx context.Context, req transport.SessionReq) ([]*RewardDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("GetReward_FAIL",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	rewards, err := t.rewardRepo.GetAllActiveRewardByOwnerId(ctx, session.UserId)
	if err != nil {
		foree_logger.Logger.Error("GetReward_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	ret := make([]*RewardDTO, len(rewards))
	for i, v := range rewards {
		ret[i] = NewRewardDTO(v)
	}

	foree_logger.Logger.Debug("GetReward_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId)
	return ret, nil
}

func (t *TransactionService) GetTxSummary(ctx context.Context, req GetTransactionReq) (*TxSummaryDetailDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("GetTxSummary_FAIL",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}

	summaryTx, err := t.txSummaryRepo.GetUniqueTxSummaryByOwnerAndId(ctx, session.UserId, req.TransactionId)
	if err != nil {
		foree_logger.Logger.Error("GetTxSummary_FAIL",
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

	wg := sync.WaitGroup{}

	var interacAcc *account.InteracAccount
	getInteracAccount := func() {
		defer wg.Done()
		interacAcc, err = t.interacAccountRepo.GetUniqueInteracAccountById(ctx, summaryTx.SrcAccId)
		if err != nil {
			foree_logger.Logger.Warn("GetTxSummary",
				"sumTxId", req.TransactionId,
				"msg", "load interacAccount failed",
				"cause", err.Error(),
			)
		}
	}
	wg.Add(1)
	go getInteracAccount()

	var contactAcc *account.ContactAccount
	getContactAccount := func() {
		defer wg.Done()
		contactAcc, err = t.contactAccountRepo.GetUniqueContactAccountById(ctx, summaryTx.DestAccId)
		if err != nil {
			foree_logger.Logger.Warn("GetTxSummary",
				"sumTxId", req.TransactionId,
				"msg", "load contactAccount failed",
				"cause", err.Error(),
			)
		}
	}
	wg.Add(1)
	go getContactAccount()

	wg.Wait()
	summaryTx.SrcAccount = interacAcc
	summaryTx.DestAccount = contactAcc

	foree_logger.Logger.Debug("GetTxSummary_SUCCESS",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
	)
	return NewTxSummaryDetailDTO(summaryTx), nil
}

func (t *TransactionService) QuerySummaryTxs(ctx context.Context, req QueryTransactionReq) ([]*TxSummaryDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("QuerySummaryTxs_FAIL",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}
	//TODO: limit, offset pruning
	var summaryTxs []*transaction.TxSummary
	var err error

	if req.Status == "" || req.Status == "All" {
		summaryTxs, err = t.txSummaryRepo.GetAllTxSummaryByOwnerIdWithPagination(ctx, session.UserId, req.Limit, req.Offset)
	} else {
		summaryTxs, err = t.txSummaryRepo.GetAllTxSummaryByOwnerIdAndStatusWithPagination(ctx, session.UserId, req.Status, req.Limit, req.Offset)
	}

	if err != nil {
		foree_logger.Logger.Error("QuerySummaryTxs_FAIL",
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

	foree_logger.Logger.Debug("QuerySummaryTxs_SUCCESS",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
	)
	return rets, nil
}

func (t *TransactionService) CountSummaryTxs(ctx context.Context, req QueryTransactionReq) (*TxSummarieCountDTO, transport.HError) {
	session, sErr := t.authService.GetSession(ctx, req.SessionId)
	if session == nil {
		foree_logger.Logger.Info("CountSummaryTxs_FAIL",
			"sessionId", req.SessionId,
			"ip", loadRealIp(ctx),
			"cause", "session no found",
		)
		return nil, sErr
	}
	var count int
	var err error

	if req.Status == "" || req.Status == "All" {
		count, err = t.txSummaryRepo.CountTxSummaryByOwnerId(ctx, session.UserId)
	} else {
		count, err = t.txSummaryRepo.CountTxSummaryByOwnerIdAndStatus(ctx, session.UserId, req.Status)
	}

	if err != nil {
		foree_logger.Logger.Error("CountSummaryTxs_FAIL",
			"ip", loadRealIp(ctx),
			"userId", session.UserId,
			"sessionId", req.SessionId,
			"cause", err.Error(),
		)
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Debug("CountSummaryTxs_SUCCESS",
		"ip", loadRealIp(ctx),
		"userId", session.UserId,
		"sessionId", req.SessionId,
	)
	return &TxSummarieCountDTO{
		Count: count,
	}, nil
}

// Check transaction status, see if is able to cancel.
// func (t *TransactionService) CancelTransaction(ctx context.Context, req CancelTransactionReq) (*TxCancelDTO, transport.HError) {
// 	session, sErr := t.authService.GetSession(ctx, req.SessionId)
// 	if session == nil {
// 		foree_logger.Logger.Info("CancelTransaction_FAIL",
// 			"sessionId", req.SessionId,
// 			"ip", loadRealIp(ctx),
// 			"cause", "session no found",
// 		)
// 		return nil, sErr
// 	}

// 	summaryTx, err := t.txSummaryRepo.GetUniqueTxSummaryByOwnerAndId(ctx, session.UserId, req.TransactionId)
// 	if err != nil {
// 		foree_logger.Logger.Error("CancelTransaction_FAIL",
// 			"ip", loadRealIp(ctx),
// 			"userId", session.UserId,
// 			"sessionId", req.SessionId,
// 			"cause", err.Error(),
// 		)
// 		return nil, transport.WrapInteralServerError(err)
// 	}

// 	if summaryTx == nil {
// 		return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "no found")
// 	}

// 	fTx, err := t.txProcessor.LoadTx(summaryTx.ParentTxId)
// 	if err != nil {
// 		foree_logger.Logger.Error("CancelTransaction_FAIL",
// 			"ip", loadRealIp(ctx),
// 			"userId", session.UserId,
// 			"sessionId", req.SessionId,
// 			"cause", err.Error(),
// 		)
// 		return nil, transport.WrapInteralServerError(err)
// 	}

// 	if fTx.Stage == transaction.TxStageInteracCI && fTx.StageStatus == transaction.TxStatusSent {
// 		resp, err := t.scotiaClient.CancelPayment(scotia.CancelPaymentRequest{
// 			PaymentId:    fTx.CI.ScotiaPaymentId,
// 			CancelReason: req.CancelReason,
// 		})
// 		//TODO: log
// 		if err != nil {
// 			foree_logger.Logger.Error("CancelTransaction_FAIL",
// 				"ip", loadRealIp(ctx),
// 				"userId", session.UserId,
// 				"sessionId", req.SessionId,
// 				"cause", err.Error(),
// 			)
// 			return nil, transport.WrapInteralServerError(err)
// 		} else if resp.StatusCode/100 != 2 {
// 			return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "transaction can not cancel")
// 		}
// 	} else if fTx.Stage == transaction.TxStageNBPCO && fTx.StageStatus == transaction.TxStatusSent && fTx.COUT.CashOutAcc.Type == foree_constant.ContactAccountTypeCash {
// 		resp, err := t.nbpClient.CancelTransaction(nbp.CancelTransactionRequest{
// 			GlobalId:           fTx.COUT.NBPReference,
// 			CancellationReason: req.CancelReason,
// 		})
// 		//TODO: log
// 		if err != nil {
// 			foree_logger.Logger.Error("CancelTransaction_FAIL",
// 				"ip", loadRealIp(ctx),
// 				"userId", session.UserId,
// 				"sessionId", req.SessionId,
// 				"cause", err.Error(),
// 			)
// 			return nil, transport.WrapInteralServerError(err)
// 		} else if resp.StatusCode/100 != 2 {
// 			return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "transaction can not cancel")
// 		}
// 	} else {
// 		foree_logger.Logger.Warn("CancelTransaction_FAIL",
// 			"ip", loadRealIp(ctx),
// 			"userId", session.UserId,
// 			"sessionId", req.SessionId,
// 			"foreeTxId", fTx.ID,
// 			"Stage", fTx.Stage,
// 			"StageStatus", fTx.StageStatus,
// 		)
// 		return nil, transport.NewFormError("Invalid transaction cancel request", "transactionId", "transaction can not cancel")
// 	}

// 	foree_logger.Logger.Info("CancelTransaction_SUCCESS",
// 		"ip", loadRealIp(ctx),
// 		"userId", session.UserId,
// 		"sessionId", req.SessionId,
// 		"foreeTxId", fTx.ID,
// 	)
// 	return &TxCancelDTO{
// 		TransactionId: req.TransactionId,
// 		Message:       "cancel successfully",
// 	}, nil
// }
