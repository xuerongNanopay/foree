package foree_service

import (
	"context"
	"database/sql"
	"fmt"

	"xue.io/go-pay/app/foree/account"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/partner/nbp"
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
	rewardRepo *promotion.RewardRepo,
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
	rewardRepo          *promotion.RewardRepo
	dailyTxLimiteRepo   *transaction.DailyTxLimitRepo
	userRepo            *auth.UserRepo
	contactAccountRepo  *account.ContactAccountRepo
	interacAccountRepo  *account.InteracAccountRepo
	interacTxProcessor  *InteracTxProcessor
	idmTxProcessor      *IDMTxProcessor
	nbpTxProcessor      *NBPTxProcessor
}

func (p *TxProcessor) reloadTransactions() {

}

func (p *TxProcessor) createAndProcessTx(tx transaction.ForeeTx) {
	foreeTx, err := p.createFullTx(tx)
	if err != nil {
		foree_logger.Logger.Error("CreateAndProcessTx_FAIL",
			"foreeTxId", tx.ID,
			"cause", err.Error(),
		)
		return
	}

	// _, err = p.loadAndProcessTx(foreeTx.ID)
	// if err != nil {
	// 	foree_logger.Logger.Error("CreateAndProcessTx_FAIL",
	// 		"foreeTxId", tx.ID,
	// 		"cause", err.Error(),
	// 	)
	// }
	p.ProcessRootTx(foreeTx.ID)
}

