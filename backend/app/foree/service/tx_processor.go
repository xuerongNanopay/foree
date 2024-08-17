package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
	time_util "xue.io/go-pay/util/time"
)

const (
	FeeName string = "FOREE_TX_CAD_FEE"
)

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

type TxProcessorConfig struct {
}

// It is the internal service for transaction process.

type TxProcessor struct {
	interacTxRepo      *transaction.InteracCITxRepo
	npbTxRepo          *transaction.NBPCOTxRepo
	idmTxRepo          *transaction.IdmTxRepo
	txHistoryRepo      *transaction.TxHistoryRepo
	txSummaryRepo      *transaction.TxSummaryRepo
	dailyTxLimiteRepo  *transaction.DailyTxLimitRepo
	foreeTxRepo        *transaction.ForeeTxRepo
	rateRepo           *transaction.RateRepo
	userRepo           *auth.UserRepo
	contactRepo        *account.ContactAccountRepo
	interacRepo        *account.InteracAccountRepo
	feeRepo            *transaction.FeeRepo
	feeJointRepo       *transaction.FeeJointRepo
	promoCodeRepo      *transaction.PromoCodeRepo
	promoCodeJointRepo *transaction.PromoCodeJointRepo
	rewardRepo         *transaction.RewardRepo
	processingMap      []map[int64]*transaction.ForeeTx // Avoid duplicate process
	processingLock     sync.RWMutex
}

