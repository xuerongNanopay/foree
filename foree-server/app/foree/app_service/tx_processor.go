package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_constant "xue.io/go-pay/app/foree/constant"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	time_util "xue.io/go-pay/util/time"
)

// Goal of Txprocessor:
//	1. dispatch tx to sub tx processor.
// 	2. handle terminal status of subprocessor.
// 	3. it use curStage to navigate to different subprocessor.
// It is the internal service for transaction process.

func NewTxProcessor(
	db *sql.DB,
	interacTxRepo *transaction.InteracCITxRepo,
	nbpTxRepo *transaction.NBPCOTxRepo,
	idmTxRepo *transaction.IdmTxRepo,
	txHistoryRepo *transaction.TxHistoryRepo,
	txSummaryRepo *transaction.TxSummaryRepo,
	foreeTxRepo *transaction.ForeeTxRepo,
	interacRefundTxRepo *transaction.InteracRefundTxRepo,
	rewardRepo *transaction.RewardRepo,
	dailyTxLimiteRepo *transaction.DailyTxLimitRepo,
	userRepo *auth.UserRepo,
	contactAccountRepo *account.ContactAccountRepo,
	interacAccountRepo *account.InteracAccountRepo,
) *TxProcessor {
	return &TxProcessor{
		db:                  db,
		interacTxRepo:       interacTxRepo,
		nbpTxRepo:           nbpTxRepo,
		idmTxRepo:           idmTxRepo,
		txHistoryRepo:       txHistoryRepo,
		txSummaryRepo:       txSummaryRepo,
		foreeTxRepo:         foreeTxRepo,
		interacRefundTxRepo: interacRefundTxRepo,
		rewardRepo:          rewardRepo,
		dailyTxLimiteRepo:   dailyTxLimiteRepo,
		userRepo:            userRepo,
		contactAccountRepo:  contactAccountRepo,
		interacAccountRepo:  interacAccountRepo,
	}
}

func (p *TxProcessor) SetCITxProcessor(ciTxProcessor *CITxProcessor) {
	p.ciTxProcessor = ciTxProcessor
}

func (p *TxProcessor) SetIDMTxProcessor(idmTxProcessor *IDMTxProcessor) {
	p.idmTxProcessor = idmTxProcessor
}

func (p *TxProcessor) SetNBPTxProcessor(nbpTxProcessor *NBPTxProcessor) {
	p.nbpTxProcessor = nbpTxProcessor
}

type TxProcessor struct {
	db                  *sql.DB
	interacTxRepo       *transaction.InteracCITxRepo
	nbpTxRepo           *transaction.NBPCOTxRepo
	idmTxRepo           *transaction.IdmTxRepo
	txHistoryRepo       *transaction.TxHistoryRepo
	txSummaryRepo       *transaction.TxSummaryRepo
	foreeTxRepo         *transaction.ForeeTxRepo
	interacRefundTxRepo *transaction.InteracRefundTxRepo
	rewardRepo          *transaction.RewardRepo
	dailyTxLimiteRepo   *transaction.DailyTxLimitRepo
	userRepo            *auth.UserRepo
	contactAccountRepo  *account.ContactAccountRepo
	interacAccountRepo  *account.InteracAccountRepo
	ciTxProcessor       *CITxProcessor
	idmTxProcessor      *IDMTxProcessor
	nbpTxProcessor      *NBPTxProcessor
}

func (p *TxProcessor) createAndProcessTx(tx transaction.ForeeTx) {
	foreeTx, err := p.createFullTx(tx)
	if err != nil {
		foree_logger.Logger.Error("CreateAndProcessTx_Fail",
			"foreeTxId", tx.ID,
			"cause", err.Error(),
		)
		return
	}

	_, err = p.processTx(*foreeTx)
	if err != nil {
		foree_logger.Logger.Error("CreateAndProcessTx_Fail",
			"foreeTxId", tx.ID,
			"cause", err.Error(),
		)
	}
}

