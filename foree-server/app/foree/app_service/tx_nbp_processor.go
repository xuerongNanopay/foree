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
	"xue.io/go-pay/constant"
	"xue.io/go-pay/partner/nbp"
)

const nbpTxRecheckInterval = 10 * time.Minute
const nbpTxRetryInterval = 10 * time.Second
const nbpTxRetryAttempts = 5
const nbpTxProcessorRefreshTicker = 5 * time.Minute

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
		retryFTxs:              make(map[int64]*transaction.ForeeTx, 64),
		waitFTxs:               make(map[int64]*transaction.ForeeTx, 256),
		clearChan:              make(chan int64, 32),               // capacity with 32 should be enough.
		forwardChan:            make(chan transaction.ForeeTx, 32), // capacity with 32 should be enough.
		checkStatusTicker:      time.NewTicker(nbpTxProcessorRefreshTicker),
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
	retryFTxs              map[int64]*transaction.ForeeTx
	waitFTxs               map[int64]*transaction.ForeeTx
	waitChan               chan transaction.ForeeTx
	forwardChan            chan transaction.ForeeTx
	checkStatusTicker      *time.Ticker
	clearChan              chan int64
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
			//TODO: check status.
			func() {}()
		}
	}
}

func (p *NBPTxProcessor) waitFTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if fTx.CurStage != transaction.TxStageInteracCI && fTx.CurStageStatus != transaction.TxStatusSent {
		return nil, fmt.Errorf("NBPTxProcessor -- waitFTx -- ForeeTx `%v` is in stage `%s` at status `%s`", fTx.ID, fTx.CurStage, fTx.CurStageStatus)
	}
	p.waitChan <- fTx
	return &fTx, nil
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
	case transaction.TxStatusCompleted:
		p.txProcessor.next(nbpTx.ParentTxId)
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
		foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
		return
	}
	req, err := p.buildLoadRemittanceRequest(*fTx)
	if err != nil {
		foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
		return
	}

	mode, err := mapNBPMode(fTx.COUT.CashOutAcc.Type)
	if err != nil {
		foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
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
			foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL", "parentTxId", parentTxId, "cause", err.Error())
			return
		}
		//Retry case: 5xx, 401, 403
		if resp.StatusCode/100 == 5 || resp.ResponseCode == "401" || resp.ResponseCode == "403" || resp.ResponseCode == "406" {
			foree_logger.Logger.Warn("IDM_Processor-loadRemittance", "parentTxId", parentTxId, "retry", i)
			time.Sleep(nbpTxRetryInterval)
		} else {
			break
		}

	}

	// Retry later manully
	if resp.StatusCode/100 == 5 || resp.ResponseCode == "401" || resp.ResponseCode == "403" {
		foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL",
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
			foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL",
				"parentTxId", parentTxId,
				"httpStatus", resp.StatusCode,
				"cause", err.Error(),
			)
			return
		}
		foree_logger.Logger.Info("IDM_Processor-loadRemittance_SUCCESS", "parentTxId", parentTxId)
		p.process(nbpTx.ParentTxId)
		return
	}

	// Reject
	foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL",
		"parentTxId", parentTxId,
		"httpStatus", resp.StatusCode,
		"httpRequest", resp.RawRequest,
		"httpResponse", resp.RawResponse,
		"msg", "please re-run the transaction.",
	)
	nbpTx.Status = transaction.TxStatusRejected
	err = p.nbpTxRepo.UpdateNBPCOTxById(context.TODO(), nbpTx)
	if err != nil {
		foree_logger.Logger.Error("IDM_Processor-loadRemittance_FAIL",
			"parentTxId", parentTxId,
			"httpStatus", resp.StatusCode,
			"cause", err.Error(),
		)
		return
	}
	p.process(nbpTx.ParentTxId)
}

// We don't use transaction here, case NBP can check duplicate.
func (p *NBPTxProcessor) processTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *NBPTxProcessor) refreshNBPStatus(fTx transaction.ForeeTx, nbpStatus string) (*transaction.ForeeTx, error) {
	newStatus := nbpToInternalStatusMapper(nbpStatus)
	if newStatus == transaction.TxStatusSent {
		return &fTx, nil
	}

	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)
	curFTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.CI.ParentTxId)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if curFTx.CurStage != transaction.TxStageNBPCO && curFTx.CurStageStatus != transaction.TxStatusSent {
		dTx.Rollback()
		p.clearChan <- fTx.ID
		return nil, fmt.Errorf("NBPTxProcessor -- refreshNBPStatus -- ForeeTx `%v` is in stage `%s` at status `%s`", curFTx.ID, curFTx.CurStage, curFTx.CurStageStatus)
	}

	fTx.COUT.Status = newStatus
	fTx.CurStageStatus = newStatus

	err = p.nbpTxRepo.UpdateNBPCOTxById(ctx, *fTx.COUT)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	err = p.foreeTxRepo.UpdateForeeTxById(ctx, fTx)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if err = dTx.Commit(); err != nil {
		return nil, err
	}

	return &fTx, nil
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

func (p *NBPTxProcessor) buildLoadRemittanceRequest(fTx transaction.ForeeTx) (*nbp.LoadRemittanceRequest, error) {
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
		RemitterContact:                 fTx.CI.CashInAcc.PhoneNumber,
		RemitterDOB:                     (*nbp.NBPDate)(fTx.Owner.Dob),
		RemitterAddress:                 generateLoadRemittanceFromInteracAccount(fTx.CI.CashInAcc),
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

// TODO: fix the issue
func mapNBPMode(accType account.ContactAccountType) (nbp.PMTMode, error) {
	// switch accType {
	// case foree_constant.ContactAccountTypeCash:
	// 	return nbp.PMTModeCash, nil
	// case foree_constant.ContactAccountTypeThirdPartyPayments:
	// 	return nbp.PMTModeThirdPartyPayments, nil
	// case foree_constant.ContactAccountTypeAccountTransfers:
	// 	return nbp.PMTModeAccountTransfers, nil
	// default:
	// 	return "", fmt.Errorf("NBPTxProcessor -- unknown contact account type `%s`", accType)
	// }
	return nbp.PMTModeCash, nil
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

func generateLoadRemittanceFromInteracAccount(acc *account.InteracAccount) string {
	if acc.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", acc.Address1, acc.City, acc.Province, acc.PostalCode, acc.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", acc.Address1, acc.Address2, acc.City, acc.Province, acc.PostalCode, acc.Country)
}

func generateLoadRemittanceFromContactAccount(acc *account.ContactAccount) string {
	if acc.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", acc.Address1, acc.City, acc.Province, acc.PostalCode, acc.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", acc.Address1, acc.Address2, acc.City, acc.Province, acc.PostalCode, acc.Country)
}

func chunkNBPIds(s []string, splitSize int) [][]string {
	numberOfSlices := len(s) / splitSize
	remainder := len(s) % splitSize

	ret := make([][]string, 0)
	start := 0
	end := 0

	for i := 0; i < numberOfSlices; i++ {
		end += splitSize
		ret = append(ret, s[start:end])
		start = end
	}

	if remainder > 0 {
		end = start + remainder
		ret = append(ret, s[start:end])
	}
	return ret
}