func (p *TxProcessor) quoteTx(user auth.User, quote QuoteTransactionReq) (*transaction.ForeeTx, error) {
	ctx := context.Background()
	rate, err := p.rateRepo.GetUniqueRateById(ctx, transaction.GenerateRateId(quote.SrcCurrency, quote.DestCurrency))
	if err != nil {
		return nil, err
	}
	if rate == nil {
		return nil, fmt.Errorf("user `%v` try to create transaction with unkown rate `%s`", user.ID, transaction.GenerateRateId(quote.SrcCurrency, quote.DestCurrency))
	}

	//Reward
	var reward *transaction.Reward
	if len(quote.RewardIds) == 1 {
		rewardId := quote.RewardIds[1]
		r, err := p.rewardRepo.GetUniqueRewardById(ctx, rewardId)
		if err != nil {
			return nil, err
		}
		if r == nil {
			return nil, fmt.Errorf("user `%v` try to redeem unknown reward `%v`", user.ID, rewardId)
		}
		if r.OwnerId != user.ID {
			return nil, fmt.Errorf("user `%v` try to redeem reward `%v` that is belong to `%v`", user.ID, rewardId, rewardId, r.OwnerId)
		}
		if r.Status != transaction.RewardStatusActive {
			return nil, fmt.Errorf("user `%v` try to redeem reward `%v` that is currently in status `%v`", user.ID, rewardId, r.Status)
		}
		if r.Amt.Currency != quote.SrcCurrency {
			return nil, fmt.Errorf("user `%v` try to redeem reward `%v` that apply currency `%v` to currency `%v`", user.ID, rewardId, r.Amt.Currency, quote.SrcCurrency)
		}
		if (quote.SrcAmount - float64(r.Amt.Amount)) < 10 {
			return nil, fmt.Errorf("user `%v` try to redeem reward `%v` with srcAmount `%v`", user.ID, rewardId, quote.SrcCurrency)
		}
		reward = r
	}

	//TODO: PromoCode
	//Don't return err. Just ignore the promocode reward.
	if quote.PromoCode != "" {
		promoCode, err := p.promoCodeRepo.GetUniquePromoCodeByCode(ctx, quote.PromoCode)
		if err != nil {
			//TODO: log
			goto existpromo
		}
		if promoCode == nil {
			//TODO: log
			goto existpromo
		}
		if !promoCode.IsValid() {
			//TODO: log
			goto existpromo
		}
		if quote.SrcCurrency != promoCode.MinAmt.Currency {
			//TODO: log
			goto existpromo
		}
		if quote.SrcAmount < float64(promoCode.MinAmt.Amount) {
			//TODO: log
			goto existpromo
		}
		//TODO: check account limit.
	}

existpromo:

	//Fee
	fee, err := p.feeRepo.GetUniqueFeeByName(FeeName)
	if err != nil {
		return nil, err
	}
	if fee == nil {
		return nil, fmt.Errorf("fee `%v` not found", FeeName)
	}

	joint, err := fee.MaybeApplyFee(types.AmountData{Amount: types.Amount(quote.SrcAmount), Currency: quote.SrcCurrency})
	if err != nil {
		return nil, err
	}
	if joint != nil {
		joint.Description = fee.Description
		joint.OwnerId = user.ID
	}

	//Total
	totalAmt := types.AmountData{}

	if joint != nil {
		totalAmt.Amount += joint.Amt.Amount
	}

	if reward != nil {
		totalAmt.Amount -= reward.Amt.Amount
	}

	//Summary

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
		SrcAccId:           quote.SrcAccId,
		DestAccId:          quote.DestAccId,
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

func (p *TxProcessor) GetDailyTxLimit(user auth.User) (*transaction.DailyTxLimit, error) {
	ctx := context.Background()
	return p.getDailyTxLimit(ctx, user)
}

func (p *TxProcessor) addDailyTxLimit(ctx context.Context, user auth.User, amt types.AmountData) (*transaction.DailyTxLimit, error) {
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

func (p *TxProcessor) minusDailyTxLimit(ctx context.Context, user auth.User, amt types.AmountData) (*transaction.DailyTxLimit, error) {
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
func (p *TxProcessor) getDailyTxLimit(ctx context.Context, user auth.User) (*transaction.DailyTxLimit, error) {
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

func (p *TxProcessor) createTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *TxProcessor) InsertTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *TxProcessor) loadTx(id int64) (*transaction.ForeeTx, error) {
	ctx := context.Background()
	foree, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load CI
	ci, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, foree.ID)
	if err != nil {
		return nil, err
	}
	if ci == nil {
		return nil, fmt.Errorf("InteracCITx no found for ForeeTx `%v`", foree.ID)
	}

	srcInteracAcc, err := p.interacRepo.GetUniqueInteracAccountById(ctx, ci.SrcInteracAccId)
	if err != nil {
		return nil, err
	}
	if srcInteracAcc == nil {
		return nil, fmt.Errorf("SrcInteracAcc no found for InteracCITx `%v`", ci.SrcInteracAccId)
	}
	ci.SrcInteracAcc = srcInteracAcc

	destInteracAcc, err := p.interacRepo.GetUniqueInteracAccountById(ctx, ci.DestInteracAccId)
	if err != nil {
		return nil, err
	}
	if destInteracAcc == nil {
		return nil, fmt.Errorf("DestInteracAcc no found for InteracCITx `%v`", ci.DestInteracAccId)
	}
	ci.DestInteracAcc = destInteracAcc
	foree.CI = ci

	// Load IDM
	idm, err := p.idmTxRepo.GetUniqueIDMTxByParentTxId(ctx, foree.ID)
	if err != nil {
		return nil, err
	}
	if idm == nil {
		return nil, fmt.Errorf("IDMTx no found for ForeeTx `%v`", foree.ID)
	}
	foree.IDM = idm

	// Load COUT
	cout, err := p.npbTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, foree.ID)
	if err != nil {
		return nil, err
	}
	if cout == nil {
		return nil, fmt.Errorf("NBPCOTx no found for ForeeTx `%v`", foree.ID)
	}

	destContactAcc, err := p.contactRepo.GetUniqueContactAccountById(ctx, cout.DestContactAccId)
	if err != nil {
		return nil, err
	}
	if destContactAcc == nil {
		return nil, fmt.Errorf("DestContactAcc no found for NBPCOTx `%v`", cout.DestContactAccId)
	}
	cout.DestContactAcc = destContactAcc
	foree.COUT = cout

	// TODO: fees?, rewards?

	return foree, nil
}

// TODO: change argument to int64
func (p *TxProcessor) processTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if tx.Type != transaction.TxTypeInteracToNBP {
		return nil, fmt.Errorf("unknow ForeeTx type `%s`", tx.Type)
	}
	var err error
	var nTx *transaction.ForeeTx
	maxLoop := 16
	i := 0
	ctx := context.Background()
	for {
		nTx, err = p.doProcessTx(ctx, tx)
		if err != nil {
			return nil, err
		}
		if tx.CurStage == nTx.CurStage && nTx.CurStageStatus == tx.CurStageStatus {
			return nTx, nil
		}
		// Record the history.
		go p.recordTxHistory(*transaction.NewTxHistory(nTx, ""))
		tx = *nTx

		if i > maxLoop {
			return nil, fmt.Errorf("unexpect looping for ForeeTx `%v`", nTx.ID)
		}
		i += 1
	}

}

func (p *TxProcessor) recordTxHistory(h transaction.TxHistory) {
	if _, err := p.txHistoryRepo.InserTxHistory(h); err != nil {
		fmt.Println(err.Error())
	}

}