func (p *TxProcessor) loadAndProcessTx(foreeId int64) (*transaction.ForeeTx, error) {
	fTx, err := p.loadTx(foreeId, true)
	if err != nil {
		return nil, err
	}

	go func() {
		_, err := p.processTx(*fTx)
		if err != nil {
			//TODO log
		}
	}()

	return fTx, nil
}

// Create CI, COUT, IDM for ForeeTx
func (p *TxProcessor) createFullTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	wg := sync.WaitGroup{}
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	_, err = p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, tx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	// Create CI
	var ciTx *transaction.InteracCITx
	var ciErr error
	createCI := func() {
		defer wg.Done()
		ciId, err := p.interacTxRepo.InsertInteracCITx(ctx, transaction.InteracCITx{
			Status:      transaction.TxStatusInitial,
			CashInAccId: tx.CinAccId,
			EndToEndId:  tx.Summary.NBPReference,
			Amt:         tx.SrcAmt,
			ParentTxId:  tx.ID,
			OwnerId:     tx.OwnerId,
		})
		if err != nil {
			ciErr = err
			return
		}
		ci, err := p.interacTxRepo.GetUniqueInteracCITxById(ctx, ciId)
		if err != nil {
			ciErr = err
			return
		}
		ciTx = ci
	}
	wg.Add(1)
	go createCI()

	// Create IDM
	var idmTx *transaction.IDMTx
	var idmErr error
	createIDM := func() {
		defer wg.Done()
		idmId, err := p.idmTxRepo.InsertIDMTx(ctx, transaction.IDMTx{
			Status:     transaction.TxStatusInitial,
			Ip:         tx.Ip,
			UserAgent:  tx.UserAgent,
			ParentTxId: tx.ID,
			OwnerId:    tx.OwnerId,
		})
		if err != nil {
			idmErr = err
			return
		}
		idm, err := p.idmTxRepo.GetUniqueIDMTxById(ctx, idmId)
		if err != nil {
			idmErr = err
			return
		}
		idmTx = idm
	}
	wg.Add(1)
	go createIDM()

	// Create Cout
	var coutTx *transaction.NBPCOTx
	var coutErr error
	createCout := func() {
		defer wg.Done()
		coutId, err := p.nbpTxRepo.InsertNBPCOTx(ctx, transaction.NBPCOTx{
			Status:       transaction.TxStatusInitial,
			Amt:          tx.DestAmt,
			NBPReference: tx.Summary.NBPReference,
			CashOutAccId: tx.CoutAccId,
			ParentTxId:   tx.ID,
			OwnerId:      tx.OwnerId,
		})
		if err != nil {
			coutErr = err
			return
		}
		cout, err := p.nbpTxRepo.GetUniqueNBPCOTxById(ctx, coutId)
		if err != nil {
			coutErr = err
			return
		}
		coutTx = cout
	}

	wg.Add(1)
	go createCout()

	wg.Wait()
	if ciErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_Fail",
			"ip", loadRealIp(ctx),
			"foreeTxId", tx.ID,
			"cause", ciErr.Error(),
		)
		return nil, ciErr
	}
	if idmErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_Fail",
			"ip", loadRealIp(ctx),
			"foreeTxId", tx.ID,
			"cause", idmErr.Error(),
		)
		return nil, idmErr
	}
	if coutErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_Fail",
			"ip", loadRealIp(ctx),
			"foreeTxId", tx.ID,
			"cause", coutErr.Error(),
		)
		return nil, coutErr
	}

	tx.CI = ciTx
	tx.IDM = idmTx
	tx.COUT = coutTx

	if err = dTx.Commit(); err != nil {
		foree_logger.Logger.Error("CreateFullTx_Fail",
			"ip", loadRealIp(ctx),
			"foreeTxId", tx.ID,
			"cause", err.Error(),
		)
		return nil, err
	}
	return &tx, nil
}

func (p *TxProcessor) LoadTx(id int64) (*transaction.ForeeTx, error) {
	return p.loadTx(id, true)
}

