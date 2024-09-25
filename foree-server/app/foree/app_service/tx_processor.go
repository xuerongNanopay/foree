package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"xue.io/go-pay/app/foree/account"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
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
	interacRefundTxRepo *transaction.ForeeRefundTxRepo,
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

func (p *TxProcessor) SetInteracTxProcessor(interacTxProcessor *InteracTxProcessor) {
	p.interacTxProcessor = interacTxProcessor
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
	foreeRefundRepo     *transaction.ForeeRefundTxRepo
	interacRefundTxRepo *transaction.ForeeRefundTxRepo
	rewardRepo          *transaction.RewardRepo
	dailyTxLimiteRepo   *transaction.DailyTxLimitRepo
	userRepo            *auth.UserRepo
	contactAccountRepo  *account.ContactAccountRepo
	interacAccountRepo  *account.InteracAccountRepo
	interacTxProcessor  *InteracTxProcessor
	idmTxProcessor      *IDMTxProcessor
	nbpTxProcessor      *NBPTxProcessor
}

func (p *TxProcessor) createAndProcessTx(tx transaction.ForeeTx) {
	// foreeTx, err := p.createFullTx(tx)
	// if err != nil {
	// 	foree_logger.Logger.Error("CreateAndProcessTx_FAIL",
	// 		"foreeTxId", tx.ID,
	// 		"cause", err.Error(),
	// 	)
	// 	return
	// }

	// _, err = p.processTx(*foreeTx)
	// if err != nil {
	// 	foree_logger.Logger.Error("CreateAndProcessTx_FAIL",
	// 		"foreeTxId", tx.ID,
	// 		"cause", err.Error(),
	// 	)
	// }
}