func (p *TxProcessor) doProcessTx(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if tx.Status == transaction.TxStatusInitial {
		tx.Status = transaction.TxStatusProcessing
		tx.CurStage = transaction.TxStageInteracCI
		tx.CurStageStatus = transaction.TxStatusInitial
		return &tx, nil
	}
	if tx.Status == transaction.TxStatusCompleted || tx.Status == transaction.TxStatusCancelled || tx.Status == transaction.TxStatusRejected {
		//TODO: log warn.
		return &tx, nil
	}

	switch tx.CurStage {
	case transaction.TxStageInteracCI:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send scotia API
			//Set to Send
		case transaction.TxStatusSent:
			//Check status from scotia API.(Webhook, or cron)
			//Just do noting waiting for cron
		case transaction.TxStatusCompleted:
			tx.CurStage = transaction.TxStageInteracCI
			tx.CurStageStatus = transaction.TxStatusInitial
			return &tx, nil
		case transaction.TxStatusRejected:
			//Set to reject
		case transaction.TxStatusCancelled:
			// set to cancel
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	case transaction.TxStageIDM:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send IDMAPI
			//Set to Send
		case transaction.TxStatusCompleted:
			//Move to next stage
			tx.CurStage = transaction.TxStageNBPCO
			tx.CurStageStatus = transaction.TxStatusInitial
			return &tx, nil
		case transaction.TxStatusRejected:
			// Set to ForeeTx to terminal status.
			tx.Status = transaction.TxStatusRejected
			tx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", tx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			// Close remaing non-terminated transactions.
			nT, err := p.closeRemainingTx(ctx, tx)
			if err != nil {
				return nil, err
			}
			go p.maybeRefund(*nT)
			return nT, nil
		case transaction.TxStatusSuspend:
			//Wait to approve
			//Log warn?
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	case transaction.TxStageNBPCO:
		switch tx.CurStageStatus {
		case transaction.TxStatusInitial:
			//TODO: call send NBP API
		case transaction.TxStatusSent:
			//Check status from NBP API.
			//Or just wait for clone
		case transaction.TxStatusCompleted:
			tx.Status = transaction.TxStatusCompleted
			tx.Conclusion = fmt.Sprintf("Complete at %s.", time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			return &tx, nil
			// set tx sum to complete
		case transaction.TxStatusRejected:
			tx.Status = transaction.TxStatusRejected
			tx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", tx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			go p.maybeRefund(tx)
			return &tx, nil
		case transaction.TxStatusCancelled:
			tx.Status = transaction.TxStatusCancelled
			tx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", tx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, tx); err != nil {
				return nil, err
			}
			go p.maybeRefund(tx)
			return &tx, nil
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", tx.ID, tx.CurStageStatus, tx.CurStage)
		}
	default:
		return nil, fmt.Errorf("transaction `%v` in unknown stage `%s`", tx.ID, tx.CurStage)
	}
	return nil, nil
}

func (p *TxProcessor) closeRemainingTx(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	switch tx.CurStage {
	case transaction.TxStageInteracCI:
		idm := tx.IDM
		co := tx.COUT
		idm.Status = transaction.TxStatusClosed
		co.Status = transaction.TxStatusClosed
		if err := p.idmTxRepo.UpdateIDMTxById(ctx, *idm); err != nil {
			return nil, err
		}
		if err := p.npbTxRepo.UpdateNBPCOTxById(ctx, *co); err != nil {
			return nil, err
		}
		return &tx, nil
	case transaction.TxStageIDM:
		co := tx.COUT
		co.Status = transaction.TxStatusClosed
		if err := p.npbTxRepo.UpdateNBPCOTxById(ctx, *co); err != nil {
			return nil, err
		}
		return &tx, nil
	default:
		//TODO: Log warn
		return &tx, nil
	}
}

func (p *TxProcessor) maybeRefund(tx transaction.ForeeTx) {
	//TODO: implement
}

// TODO: change argement to id.
func (p *TxProcessor) approveIDM(ctx context.Context, tx transaction.ForeeTx) {
	if tx.CurStage == transaction.TxStageIDM && tx.CurStageStatus == transaction.TxStatusSuspend {

	}
	//TODO: implement
}

// TODO: change argement to id.
func (p *TxProcessor) rejectIDM(ctx context.Context, tx transaction.ForeeTx) {
	if tx.CurStage == transaction.TxStageIDM && tx.CurStageStatus == transaction.TxStatusSuspend {

	}
}
