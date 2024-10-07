package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/partner/nbp"
	string_util "xue.io/go-pay/util/string"
)

const nbpTxRecheckInterval = 10 * time.Minute
const nbpTxRetryInterval = 10 * time.Second
const nbpTxRetryAttempts = 5
const nbpTxStatusRefreshTicker = 2 * time.Minute

func NewNBPTxProcessor(
	db *sql.DB,
	foreeTxRepo *transaction.ForeeTxRepo,
	txProcessor *TxProcessor,
	nbpTxRepo *transaction.NBPCOTxRepo,
	nbpClient nbp.NBPClient,
	userExtraRepo *foree_auth.UserExtraRepo,
	userIdentificationRepo *foree_auth.UserIdentificationRepo,
) *NBPTxProcessor {
	ret := &NBPTxProcessor{
		db:                     db,
		foreeTxRepo:            foreeTxRepo,
		txProcessor:            txProcessor,
		nbpTxRepo:              nbpTxRepo,
		nbpClient:              nbpClient,
		userExtraRepo:          userExtraRepo,
		userIdentificationRepo: userIdentificationRepo,
		statusRefreshChan:      make(chan transaction.NBPCOTx, 64),
		statusRefreshTicker:    time.NewTicker(nbpTxStatusRefreshTicker),
	}
	go ret.start()
	return ret
}

type NBPTxWaitWrapper struct {
	nbpTx     transaction.NBPCOTx
	recheckAt time.Time
}

type NBPTxProcessor struct {
	db                     *sql.DB
	foreeTxRepo            *transaction.ForeeTxRepo
	txProcessor            *TxProcessor
	nbpTxRepo              *transaction.NBPCOTxRepo
	nbpClient              nbp.NBPClient
	userExtraRepo          *foree_auth.UserExtraRepo
	userIdentificationRepo *foree_auth.UserIdentificationRepo
	waits                  sync.Map
	statusRefreshChan      chan transaction.NBPCOTx
	statusRefreshTicker    *time.Ticker
}

func (p *NBPTxProcessor) start() {
	for {
		select {
		case nbpTx := <-p.statusRefreshChan:
			_, ok := p.waits.Load(nbpTx.NBPReference)
			if ok {
				foree_logger.Logger.Warn("NBPTxProcessor-statusRefreshChan",
					"nbpReference", nbpTx.NBPReference,
					"msg", "nbpTx is in waiting aleardy",
				)
				continue
			}
			p.waits.Store(nbpTx.NBPReference, NBPTxWaitWrapper{
				nbpTx:     nbpTx,
				recheckAt: time.Now().Add(nbpTxRecheckInterval),
			})
		case <-p.statusRefreshTicker.C:
			waitNBPReferences := make([]string, 0)
			p.waits.Range(func(k, _ interface{}) bool {
				nbpReference, _ := k.(string)
				waitNBPReferences = append(waitNBPReferences, nbpReference)
				return true
			})
			chunks := string_util.ChunkSlice(waitNBPReferences, 8)
			for _, c := range chunks {
				p.refreshNBPStatuses(c)
			}
		}
	}
}

func (p *NBPTxProcessor) process(parentTxId int64) {
	ctx := context.TODO()
	nbpTx, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, parentTxId)
	if err != nil {
		foree_logger.Logger.Error("NBPTxProcessor-process", "parentTxId", parentTxId, "cause", err.Error())
		return
	}
	if nbpTx == nil {
		foree_logger.Logger.Error("NBPTxProcessor-process", "parentTxId", parentTxId, "cause", "interacTx no found")
		return
	}

	switch nbpTx.Status {
	case transaction.TxStatusInitial:
		p.loadRemittance(nbpTx.ParentTxId)
	case transaction.TxStatusSent:
		p.statusRefreshChan <- *nbpTx
		go p.txProcessor.updateSummaryTx(nbpTx.ParentTxId)
	case transaction.TxStatusCompleted:
		p.txProcessor.next(nbpTx.ParentTxId)
		go p.txProcessor.updateSummaryTx(nbpTx.ParentTxId)
	case transaction.TxStatusRejected:
		p.txProcessor.rollback(nbpTx.ParentTxId)
	case transaction.TxStatusCancelled:
		p.txProcessor.rollback(nbpTx.ParentTxId)
	default:
		foree_logger.Logger.Error(
			"NBPTxProcessor-process",
			"parentTxId", parentTxId,
			"nbpCITxId", nbpTx.ID,
			"nbpCITxStatus", nbpTx.Status,
			"cause", "unsupport status",
		)
	}
}