func (p *TxProcessor) loadTx(id int64, isEmptyCheck bool) (*transaction.ForeeTx, error) {
	ctx := context.Background()
	foreeTx, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, id)
	if err != nil {
		foree_logger.Logger.Error("loadTx_Fail", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if foreeTx == nil {
		foree_logger.Logger.Warn("loadTx_Fail",
			"foreeTxId", id,
			"cause", "foreeTx no found",
		)
		return nil, fmt.Errorf("ForeeTx no found with id `%v`", id)
	}

	// Load CI
	ci, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, foreeTx.ID)
	if err != nil {
		foree_logger.Logger.Error("loadTx_Fail", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && ci == nil {
		foree_logger.Logger.Warn("loadTx_Fail",
			"foreeTxId", id,
			"cause", "InteracCITx no found",
		)
		return nil, fmt.Errorf("InteracCITx no found for ForeeTx `%v`", foreeTx.ID)
	}

	CashInAcc, err := p.interacAccountRepo.GetUniqueInteracAccountById(ctx, ci.CashInAccId)
	if err != nil {
		foree_logger.Logger.Error("loadTx_Fail", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && CashInAcc == nil {
		foree_logger.Logger.Warn("loadTx_Fail",
			"foreeTxId", id,
			"interactTxId", ci.ID,
			"interacAccountId", ci.CashInAccId,
			"cause", "interac account no found",
		)
		return nil, fmt.Errorf("CashInAcc no found for InteracCITx `%v`", ci.CashInAccId)
	}
	ci.CashInAcc = CashInAcc

	foreeTx.CI = ci

	// Load IDM
	idm, err := p.idmTxRepo.GetUniqueIDMTxByParentTxId(ctx, foreeTx.ID)
	if err != nil {
		foree_logger.Logger.Error("loadTx_Fail", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && idm == nil {
		foree_logger.Logger.Warn("loadTx_Fail",
			"foreeTxId", id,
			"cause", "idmTx no found",
		)
		return nil, fmt.Errorf("IDMTx no found for ForeeTx `%v`", foreeTx.ID)
	}
	foreeTx.IDM = idm

	// Load COUT
	cout, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, foreeTx.ID)
	if err != nil {
		foree_logger.Logger.Error("loadTx_Fail", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && cout == nil {
		foree_logger.Logger.Warn("loadTx_Fail",
			"foreeTxId", id,
			"cause", "nbpCOTx no found",
		)
		return nil, fmt.Errorf("NBPCOTx no found for ForeeTx `%v`", foreeTx.ID)
	}

	CashOutAcc, err := p.contactAccountRepo.GetUniqueContactAccountById(ctx, cout.CashOutAccId)
	if err != nil {
		foree_logger.Logger.Error("loadTx_Fail", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && CashOutAcc == nil {
		foree_logger.Logger.Warn("loadTx_Fail",
			"foreeTxId", id,
			"nbpCoTxId", cout.ID,
			"contactAccountId", cout.CashOutAccId,
			"cause", "CashOutAcc no found",
		)
		return nil, fmt.Errorf("CashOutAcc no found for NBPCOTx `%v`", cout.CashOutAccId)
	}
	cout.CashOutAcc = CashOutAcc
	foreeTx.COUT = cout

	// Load User
	user, err := p.userRepo.GetUniqueUserById(foreeTx.OwnerId)
	if err != nil {
		foree_logger.Logger.Error("loadTx_Fail", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && user == nil {
		foree_logger.Logger.Warn("loadTx_Fail",
			"foreeTxId", id,
			"ownerId", foreeTx.OwnerId,
			"cause", "owner no found",
		)
		return nil, fmt.Errorf("owner `%v` no found for ForeeTx `%v`", foreeTx.OwnerId, foreeTx.ID)
	}
	foreeTx.Owner = user

	// TODO: fees?, rewards?

	return foreeTx, nil
}

func (p *TxProcessor) processRootTx(fTx transaction.ForeeTx) error {
	dbTx, err := p.db.Begin()
	if err != nil {
		dbTx.Rollback()
		foree_logger.Logger.Error("processRootTx", "foreeTxId", fTx.ID, "cause", err.Error())
		return err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dbTx)

	curForeeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
	if err != nil {
		dbTx.Rollback()
		foree_logger.Logger.Error("processRootTx", "foreeTxId", fTx.ID, "cause", err.Error())
		return err
	}
	if curForeeTx.CurStage != fTx.CurStage {
		foree_logger.Logger.Warn("processRootTx",
			"foreeTxId", fTx.ID,
			"transactionStage", fTx.CurStage,
			"transactionStageIndb", curForeeTx.CurStage,
			"cause", "transaction stage mismatch",
		)
		return fmt.Errorf("ForeeTx stage mismatch expect `%v` but `%v`", fTx.CurStage, curForeeTx.CurStage)
	}

	switch fTx.CurStage {
	case "":
		fTx.CurStage = transaction.TxStageInteracCI
		err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx)
		if err != nil {
			dbTx.Rollback()
			foree_logger.Logger.Error("processRootTx", "foreeTxId", fTx.ID, "cause", err.Error())
			return err
		}
		if err = dbTx.Commit(); err != nil {
			foree_logger.Logger.Error("processRootTx", "foreeTxId", fTx.ID, "cause", err.Error())
			return err
		}
		//TODO: go update summaryTx
		fallthrough
	case transaction.TxStageInteracCI:
		return nil
	default:
		foree_logger.Logger.Error("processRootTx",
			"foreeTxId", fTx.ID,
			"transactionStage", fTx.CurStage,
			"cause", "unkown foreeTx stage",
		)
		return fmt.Errorf("unknmow foreeTx stage `%v`", fTx.CurStage)
	}
}

func (p *TxProcessor) processTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {

	if tx.Type != transaction.TxTypeInteracToNBP {
		return nil, fmt.Errorf("unknow ForeeTx type `%s` for transaction `%v`", tx.Type, tx.ID)
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
			p.updateTxSummary(context.TODO(), *nTx)
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

func (p *TxProcessor) doProcessTx(ctx context.Context, fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if fTx.Status == transaction.TxStatusInitial {
		fTx.Status = transaction.TxStatusProcessing
		fTx.CurStage = transaction.TxStageInteracCI
		fTx.CurStageStatus = transaction.TxStatusInitial
		return &fTx, nil
	}
	if fTx.Status == transaction.TxStatusCompleted || fTx.Status == transaction.TxStatusCancelled || fTx.Status == transaction.TxStatusRejected {
		//TODO: log warn.
		return &fTx, nil
	}

	switch fTx.CurStage {
	case transaction.TxStageInteracCI:
		switch fTx.CurStageStatus {
		case transaction.TxStatusInitial:
			return p.ciTxProcessor.processTx(fTx)
		case transaction.TxStatusSent:
			return p.ciTxProcessor.waitFTx(fTx)
		case transaction.TxStatusCompleted:
			fTx.CurStage = transaction.TxStageIDM
			fTx.CurStageStatus = transaction.TxStatusInitial
			return &fTx, nil
		case transaction.TxStatusRejected:
			fTx.Status = transaction.TxStatusRejected
			fTx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", fTx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx); err != nil {
				return nil, err
			}
			// Close remaing non-terminated transactions.
			nT, err := p.closeRemainingTx(ctx, fTx)
			if err != nil {
				return nil, err
			}
			go p.MaybeRefund(ctx, *nT)
			return nT, nil
		case transaction.TxStatusCancelled:
			fTx.Status = transaction.TxStatusCancelled
			fTx.Conclusion = fmt.Sprintf("Cancelled in `%s` at %s", fTx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx); err != nil {
				return nil, err
			}
			// Close remaing non-terminated transactions.
			nT, err := p.closeRemainingTx(ctx, fTx)
			if err != nil {
				return nil, err
			}
			go p.MaybeRefund(ctx, *nT)
			return nT, nil
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", fTx.ID, fTx.CurStageStatus, fTx.CurStage)
		}
	case transaction.TxStageIDM:
		switch fTx.CurStageStatus {
		case transaction.TxStatusInitial:
			return p.idmTxProcessor.processTx(fTx)
		case transaction.TxStatusCompleted:
			//Move to next stage
			fTx.CurStage = transaction.TxStageNBPCO
			fTx.CurStageStatus = transaction.TxStatusInitial
			return &fTx, nil
		case transaction.TxStatusRejected:
			// Set to ForeeTx to terminal status.
			fTx.Status = transaction.TxStatusRejected
			fTx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", fTx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx); err != nil {
				return nil, err
			}
			// Close remaing non-terminated transactions.
			nT, err := p.closeRemainingTx(ctx, fTx)
			if err != nil {
				return nil, err
			}
			go p.MaybeRefund(ctx, *nT)
			return nT, nil
		case transaction.TxStatusSuspend:
			//Wait to approve
			//Log warn?
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", fTx.ID, fTx.CurStageStatus, fTx.CurStage)
		}
	case transaction.TxStageNBPCO:
		switch fTx.CurStageStatus {
		case transaction.TxStatusInitial:
			return p.nbpTxProcessor.processTx(fTx)
		case transaction.TxStatusSent:
			return p.nbpTxProcessor.waitFTx(fTx)
		case transaction.TxStatusCompleted:
			fTx.Status = transaction.TxStatusCompleted
			fTx.Conclusion = fmt.Sprintf("Complete at %s.", time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx); err != nil {
				return nil, err
			}
			go p.RedeemReward(ctx, fTx)
			return &fTx, nil
		case transaction.TxStatusRejected:
			fTx.Status = transaction.TxStatusRejected
			fTx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", fTx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx); err != nil {
				return nil, err
			}
			go p.MaybeRefund(ctx, fTx)
			return &fTx, nil
		case transaction.TxStatusCancelled:
			fTx.Status = transaction.TxStatusCancelled
			fTx.Conclusion = fmt.Sprintf("Rejected in `%s` at %s", fTx.CurStage, time_util.NowInToronto().Format(time.RFC3339))
			if err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx); err != nil {
				return nil, err
			}
			go p.MaybeRefund(ctx, fTx)
			return &fTx, nil
		default:
			return nil, fmt.Errorf("transaction `%v` in unknown status `%s` at statge `%s`", fTx.ID, fTx.CurStageStatus, fTx.CurStage)
		}
	default:
		return nil, fmt.Errorf("transaction `%v` in unknown stage `%s`", fTx.ID, fTx.CurStage)
	}
	return &fTx, nil
}

func (p *TxProcessor) closeRemainingTx(ctx context.Context, fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	switch fTx.CurStage {
	case transaction.TxStageInteracCI:
		idm := fTx.IDM
		co := fTx.COUT
		idm.Status = transaction.TxStatusClosed
		co.Status = transaction.TxStatusClosed
		if err := p.idmTxRepo.UpdateIDMTxById(ctx, *idm); err != nil {
			return nil, err
		}
		if err := p.nbpTxRepo.UpdateNBPCOTxById(ctx, *co); err != nil {
			return nil, err
		}
		return &fTx, nil
	case transaction.TxStageIDM:
		co := fTx.COUT
		co.Status = transaction.TxStatusClosed
		if err := p.nbpTxRepo.UpdateNBPCOTxById(ctx, *co); err != nil {
			return nil, err
		}
		return &fTx, nil
	default:
		//TODO: Log warn
		return &fTx, nil
	}
}

// TODO: reDesign.
func (p *TxProcessor) updateTxSummary(ctx context.Context, fTx transaction.ForeeTx) {
	txSummary := *fTx.Summary
	txSummary.IsCancelAllowed = false

	if fTx.Status == transaction.TxStatusInitial {
		txSummary.Status = transaction.TxSummaryStatusInitial
	} else if fTx.Status == transaction.TxStatusProcessing {
		if fTx.CurStage == transaction.TxStageInteracCI && fTx.CurStageStatus == transaction.TxStatusSent {
			txSummary.Status = transaction.TxSummaryStatusAwaitPayment
			txSummary.IsCancelAllowed = true
		} else if fTx.CurStage == transaction.TxStageNBPCO && fTx.CurStageStatus == transaction.TxStatusSent && fTx.COUT.CashOutAcc.Type == foree_constant.ContactAccountTypeCash {
			txSummary.Status = transaction.TxSummaryStatusPickup
			txSummary.IsCancelAllowed = true
		} else {
			txSummary.Status = transaction.TxSummaryStatusInProgress
		}
	} else if fTx.Status == transaction.TxStatusCompleted {
		txSummary.Status = transaction.TxSummaryStatusCompleted
	} else if fTx.Status == transaction.TxStatusCancelled || fTx.Status == transaction.TxStatusRejected {
		//TODO: check refund.
		txSummary.Status = transaction.TxSummaryStatusCancelled
	} else {
		//TODO: log error
		return
	}

	if txSummary.Status != fTx.Summary.Status {
		err := p.txSummaryRepo.UpdateTxSummaryById(ctx, txSummary)
		if err != nil {
			//TODO: log
			return
		}
	}

}

func (p *TxProcessor) recordTxHistory(h transaction.TxHistory) {
	if _, err := p.txHistoryRepo.InserTxHistory(context.Background(), h); err != nil {
		fmt.Println(err.Error())
	}
}

func (p *TxProcessor) RedeemReward(ctx context.Context, fTx transaction.ForeeTx) {
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return
	}

	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)
	rewards, err := p.rewardRepo.GetAllRewardByAppliedTransactionId(ctx, fTx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: Log error
		return
	}

	for _, v := range rewards {
		v.Status = transaction.RewardStatusRedeemed
		err := p.rewardRepo.UpdateRewardTxById(ctx, *v)
		if err != nil {
			dTx.Rollback()
			//TODO: Log error
			return
		}
	}

	if err = dTx.Commit(); err != nil {
		//TODO: Log error
		return
	}
}

func (p *TxProcessor) MaybeRefund(ctx context.Context, fTx transaction.ForeeTx) {
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return
	}

	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	foreeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return
	}

	if foreeTx.Status != transaction.TxStatusCancelled && foreeTx.Status != transaction.TxStatusRejected {
		dTx.Rollback()
		//TODO: log err
		return
	}

	if foreeTx.CurStage == transaction.TxStageRefund {
		dTx.Rollback()
		//TODO: double refund.
		return
	}

	rewards, err := p.rewardRepo.GetAllRewardByAppliedTransactionId(ctx, fTx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: Log error
		return
	}

	// Refund rewards.
	for _, v := range rewards {
		v.Status = transaction.RewardStatusDelete
		err := p.rewardRepo.UpdateRewardTxById(ctx, *v)
		if err != nil {
			dTx.Rollback()
			//TODO: Log error
			return
		}
		v.Status = transaction.RewardStatusActive
		_, err = p.rewardRepo.InsertReward(ctx, *v)
		if err != nil {
			dTx.Rollback()
			//TODO: Log error
			return
		}
	}

	// Refund limit.
	reference := transaction.GetDailyTxLimitReference(&fTx)
	dailyLimit, err := p.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
	if err != nil {
		dTx.Rollback()
		//TODO: Log error
		return
	}

	dailyLimit.UsedAmt.Amount += fTx.SrcAmt.Amount

	if err := p.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		dTx.Rollback()
		//TODO: Log error
		return
	}

	// Create refund transaction
	if fTx.CI.Status == transaction.TxStatusCompleted {
		_, err := p.interacRefundTxRepo.InsertInteracRefundTx(ctx, transaction.InteracRefundTx{
			Status:             transaction.RefundTxStatusInitial,
			RefundInteracAccId: fTx.CI.ID,
			ParentTxId:         fTx.ID,
			OwnerId:            fTx.OwnerId,
		})
		if err != nil {
			dTx.Rollback()
			//TODO: Log error
			return
		}
	}

	fTx.CurStage = transaction.TxStageRefund

	if err := p.foreeTxRepo.UpdateForeeTxById(ctx, fTx); err != nil {
		dTx.Rollback()
		//TODO: Log error
		return
	}

	if err = dTx.Commit(); err != nil {
		//TODO: Log error
		return
	}
}