// Create CI, COUT, IDM for ForeeTx
func (p *TxProcessor) createFullTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	_, err = p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	// Create CI
	var ciTx *transaction.InteracCITx
	var ciErr error
	createCI := func() {
		ciId, err := p.interacTxRepo.InsertInteracCITx(ctx, transaction.InteracCITx{
			Status:      transaction.TxStatusInitial,
			CashInAccId: fTx.CinAccId,
			EndToEndId:  fTx.Summary.NBPReference,
			Amt:         fTx.SrcAmt,
			ParentTxId:  fTx.ID,
			OwnerId:     fTx.OwnerId,
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
	createCI()

	// Create IDM
	var idmTx *transaction.IDMTx
	var idmErr error
	createIDM := func() {
		idmId, err := p.idmTxRepo.InsertIDMTx(ctx, transaction.IDMTx{
			Status:     transaction.TxStatusInitial,
			Ip:         fTx.Ip,
			UserAgent:  fTx.UserAgent,
			ParentTxId: fTx.ID,
			OwnerId:    fTx.OwnerId,
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
	createIDM()

	// Create Cout
	var coutTx *transaction.NBPCOTx
	var coutErr error
	createCout := func() {
		coutAcc, err := p.contactAccountRepo.GetUniqueContactAccountById(ctx, fTx.CoutAccId)
		if err != nil {
			coutErr = err
			return
		}

		if coutAcc == nil {
			coutErr = fmt.Errorf("cash out account no found with id `%v`", fTx.CoutAccId)
			return
		}

		mode, err := mapNBPMode(coutAcc)
		if err != nil {
			coutErr = err
			return
		}

		coutId, err := p.nbpTxRepo.InsertNBPCOTx(ctx, transaction.NBPCOTx{
			Status:       transaction.TxStatusInitial,
			Mode:         mode,
			Amt:          fTx.DestAmt,
			NBPReference: fTx.Summary.NBPReference,
			CashOutAccId: fTx.CoutAccId,
			ParentTxId:   fTx.ID,
			OwnerId:      fTx.OwnerId,
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

	createCout()

	if ciErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_FAIL",
			"ip", loadRealIp(ctx),
			"foreeTxId", fTx.ID,
			"cause", ciErr.Error(),
		)
		return nil, ciErr
	}
	if idmErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_FAIL",
			"ip", loadRealIp(ctx),
			"foreeTxId", fTx.ID,
			"cause", idmErr.Error(),
		)
		return nil, idmErr
	}
	if coutErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateFullTx_FAIL",
			"ip", loadRealIp(ctx),
			"foreeTxId", fTx.ID,
			"cause", coutErr.Error(),
		)
		return nil, coutErr
	}

	fTx.CI = ciTx
	fTx.IDM = idmTx
	fTx.COUT = coutTx

	if err = dTx.Commit(); err != nil {
		foree_logger.Logger.Error("CreateFullTx_FAIL",
			"ip", loadRealIp(ctx),
			"foreeTxId", fTx.ID,
			"cause", err.Error(),
		)
		return nil, err
	}
	return &fTx, nil
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
	foree_logger.Logger.Debug("TxProcessor--processRootTx",
		"foreeTxId", fTxId,
		"stage", fTx.Stage,
		"msg", "transaction process",
	)
	switch fTx.Stage {
	case transaction.TxStageBegin:
		fTx.Stage = transaction.TxStageInteracCI
		err := p.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--processRootTx", "foreeTxId", fTx.ID, "cause", err.Error())
			return
		}
		//TODO: go update summaryTx
		fallthrough
	case transaction.TxStageInteracCI:
		p.interacTxProcessor.process(fTxId)
	case transaction.TxStageIDM:
		p.idmTxProcessor.process(fTxId)
	case transaction.TxStageNBPCO:
		p.nbpTxProcessor.process(fTxId)
	case transaction.TxStageRefunding:
		p.next(fTxId)
	case transaction.TxStageSuccess:
		foree_logger.Logger.Warn("TxProcessor--processRootTx", "foreeTxId", fTx.ID, "cause", "process transaction that is SUCCESS already")
	case transaction.TxStageCancel:
		foree_logger.Logger.Warn("TxProcessor--processRootTx", "foreeTxId", fTx.ID, "cause", "process transaction that is CANCEL already")
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
		fTx.Stage = transaction.TxStageNBPCO
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
		//TODO: close up rewards.
		go p.updateRewardsToComplete(fTx.ID)
		fTx.Stage = transaction.TxStageSuccess
	case transaction.TxStageRefunding:
		refundTx, err := p.foreeRefundRepo.GetUniqueForeeRefundTxByParentTxId(ctx, fTxId)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--next_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}
		if refundTx.Status == transaction.TxStatusCompleted || refundTx.Status == transaction.TxStatusClosed {
			fTx.Stage = transaction.TxStageCancel
		}
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

	if fTx.Stage != transaction.TxStageSuccess && fTx.Stage != transaction.TxStageCancel {
		p.ProcessRootTx(fTxId)
	} else {
		foree_logger.Logger.Info("TxProcessor", "foreeTxId", fTxId, "msg", "transaction terminate")
	}
}

// If stage can not process transaction, then it will call rollback to rolling back the transaction.
// Close remaining
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

	if fTx.Stage == transaction.TxStageRefunding || fTx.Stage == transaction.TxStageCancel || fTx.Stage == transaction.TxStageSuccess {
		foree_logger.Logger.Warn("tx_processor-rollback_FAIL",
			"foreeTxId", fTxId,
			"curStage", fTx.Stage,
			"cause", "transaction is in stage that can't rollback",
		)
		return
	}

	dbTx, err := p.db.Begin()
	if err != nil {
		foree_logger.Logger.Error("tx_processor-rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		dbTx.Rollback()
		return
	}
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dbTx)

	//Close remaining transaction.
	err = p.closeRemainingTx(ctx, fTxId)
	if err != nil {
		foree_logger.Logger.Error("tx_processor-rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		dbTx.Rollback()
		return
	}

	//Revert rewards and dayly limit.
	err = p.revertRewardAndTxLimit(ctx, *fTx)
	if err != nil {
		foree_logger.Logger.Error("tx_processor-rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		dbTx.Rollback()
		return
	}

	if fTx.Stage == transaction.TxStageBegin {
		goto NO_Refund
	}

	if fTx.Stage == transaction.TxStageInteracCI {
		interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, fTxId)
		if err != nil {
			foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			dbTx.Rollback()
			return
		}

		// The case that we can safely cancel transaction.
		if interacTx.Status == transaction.TxStatusInitial || interacTx.Status == transaction.TxStatusCancelled || interacTx.Status == transaction.TxStatusRejected {
			goto NO_Refund
		}
	}

	_, err = p.foreeRefundRepo.InsertForeeRefundTx(ctx, transaction.ForeeRefundTx{
		Status:     transaction.TxStatusInitial,
		RefundAmt:  fTx.TotalAmt,
		ParentTxId: fTx.ID,
		OwnerId:    fTx.OwnerId,
	})

	if err != nil {
		foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		dbTx.Rollback()
		return
	}
	fTx.Stage = transaction.TxStageRefunding
	err = p.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
	if err != nil {
		dbTx.Rollback()
		foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	goto COMMIT

NO_Refund:
	// No refund need.
	fTx.Stage = transaction.TxStageCancel
	err = p.foreeTxRepo.UpdateForeeTxById(ctx, *fTx)
	if err != nil {
		dbTx.Rollback()
		foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
COMMIT:
	if err = dbTx.Commit(); err != nil {
		foree_logger.Logger.Error("TxProcessor--rollback_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	go p.updateSummaryTx(fTxId)
}

func (p *TxProcessor) updateRewardsToComplete(fTxId int64) {

}

func (p *TxProcessor) closeRemainingTx(ctx context.Context, fTxId int64) error {
	interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, fTxId)
	if err != nil {
		return err
	}
	if interacTx == nil {
		return fmt.Errorf("interacTx no found with parentTxId `%v`", fTxId)
	}

	idmTx, err := p.idmTxRepo.GetUniqueIDMTxByParentTxId(ctx, fTxId)
	if err != nil {
		return err
	}
	if idmTx == nil {
		return fmt.Errorf("idmTx no found with parentTxId `%v`", fTxId)
	}

	nbpTx, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, fTxId)
	if err != nil {
		return err
	}
	if nbpTx == nil {
		return fmt.Errorf("nbpTx no found with parentTxId `%v`", fTxId)
	}

	if interacTx.Status == transaction.TxStatusInitial {
		newInteracTx := *interacTx
		newInteracTx.Status = transaction.TxStatusClosed
		if err := p.interacTxRepo.UpdateInteracCITxById(ctx, newInteracTx); err != nil {
			return err
		}
	}

	if idmTx.Status == transaction.TxStatusInitial {
		newInteracTx := *idmTx
		newInteracTx.Status = transaction.TxStatusClosed
		if err := p.idmTxRepo.UpdateIDMTxById(ctx, newInteracTx); err != nil {
			return err
		}
	}

	if nbpTx.Status == transaction.TxStatusInitial {
		newInteracTx := *nbpTx
		newInteracTx.Status = transaction.TxStatusClosed
		if err := p.nbpTxRepo.UpdateNBPCOTxById(ctx, newInteracTx); err != nil {
			return err
		}
	}
	return nil
}

func (p *TxProcessor) updateSummaryTx(fTxId int64) {
	ctx := context.TODO()
	fTx, err := p.foreeTxRepo.GetUniqueForeeTxById(ctx, fTxId)
	if err != nil {
		foree_logger.Logger.Error("tx_processor-updateSummaryTx_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	if fTx == nil {
		foree_logger.Logger.Warn("tx_processor-updateSummaryTx_FAIL",
			"foreeTxId", fTxId,
			"cause", "ForeeTx no found",
		)
		return
	}
	sumTx, err := p.txSummaryRepo.GetUniqueTxSummaryByParentTxId(context.TODO(), fTx.ID)
	if err != nil {
		foree_logger.Logger.Error("tx_processor-updateSummaryTx_FAIL", "foreeTxId", fTxId, "cause", err.Error())
		return
	}
	newSumTx := *sumTx
	if fTx.Stage == transaction.TxStageBegin {
		newSumTx.Status = transaction.TxSummaryStatusInitial
	} else if fTx.Stage == transaction.TxStageInteracCI {
		interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(context.TODO(), fTx.ID)
		if err != nil {
			foree_logger.Logger.Error("tx_processor-updateSummaryTx_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}
		if interacTx == nil {
			foree_logger.Logger.Warn("tx_processor-updateSummaryTx_FAIL",
				"foreeTxId", fTxId,
				"cause", "InteracTx no found",
			)
			return
		}
		if interacTx.Status == transaction.TxStatusInitial {
			newSumTx.Status = transaction.TxSummaryStatusInitial
		} else if interacTx.Status == transaction.TxStatusSent {
			newSumTx.Status = transaction.TxSummaryStatusAwaitPayment
			newSumTx.PaymentUrl = interacTx.PaymentUrl
			newSumTx.IsCancelAllowed = true
		} else {
			newSumTx.Status = transaction.TxSummaryStatusInProgress
		}
	} else if fTx.Stage == transaction.TxStageIDM {
		newSumTx.Status = transaction.TxSummaryStatusInProgress
	} else if fTx.Stage == transaction.TxStageNBPCO {
		nbpTx, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(context.TODO(), fTx.ID)
		if err != nil {
			foree_logger.Logger.Error("tx_processor-updateSummaryTx_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}
		if nbpTx.Mode == nbp.PMTModeCash && nbpTx.Status == transaction.TxStatusSent {
			newSumTx.Status = transaction.TxSummaryStatusPickup
			newSumTx.IsCancelAllowed = true
		} else {
			newSumTx.Status = transaction.TxSummaryStatusInProgress
		}
	} else if fTx.Stage == transaction.TxStageSuccess {
		newSumTx.Status = transaction.TxSummaryStatusCompleted
	} else if fTx.Stage == transaction.TxStageCancel {
		newSumTx.Status = transaction.TxSummaryStatusCancelled
	} else if fTx.Stage == transaction.TxStageRefunding {
		//TODO:
		newSumTx.Status = transaction.TxSummaryStatusRefunding
	} else {
		foree_logger.Logger.Error("TxProcessor--updateSummaryTx_FAIL", "stage", fTx.Stage, "cause", "unknow stage")
		return
	}

	if sumTx.Status != newSumTx.Status || newSumTx.IsCancelAllowed != sumTx.IsCancelAllowed {
		err = p.txSummaryRepo.UpdateTxSummaryById(context.TODO(), newSumTx)
		if err != nil {
			foree_logger.Logger.Error("tx_processor-updateSummaryTx_FAIL", "foreeTxId", fTxId, "cause", err.Error())
			return
		}

		foree_logger.Logger.Debug("tx_processor-updateSummaryTx_Success", "oldSummaryStatus", sumTx.Status, "newSummaryStatus", newSumTx.Status)
	}
}

func (p *TxProcessor) revertRewardAndTxLimit(ctx context.Context, fTx transaction.ForeeTx) error {
	//Refund Reward.
	rewards, err := p.rewardRepo.GetAllRewardByAppliedTransactionId(ctx, fTx.ID)
	if err != nil {
		return err
	}

	//TODO: check the reward type for non_refundable rewards.
	for _, reward := range rewards {
		r := *reward
		r.Status = nbp.AccStatusActive
		r.AppliedTransactionId = 0
		if err := p.rewardRepo.UpdateRewardTxById(ctx, r); err != nil {
			return err
		}
	}

	//Refund daily limit.
	dailyLimit, err := p.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, fTx.LimitReference)
	if err != nil {
		return err
	}

	if dailyLimit == nil {
		return nil
	}

	dailyLimit.UsedAmt.Amount -= fTx.SrcAmt.Amount

	if err := p.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return err
	}
	foree_logger.Logger.Info("Refund Transaction limit", "foreeTxId", fTx.ID, "refundAmout", fTx.SrcAmt.Amount)
	return nil
}