func (p *NBPTxProcessor) loadRemittance(parentTxId int64) {
	fTx, err := p.txProcessor.loadTx(parentTxId, true)
	if err != nil {
		foree_logger.Logger.Error("NBPTxProcessor--loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
		return
	}
	req, err := p.buildLoadRemittanceRequest(fTx)
	if err != nil {
		foree_logger.Logger.Error("NBPTxProcessor--loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
		return
	}

	mode, err := mapNBPMode(fTx.COUT.CashOutAcc)
	if err != nil {
		foree_logger.Logger.Error("NBPTxProcessor--loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
		return
	}

	//403: Internal Error
	//405: Duplicate Global ID
	//406: Remittance could not be load
	//401: token error

	// Retry 5 times with 15 second interval.
	var resp *nbp.LoadRemittanceResponse
	for i := 0; i < nbpTxRetryAttempts; i++ {
		resp, err = p.sendPaymentWithMode(*req, mode)
		if err != nil {
			foree_logger.Logger.Error("NBPTxProcessor--loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
			return
		}
		//Retry case: 5xx, 401, 403
		if resp.StatusCode/100 == 5 || resp.ResponseCode == "401" || resp.ResponseCode == "403" || resp.ResponseCode == "406" {
			foree_logger.Logger.Warn("NBPTxProcessor--loadRemittance", "parentTxId", parentTxId, "retry", i)
			time.Sleep(nbpTxRetryInterval)
		} else {
			break
		}

	}

	// Retry later manully
	if resp.StatusCode/100 == 5 || resp.ResponseCode == "401" || resp.ResponseCode == "403" {
		foree_logger.Logger.Warn("NBPTxProcessor--loadRemittance_FAIL",
			"parentTxId", parentTxId,
			"httpStatus", resp.StatusCode,
			"httpResponse", resp.RawResponse,
			"msg", "please re-run the transaction.",
		)
		return
	}

	nbpTx := *fTx.COUT

	// Success
	if resp.StatusCode/100 == 2 || resp.ResponseCode == "405" {
		nbpTx.Status = transaction.TxStatusSent
		err := p.nbpTxRepo.UpdateNBPCOTxById(context.TODO(), nbpTx)
		if err != nil {
			foree_logger.Logger.Error("NBPTxProcessor--loadRemittance_FAIL",
				"parentTxId", parentTxId,
				"httpStatus", resp.StatusCode,
				"cause", err.Error(),
			)
			return
		}
		foree_logger.Logger.Info("NBPTxProcessor--loadRemittance_SUCCESS", "parentTxId", parentTxId)
		p.statusRefreshChan <- nbpTx
		return
	}

	// Reject
	foree_logger.Logger.Warn("NBPTxProcessor--loadRemittance_FAIL",
		"parentTxId", parentTxId,
		"httpStatus", resp.StatusCode,
		"httpRequest", resp.RawRequest,
		"httpResponse", resp.RawResponse,
		"msg", "nbp call failed",
	)
	nbpTx.Status = transaction.TxStatusRejected
	err = p.nbpTxRepo.UpdateNBPCOTxById(context.TODO(), nbpTx)
	if err != nil {
		foree_logger.Logger.Error("NBPTxProcessor--loadRemittance_FAIL",
			"parentTxId", parentTxId,
			"httpStatus", resp.StatusCode,
			"cause", err.Error(),
		)
		return
	}
	p.process(nbpTx.ParentTxId)
}

func (p *NBPTxProcessor) refreshNBPStatuses(nbpReferences []string) {
	resp, err := p.nbpClient.TransactionStatusByIds(nbp.TransactionStatusByIdsRequest{
		Ids: nbpReferences,
	})

	if err != nil {
		foree_logger.Logger.Error("NBPTxProcessor--refreshNBPStatuses_FAIL",
			"cause", err.Error(),
		)
		return
	}

	if resp.StatusCode/100 != 2 {
		foree_logger.Logger.Error("NBPTxProcessor--refreshNBPStatuses_FAIL",
			"httpStatus", resp.StatusCode,
			"httpRawRequest", resp.RawRequest,
			"httpRawResponse", resp.RawResponse,
		)
		return
	}

	for _, nbpRef := range resp.TransactionStatuses {
		newTxStatus := nbpToInternalStatusMapper(nbpRef.Status)
		if newTxStatus == transaction.TxStatusSent {
			continue
		}

		curNBPTx, err := p.nbpTxRepo.GetUniqueNBPCOTxByNBPReference(context.TODO(), nbpRef.GlobalId)
		if err != nil {
			foree_logger.Logger.Error("NBPTxProcessor--refreshNBPStatuses_FAIL",
				"cause", err.Error(),
			)
			continue
		}

		if curNBPTx.Status != transaction.TxStatusSent {
			foree_logger.Logger.Warn("NBPTxProcessor--refreshNBPStatuses_FAIL",
				"nbpTxId", curNBPTx.ID,
				"nbpReference", curNBPTx.NBPReference,
				"currentNbpTxStatus", curNBPTx.Status,
				"cause", "try to upddate nbpTx that is not in SENT status",
			)
			p.waits.Delete(curNBPTx.NBPReference)
			continue
		}

		curNBPTx.Status = newTxStatus
		err = p.nbpTxRepo.UpdateNBPCOTxById(context.TODO(), *curNBPTx)
		if err != nil {
			foree_logger.Logger.Error("NBPTxProcessor--refreshNBPStatuses_FAIL",
				"cause", err.Error(),
			)
			continue
		}
		p.waits.Delete(curNBPTx.NBPReference)
		//moving forward
		go p.process(curNBPTx.ParentTxId)
		foree_logger.Logger.Info("NBPTxProcessor--refreshNBPStatuses",
			"nbpTxId", curNBPTx.ID,
			"nbpReference", curNBPTx.NBPReference,
			"newNBPTxStatus", curNBPTx.Status,
			"msg", "NBP wait complete",
		)
	}
}

func (p *NBPTxProcessor) ManualUpdate(parentTxId int64, newTxStatus transaction.TxStatus) (bool, error) {
	if newTxStatus != transaction.TxStatusRejected && newTxStatus != transaction.TxStatusCompleted {
		return false, fmt.Errorf("unsupport transaction status `%v`", newTxStatus)
	}

	ctx := context.TODO()
	nbpTx, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, parentTxId)
	if err != nil {
		return false, err
	}
	if nbpTx == nil {
		return false, fmt.Errorf("InteracCITx no found with parentTxId `%v`", parentTxId)
	}
	if nbpTx.Status != transaction.TxStatusSent {
		return false, fmt.Errorf("expect InteracCITx in `%v`, but got `%v`", transaction.TxStatusSent, nbpTx.Status)
	}

	nbpTx.Status = transaction.TxStatusCompleted
	err = p.nbpTxRepo.UpdateNBPCOTxById(context.TODO(), *nbpTx)
	if err != nil {
		return false, err
	}

	p.waits.Delete(nbpTx.NBPReference)
	go p.txProcessor.next(nbpTx.ParentTxId)
	return true, nil
}

// TODO: call scotial cancel api
func (p *NBPTxProcessor) Cancel(parentTxId int64) (bool, error) {
	ctx := context.TODO()
	nbpTx, err := p.nbpTxRepo.GetUniqueNBPCOTxByParentTxId(ctx, parentTxId)
	if err != nil {
		return false, err
	}
	if nbpTx == nil {
		return false, fmt.Errorf("nbpTx no found with parentTxId `%v`", parentTxId)
	}
	if nbpTx.Status != transaction.TxStatusSent {
		return false, fmt.Errorf("expect nbpTx in `%v`, but got `%v`", transaction.TxStatusSent, nbpTx.Status)
	}

	//TODO: call scotial cancel api
	//TODO: if error return.

	nbpTx.Status = transaction.TxStatusCancelled
	err = p.nbpTxRepo.UpdateNBPCOTxById(context.TODO(), *nbpTx)
	if err != nil {
		return false, err
	}
	p.waits.Delete(nbpTx.NBPReference)
	go p.txProcessor.rollback(nbpTx.ParentTxId)
	return true, nil
}

func (p *NBPTxProcessor) sendPaymentWithMode(r nbp.LoadRemittanceRequest, mode nbp.PMTMode) (*nbp.LoadRemittanceResponse, error) {
	switch mode {
	case nbp.PMTModeCash:
		return p.nbpClient.LoadRemittanceCash(r)
	case nbp.PMTModeThirdPartyPayments:
		return p.nbpClient.LoadRemittanceThirdParty(r)
	case nbp.PMTModeAccountTransfers:
		return p.nbpClient.LoadRemittanceAccounts(r)
	default:
		return nil, fmt.Errorf("NBPTxProcessor -- unknow mode %s", mode)
	}
}

func (p *NBPTxProcessor) buildLoadRemittanceRequest(fTx *transaction.ForeeTx) (*nbp.LoadRemittanceRequest, error) {
	userExtra, err := p.userExtraRepo.GetUniqueUserExtraByOwnerId(fTx.OwnerId)
	if err != nil {
		return nil, err
	}
	identifications, err := p.userIdentificationRepo.GetAllUserIdentificationByOwnerId(fTx.OwnerId)
	if err != nil {
		return nil, err
	}

	var identification *foree_auth.UserIdentification = nil

	for _, v := range identifications {
		if v.Status == foree_auth.IdentificationStatusActive {
			identification = v
		}
	}

	if identification == nil {
		return nil, fmt.Errorf("NBPTxProcessor -- user `%v` do not find a proper identification", fTx.OwnerId)
	}

	transactionDate := time.Now()
	lrr := &nbp.LoadRemittanceRequest{
		GlobalId:                        fTx.COUT.NBPReference,
		Amount:                          nbp.NBPAmount(fTx.COUT.Amt.Amount),
		Currency:                        fTx.COUT.Amt.Currency,
		TransactionDate:                 (*nbp.NBPDate)(&transactionDate),
		OriginatingCountry:              "Canada",
		PurposeRemittance:               fTx.TransactionPurpose,
		RemitterName:                    fTx.CI.CashInAcc.GetLegalName(),
		RemitterEmail:                   fTx.CI.CashInAcc.Email,
		RemitterContact:                 fTx.Owner.PhoneNumber,
		RemitterDOB:                     (*nbp.NBPDate)(fTx.Owner.Dob),
		RemitterAddress:                 generateLoadRemittanceFromInteracAccount(fTx.Owner),
		RemitterIdType:                  mapNBPRemitterIdType(identification.Type),
		RemitterId:                      identification.Value,
		RemitterPOB:                     userExtra.Pob,
		BeneficiaryName:                 fTx.COUT.CashOutAcc.GetLegalName(),
		BeneficiaryAddress:              generateLoadRemittanceFromContactAccount(fTx.COUT.CashOutAcc),
		BeneficiaryCity:                 fTx.COUT.CashOutAcc.City,
		RemitterBeneficiaryRelationship: fTx.COUT.CashOutAcc.RelationshipToContact,
	}

	if fTx.COUT.CashOutAcc.Type != foree_constant.ContactAccountTypeCash {
		lrr.BeneficiaryBank = fTx.COUT.CashOutAcc.InstitutionName
		lrr.BeneficiaryAccount = fTx.COUT.CashOutAcc.AccountNumber
	}

	return lrr, nil

}

func nbpToInternalStatusMapper(nbpStatus string) transaction.TxStatus {
	switch nbpStatus {
	case "REJECTED":
		fallthrough
	case "ERROR":
		return transaction.TxStatusRejected
	case "PAID":
		return transaction.TxStatusCompleted
	case "CANCELLED":
		return transaction.TxStatusCancelled
	default:
		return transaction.TxStatusSent
	}
}

func mapNBPRemitterIdType(idType foree_auth.IdentificationType) nbp.RemitterIdType {
	switch idType {
	case foree_constant.IDTypePassport:
		return nbp.RemitterIdTypePassport
	case foree_constant.IDTypeDriverLicense:
		return nbp.RemitterIdTypeDrivinglicense
	default:
		return nbp.RemitterIdTypeOther
	}
}

func generateLoadRemittanceFromInteracAccount(user *auth.User) string {
	if user.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", user.Address1, user.City, user.Province, user.PostalCode, user.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", user.Address1, user.Address2, user.City, user.Province, user.PostalCode, user.Country)
}

func generateLoadRemittanceFromContactAccount(acc *account.ContactAccount) string {
	if acc.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", acc.Address1, acc.City, acc.Province, acc.PostalCode, acc.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", acc.Address1, acc.Address2, acc.City, acc.Province, acc.PostalCode, acc.Country)
}