func (p *TxProcessor) loadAndProcessTx(foreeId int64) (*transaction.ForeeTx, error) {
	fTx, err := p.loadTx(foreeId, true)
	if err != nil {
		return nil, err
	}

	go func() {
		// _, err := p.processTx(*fTx)
		// if err != nil {
		// 	//TODO log
		// }
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
		foree_logger.Logger.Error("CreateFullTx_FAIL",
			"ip", loadRealIp(ctx),
			"foreeTxId", tx.ID,
			"cause", ciErr.Error(),
		)
		return nil, ciErr
	}
	if idmErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_FAIL",
			"ip", loadRealIp(ctx),
			"foreeTxId", tx.ID,
			"cause", idmErr.Error(),
		)
		return nil, idmErr
	}
	if coutErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_FAIL",
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
		foree_logger.Logger.Error("CreateFullTx_FAIL",
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
		foree_logger.Logger.Error("loadTx_FAIL", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if foreeTx == nil {
		foree_logger.Logger.Warn("loadTx_FAIL",
			"foreeTxId", id,
			"cause", "foreeTx no found",
		)
		return nil, fmt.Errorf("ForeeTx no found with id `%v`", id)
	}

	// Load CI
	ci, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, foreeTx.ID)
	if err != nil {
		foree_logger.Logger.Error("loadTx_FAIL", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && ci == nil {
		foree_logger.Logger.Warn("loadTx_FAIL",
			"foreeTxId", id,
			"cause", "InteracCITx no found",
		)
		return nil, fmt.Errorf("InteracCITx no found for ForeeTx `%v`", foreeTx.ID)
	}

	CashInAcc, err := p.interacAccountRepo.GetUniqueInteracAccountById(ctx, ci.CashInAccId)
	if err != nil {
		foree_logger.Logger.Error("loadTx_FAIL", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && CashInAcc == nil {
		foree_logger.Logger.Warn("loadTx_FAIL",
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
		foree_logger.Logger.Error("loadTx_FAIL", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && idm == nil {
		foree_logger.Logger.Warn("loadTx_FAIL",
			"foreeTxId", id,
			"cause", "idmTx no found",
		)
		return nil, fmt.Errorf("IDMTx no found for ForeeTx `%v`", foreeTx.ID)
	}
	foreeTx.IDM = idm

	// Load COUT
	cout, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, foreeTx.ID)
	if err != nil {
		foree_logger.Logger.Error("loadTx_FAIL", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && cout == nil {
		foree_logger.Logger.Warn("loadTx_FAIL",
			"foreeTxId", id,
			"cause", "nbpCOTx no found",
		)
		return nil, fmt.Errorf("NBPCOTx no found for ForeeTx `%v`", foreeTx.ID)
	}

	CashOutAcc, err := p.contactAccountRepo.GetUniqueContactAccountById(ctx, cout.CashOutAccId)
	if err != nil {
		foree_logger.Logger.Error("loadTx_FAIL", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && CashOutAcc == nil {
		foree_logger.Logger.Warn("loadTx_FAIL",
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
		foree_logger.Logger.Error("loadTx_FAIL", "foreeTxId", id, "cause", err.Error())
		return nil, err
	}
	if isEmptyCheck && user == nil {
		foree_logger.Logger.Warn("loadTx_FAIL",
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

// Stage: Begin/"" -> CI -> IDM -> COUT -> End
//
// Stage: Begin/"" -> CI -> IDM -> COUT -> Rollback -> End
//
// This is internal process.
// Yes, in theory the race condition exists, but unlikely to happen.
// To avoid race condition, the simple strategy is pull from DB when we need.
// Always use goroutine on this method.
func (p *TxProcessor) ProcessRootTx(fTxId int64) {
	ctx := context.TODO()
	fTx, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, fTxId)
	if err != nil {
		foree_logger.Logger.Error("TxProcessor--processRootTx", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	if fTx == nil {
		foree_logger.Logger.Error("TxProcessor--processRootTx",
			"foreeTxId", fTxId,
			"cause", "unknown ForeeTx",
		)
		return
	}

	switch fTx.Stage {
	case transaction.TxStageBegin:
		fTx.Stage = transaction.TxStageInteracCI
		err := p.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--processRootTx", "foreeTxId", fTx.ID, "cause", err.Error())
		}
		//TODO: go update summaryTx
		fallthrough
	case transaction.TxStageInteracCI:
		p.interacTxProcessor.process(fTxId)
	case transaction.TxStageIDM:
		p.idmTxProcessor.process(fTxId)
	case transaction.TxStageNBPCO:
		p.nbpTxProcessor.process(fTxId)
	case transaction.TxStageEnd:
		foree_logger.Logger.Warn("TxProcessor--processRootTx", "foreeTxId", fTx.ID, "cause", "process transaction that is in END stage already")
	default:
		foree_logger.Logger.Error("TxProcessor--processRootTx",
			"foreeTxId", fTx.ID,
			"transactionStage", fTx.Stage,
			"cause", "unkown foreeTx stage",
		)
	}
}

func (p *TxProcessor) next(fTxId int64) {
	ctx := context.TODO()
	fTx, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, fTxId)
	if err != nil {
		foree_logger.Logger.Error("TxProcessor--next_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	if fTx == nil {
		foree_logger.Logger.Error("TxProcessor--next_FAIL",
			"foreeTxId", fTxId,
			"cause", "unknown ForeeTx",
		)
		return
	}

	switch fTx.Stage {
	case transaction.TxStageInteracCI:
		interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, fTxId)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--next_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}
		if interacTx.Status != transaction.TxStatusCompleted {
			foree_logger.Logger.Error(
				"TxProcessor--next_FAIL",
				"foreeTxId", fTxId,
				"curState", fTx.Stage,
				"interacTxStatus", interacTx.Status,
				"cause", "interacTx is not COMPLETED",
			)
			return
		}
		fTx.Stage = transaction.TxStageIDM
	case transaction.TxStageIDM:
		idmTx, err := p.idmTxRepo.GetUniqueIDMTxByParentTxId(ctx, fTxId)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--next_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}
		if idmTx.Status != transaction.TxStatusCompleted {
			foree_logger.Logger.Error(
				"TxProcessor--next_FAIL",
				"foreeTxId", fTxId,
				"curState", fTx.Stage,
				"idmTxStatus", idmTx.Status,
				"cause", "idmTx is not COMPLETED",
			)
			return
		}
		fTx.Stage = transaction.TxStageIDM
	case transaction.TxStageNBPCO:
		nbpTx, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, fTxId)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--next_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}
		if nbpTx.Status != transaction.TxStatusCompleted {
			foree_logger.Logger.Error(
				"TxProcessor--next_FAIL",
				"foreeTxId", fTxId,
				"curState", fTx.Stage,
				"nbpTxStatus", nbpTx.Status,
				"cause", "nbpTx is not COMPLETED",
			)
			return
		}
		fTx.Stage = transaction.TxStageEnd
	default:
		foree_logger.Logger.Error("TxProcessor--next_FAIL",
			"curStage", fTx.Stage,
		)
		return
	}

	err = p.foreeTxRepo.UpdateForeeTxById(context.TODO(), *fTx)
	if err != nil {
		foree_logger.Logger.Error("TxProcessor--next_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	p.ProcessRootTx(fTxId)
}

// If stage can not process transaction, then it will call rollback to rolling back the transaction.
func (p *TxProcessor) rollback(fTxId int64) {
	ctx := context.TODO()
	fTx, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, fTxId)
	if err != nil {
		foree_logger.Logger.Error("tx_processor-rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	if fTx == nil {
		foree_logger.Logger.Warn("tx_processor-rollback_FAIL",
			"foreeTxId", fTxId,
			"cause", "unknown ForeeTx",
		)
		return
	}

	if fTx.Stage == transaction.TxStageBegin {
		goto NO_Refund
	}

	if fTx.Stage == transaction.TxStageInteracCI {
		interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, fTxId)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}

		if interacTx.Status == transaction.TxStatusInitial || interacTx.Status == transaction.TxStatusCancelled {
			goto NO_Refund
		}
	}

	_, err = p.foreeRefundRepo.InsertForeeRefundTx(context.TODO(), transaction.ForeeRefundTx{
		Status:     transaction.RefundTxStatusInitial,
		RefundAmt:  fTx.TotalAmt,
		ParentTxId: fTx.ID,
		OwnerId:    fTx.OwnerId,
	})
	if err != nil {
		foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	fTx.Stage = transaction.TxStageRefund
	err = p.foreeTxRepo.UpdateForeeTxById(context.TODO(), *fTx)
	if err != nil {
		foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	return

NO_Refund:
	// No refund need.
	fTx.Stage = transaction.TxStageEnd
	err = p.foreeTxRepo.UpdateForeeTxById(context.TODO(), *fTx)
	if err != nil {
		foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
}

func (p *TxProcessor) onStatusUpdate(fTxId int64) {
	ctx := context.TODO()
	fTx, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, fTxId)
	if err != nil {
		foree_logger.Logger.Error("tx_processor-onStatusUpdate_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	if fTx == nil {
		foree_logger.Logger.Warn("tx_processor-onStatusUpdate_FAIL",
			"foreeTxId", fTxId,
			"cause", "unknown ForeeTx",
		)
		return
	}
	var newSummaryStatus transaction.TxSummaryStatus
	if fTx.Stage == transaction.TxStageBegin {
		newSummaryStatus = transaction.TxSummaryStatusInitial
	}
	if fTx.Stage == transaction.TxStageInteracCI {
		newSummaryStatus = transaction.TxSummaryStatusAwaitPayment
	}
	if fTx.Stage == transaction.TxStageIDM {
		newSummaryStatus = transaction.TxSummaryStatusInProgress
	}
	if fTx.Stage == transaction.TxStageNBPCO {
		//Specia case.
	}

	if newSummaryStatus != "TODO" {

	}
}

// TODO: reDesign.
func (p *TxProcessor) updateTxSummary(ctx context.Context, fTx transaction.ForeeTx) {
	// txSummary := *fTx.Summary
	// txSummary.IsCancelAllowed = false

	// if fTx.Status == transaction.TxStatusInitial {
	// 	txSummary.Status = transaction.TxSummaryStatusInitial
	// } else if fTx.Status == transaction.TxStatusProcessing {
	// 	if fTx.Stage == transaction.TxStageInteracCI && fTx.StageStatus == transaction.TxStatusSent {
	// 		txSummary.Status = transaction.TxSummaryStatusAwaitPayment
	// 		txSummary.IsCancelAllowed = true
	// 	} else if fTx.Stage == transaction.TxStageNBPCO && fTx.StageStatus == transaction.TxStatusSent && fTx.COUT.CashOutAcc.Type == foree_constant.ContactAccountTypeCash {
	// 		txSummary.Status = transaction.TxSummaryStatusPickup
	// 		txSummary.IsCancelAllowed = true
	// 	} else {
	// 		txSummary.Status = transaction.TxSummaryStatusInProgress
	// 	}
	// } else if fTx.Status == transaction.TxStatusCompleted {
	// 	txSummary.Status = transaction.TxSummaryStatusCompleted
	// } else if fTx.Status == transaction.TxStatusCancelled || fTx.Status == transaction.TxStatusRejected {
	// 	//TODO: check refund.
	// 	txSummary.Status = transaction.TxSummaryStatusCancelled
	// } else {
	// 	//TODO: log error
	// 	return
	// }

	// if txSummary.Status != fTx.Summary.Status {
	// 	err := p.txSummaryRepo.UpdateTxSummaryById(ctx, txSummary)
	// 	if err != nil {
	// 		//TODO: log
	// 		return
	// 	}
	// }

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

// func (p *TxProcessor) MaybeRefund(ctx context.Context, fTx transaction.ForeeTx) {
// 	dTx, err := p.db.Begin()
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: log err
// 		return
// 	}

// 	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

// 	foreeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: log err
// 		return
// 	}

// 	if foreeTx.Status != transaction.TxStatusCancelled && foreeTx.Status != transaction.TxStatusRejected {
// 		dTx.Rollback()
// 		//TODO: log err
// 		return
// 	}

// 	if foreeTx.Stage == transaction.TxStageRefund {
// 		dTx.Rollback()
// 		//TODO: double refund.
// 		return
// 	}

// 	rewards, err := p.rewardRepo.GetAllRewardByAppliedTransactionId(ctx, fTx.ID)
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: Log error
// 		return
// 	}

// 	// Refund rewards.
// 	for _, v := range rewards {
// 		v.Status = transaction.RewardStatusDelete
// 		err := p.rewardRepo.UpdateRewardTxById(ctx, *v)
// 		if err != nil {
// 			dTx.Rollback()
// 			//TODO: Log error
// 			return
// 		}
// 		v.Status = transaction.RewardStatusActive
// 		_, err = p.rewardRepo.InsertReward(ctx, *v)
// 		if err != nil {
// 			dTx.Rollback()
// 			//TODO: Log error
// 			return
// 		}
// 	}

// 	// Refund limit.
// 	reference := transaction.GetDailyTxLimitReference(&fTx)
// 	dailyLimit, err := p.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
// 	if err != nil {
// 		dTx.Rollback()
// 		//TODO: Log error
// 		return
// 	}

// 	dailyLimit.UsedAmt.Amount += fTx.SrcAmt.Amount

// 	if err := p.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
// 		dTx.Rollback()
// 		//TODO: Log error
// 		return
// 	}

// 	// Create refund transaction
// 	if fTx.CI.Status == transaction.TxStatusCompleted {
// 		_, err := p.interacRefundTxRepo.InsertForeeRefundTx(ctx, transaction.ForeeRefundTx{
// 			Status:     transaction.RefundTxStatusInitial,
// 			ParentTxId: fTx.ID,
// 			OwnerId:    fTx.OwnerId,
// 		})
// 		if err != nil {
// 			dTx.Rollback()
// 			//TODO: Log error
// 			return
// 		}
// 	}
// }
