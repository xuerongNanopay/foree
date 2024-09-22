package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/partner/scotia"
)

const DefaultScotiaProfileName = "FOREE"

type ScotiaProfile struct {
	Name                          string
	LegalName                     string
	CompanyName                   string
	InitiatingPartyName           string
	ExpireInHours                 int
	SupressResponderNotifications bool
	InteracId                     string
	CountryOfResidence            string
	ProfileName                   string
	Email                         string
	AccountNumber                 string
	AccountCurrency               string
	AmountModificationAllowed     bool
	EarlyPaymentAllowed           bool
	GuaranteedPaymentRequested    bool
}

func NewCITxProcessor(
	db *sql.DB,
	scotiaProfile ScotiaProfile,
	scotiaClient scotia.ScotiaClient,
	interacTxRepo *transaction.InteracCITxRepo,
	foreeTxRepo *transaction.ForeeTxRepo,
	txSummaryRepo *transaction.TxSummaryRepo,
	txProcessor *TxProcessor,
) *CITxProcessor {
	return &CITxProcessor{
		db:            db,
		scotiaProfile: scotiaProfile,
		scotiaClient:  scotiaClient,
		interacTxRepo: interacTxRepo,
		foreeTxRepo:   foreeTxRepo,
		txSummaryRepo: txSummaryRepo,
		txProcessor:   txProcessor,
		waitFTxs:      make(map[int64]*transaction.ForeeTx, 256),
		webhookChan:   make(chan int64, 32),               // capacity with 32 should be enough.
		clearChan:     make(chan int64, 32),               // capacity with 32 should be enough.
		forwardChan:   make(chan transaction.ForeeTx, 32), // capacity with 32 should be enough.
		waitChan:      make(chan transaction.ForeeTx, 32), // capacity with 32 should be enough.
		ticker:        time.NewTicker(5 * time.Minute),
	}
}

type CITxProcessor struct {
	db            *sql.DB
	scotiaProfile ScotiaProfile
	scotiaClient  scotia.ScotiaClient
	interacTxRepo *transaction.InteracCITxRepo
	foreeTxRepo   *transaction.ForeeTxRepo
	txSummaryRepo *transaction.TxSummaryRepo
	txProcessor   *TxProcessor
	waitFTxs      map[int64]*transaction.ForeeTx
	webhookChan   chan int64
	clearChan     chan int64
	forwardChan   chan transaction.ForeeTx
	waitChan      chan transaction.ForeeTx
	ticker        *time.Ticker
}

// Loading from DB at beginning. OR, let foree processor do it.

func (p *CITxProcessor) Start() error {
	go p.startProcessor()
	return nil
}

func (p *CITxProcessor) startProcessor() {
	for {
		select {
		case fTx := <-p.waitChan:
			_, ok := p.waitFTxs[fTx.ID]
			if ok {
				//Log duplicate
			} else {
				p.waitFTxs[fTx.ID] = &fTx
			}
		case fTx := <-p.forwardChan:
			_, ok := p.waitFTxs[fTx.ID]
			if !ok {
				//Log miss
			} else {
				delete(p.waitFTxs, fTx.ID)
			}
			go func() {
				_, err := p.txProcessor.processTx(fTx)
				if err != nil {
					//log err
				}
			}()
		case foreeTxId := <-p.clearChan:
			delete(p.waitFTxs, foreeTxId)
		case foreeTxId := <-p.webhookChan:
			v, ok := p.waitFTxs[foreeTxId]
			if !ok {
				//Log error: transaction no found
			} else {
				func() {
					nTx, err := p.refreshScotiaStatus(*v)
					if err != nil {
						//Log error
						return
					}

					if nTx.CurStageStatus != v.CurStageStatus {
						p.forwardChan <- *nTx
					}
				}()
			}
		case <-p.ticker.C:
			for _, tx := range p.waitFTxs {
				func() {
					nTx, err := p.refreshScotiaStatus(*tx)
					if err != nil {
						//Log error
						return
					}
					//TODO: end process of over one day.
					if nTx.CurStageStatus != tx.CurStageStatus {
						p.forwardChan <- *nTx
					}
				}()
			}
		}
	}
}

func (p *CITxProcessor) waitFTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if fTx.CurStage != transaction.TxStageInteracCI && fTx.CurStageStatus != transaction.TxStatusSent {
		return nil, fmt.Errorf("CITxProcessor -- waitFTx -- ForeeTx `%v` is in stage `%s` at status `%s`", fTx.ID, fTx.CurStage, fTx.CurStageStatus)
	}
	p.waitChan <- fTx
	return &fTx, nil
}

func (p *CITxProcessor) process(parentTxId int64) {
	ctx := context.TODO()
	ciTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, parentTxId)
	if err != nil {
		foree_logger.Logger.Error("CI_Processor-process", "parentTxId", parentTxId, "cause", err.Error())
		return
	}
	if ciTx == nil {
		foree_logger.Logger.Error("CI_Processor-process", "parentTxId", parentTxId, "cause", "interacTx no found")
		return
	}
	switch ciTx.Status {
	case transaction.TxStatusInitial:
		//call scotia
	case transaction.TxStatusSent:
		//wait loop
	case transaction.TxStatusCompleted:
		p.txProcessor.next(ciTx.ParentTxId)
	case transaction.TxStatusRejected:
		p.txProcessor.rollback(ciTx.ParentTxId)
	case transaction.TxStatusCancelled:
		p.txProcessor.rollback(ciTx.ParentTxId)
	default:
		foree_logger.Logger.Error(
			"CI_Processor-process",
			"parentTxId", parentTxId,
			"interacCITxId", ciTx.ID,
			"interacCITxStatus", ciTx.Status,
			"cause", "unsupport status",
		)
	}
}

func (p *CITxProcessor) requestPayment(ciTx transaction.InteracCITx) {
	resp, err := p.scotiaClient.RequestPayment(*p.createRequestPaymentReq(&ciTx))

	if err != nil {
		foree_logger.Logger.Error("CI_Processor-requestPayment", "interacTxId", ciTx.ID, "cause", err.Error())
	}
	if resp.StatusCode/100 != 2 {
		foree_logger.Logger.Warn("CI_Processor-requestPayment",
			"interacTxId", ciTx.ID,
			"httpResponseStatus", resp.StatusCode,
			"httpRequest", resp.RawRequest,
			"httpResponseBody", resp.RawResponse,
			"cause", "scotia response error",
		)
	}

	if err != nil && resp.StatusCode/100 != 2 {
		ciTx.Status = transaction.TxStatusRejected
		err := p.interacTxRepo.UpdateInteracCITxById(context.TODO(), ciTx)
		if err != nil {
			foree_logger.Logger.Error("CI_Processor-requestPayment", "interacTxId", ciTx.ID, "cause", err.Error())
		}
		p.txProcessor.rollback(ciTx.ParentTxId)
		return
	}

	ciTx.ScotiaPaymentId = resp.Data.PaymentId
	statusResp, err := p.scotiaClient.PaymentStatus(scotia.PaymentStatusRequest{
		PaymentId:  ciTx.ScotiaPaymentId,
		EndToEndId: ciTx.EndToEndId,
	})

	if err != nil {
		foree_logger.Logger.Error("CI_Processor-requestPayment", "interacTxId", ciTx.ID, "cause", err.Error())
	}
	if statusResp.StatusCode/100 != 2 {
		foree_logger.Logger.Warn("CI_Processor-requestPayment",
			"interacTxId", ciTx.ID,
			"httpResponseStatus", statusResp.StatusCode,
			"httpRequest", statusResp.RawRequest,
			"httpResponseBody", statusResp.RawResponse,
			"cause", "scotia status response error",
		)
	}

	if err != nil && statusResp.StatusCode/100 != 2 {
		ciTx.Status = transaction.TxStatusRejected
		err := p.interacTxRepo.UpdateInteracCITxById(context.TODO(), ciTx)
		if err != nil {
			foree_logger.Logger.Error("CI_Processor-requestPayment", "interacTxId", ciTx.ID, "cause", err.Error())
		}
		p.txProcessor.rollback(ciTx.ParentTxId)
		return
	}

	//Success
}

// Scotia APi Call
func (p *CITxProcessor) processTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	// Lock transaction and safety check.
	nForeeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, fTx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	if nForeeTx.CurStage != transaction.TxStageInteracCI && nForeeTx.CurStageStatus != transaction.TxStatusInitial {
		dTx.Rollback()
		return nil, fmt.Errorf("CITxProcessor -- transaction `%v` is in status `%s` at stage `%s`", nForeeTx.ID, nForeeTx.CurStageStatus, nForeeTx.CurStage)
	}

	resp, err := p.scotiaClient.RequestPayment(*p.createRequestPaymentReq(fTx.CI))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		//Directly reject to avoid mess up.
		fTx.CI.Status = transaction.TxStatusRejected
		fTx.CurStageStatus = transaction.TxStatusRejected

		err = p.interacTxRepo.UpdateInteracCITxById(ctx, *fTx.CI)
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

		//TODO: log
		return &fTx, nil
	}

	//TODO: log success

	fTx.CI.ScotiaPaymentId = resp.Data.PaymentId

	// Get url payment link
	statusResp, err := p.scotiaClient.PaymentStatus(scotia.PaymentStatusRequest{
		PaymentId:  fTx.CI.ScotiaPaymentId,
		EndToEndId: fTx.CI.EndToEndId,
	})
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if statusResp.StatusCode/100 != 2 {
		//TODO: logging?
		dTx.Rollback()
		return nil, fmt.Errorf("CITxProcessor -- scotial paymentstatus error: (httpCode: `%v`, request: `%s`, response: `%s`)", statusResp.StatusCode, statusResp.RawRequest, statusResp.RawResponse)
	}

	if len(statusResp.PaymentStatuses) != 1 {
		dTx.Rollback()
		return nil, fmt.Errorf("CITxProcessor -- scotial paymentstatus error: (httpCode: `%v`, request: `%s`, response: `%s`)", statusResp.StatusCode, statusResp.RawRequest, statusResp.RawResponse)
	}

	// Update CI
	fTx.CI.PaymentUrl = statusResp.PaymentStatuses[0].GatewayUrl
	fTx.CI.ScotiaPaymentId = resp.Data.PaymentId
	fTx.CI.Status = transaction.TxStatusSent
	fTx.CI.ScotiaClearingReference = resp.Data.ClearingSystemReference

	// Update Foree
	fTx.CurStageStatus = transaction.TxStatusSent
	fTx.CI.ScotiaPaymentId = resp.Data.PaymentId

	// Update summary
	fTx.Summary.Status = transaction.TxSummaryStatusAwaitPayment
	fTx.Summary.PaymentUrl = statusResp.PaymentStatuses[0].GatewayUrl

	err = p.interacTxRepo.UpdateInteracCITxById(ctx, *fTx.CI)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	err = p.foreeTxRepo.UpdateForeeTxById(ctx, fTx)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	err = p.txSummaryRepo.UpdateTxSummaryById(ctx, *fTx.Summary)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if err = dTx.Commit(); err != nil {
		return nil, err
	}

	p.waitChan <- fTx

	return &fTx, nil
}

// The function normally in a goroutine
func (p *CITxProcessor) refreshScotiaStatus(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {

	detailResp, err := p.scotiaClient.PaymentDetail(scotia.PaymentDetailRequest{
		PaymentId:  fTx.CI.ScotiaPaymentId,
		EndToEndId: fTx.CI.EndToEndId,
	})
	if err != nil {
		return nil, err
	}

	if detailResp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("CITxProcessor -- refreshScotiaStatusAndProcess -- scotia paymentdetail error: (httpCode: `%v`, request: `%s`, response: `%s`)", detailResp.StatusCode, detailResp.RawRequest, detailResp.RawResponse)
	}

	scotiaStatus := detailResp.PaymentDetail.TransactionStatus
	if scotiaStatus == "" {
		scotiaStatus = detailResp.PaymentDetail.RequestForPaymentStatus
	}

	newStatus := scotiaToInternalStatusMapper(scotiaStatus)
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

	if curFTx.CurStage != transaction.TxStageInteracCI && curFTx.CurStageStatus != transaction.TxStatusSent {
		dTx.Rollback()
		p.clearChan <- fTx.ID
		return nil, fmt.Errorf("CITxProcessor -- refreshScotiaStatusAndProcess -- ForeeTx `%v` is in stage `%s` at status `%s`", curFTx.ID, curFTx.CurStage, curFTx.CurStageStatus)
	}

	// Update Foree Tx and CI Tx.
	fTx.CI.Status = newStatus
	fTx.CI.ScotiaStatus = scotiaStatus
	fTx.CurStageStatus = newStatus

	err = p.interacTxRepo.UpdateInteracCITxById(ctx, *fTx.CI)
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

func (p *CITxProcessor) Webhook(paymentId string) {
	ciTx, err := p.interacTxRepo.GetUniqueInteracCITxByScotiaPaymentId(context.TODO(), paymentId)
	if err != nil {
		//TODO: Log error
	}
	if ciTx == nil {
		//TODO: Log error
	}
	p.webhookChan <- ciTx.ParentTxId
}

func scotiaToInternalStatusMapper(scotiaStatus string) transaction.TxStatus {
	switch scotiaStatus {
	case "ACCC":
		fallthrough
	case "ACSP":
		fallthrough
	case "COMPLETED":
		fallthrough
	case "REALTIME_DEPOSIT_COMPLETED":
		fallthrough
	case "DEPOSIT_COMPLETE":
		return transaction.TxStatusCompleted
	case "RJCT":
		fallthrough
	case "DECLINED":
		fallthrough
	case "REALTIME_DEPOSIT_FAILED":
		fallthrough
	case "DIRECT_DEPOSIT_FAILED":
		return transaction.TxStatusRejected
	case "CANCELLED":
		fallthrough
	case "EXPIRED":
		return transaction.TxStatusCancelled
	default:
		return transaction.TxStatusSent
	}
}

func (p *CITxProcessor) createRequestPaymentReq(ciTx *transaction.InteracCITx) *scotia.RequestPaymentRequest {
	expireDate := time.Now().Add(time.Hour * time.Duration(p.scotiaProfile.ExpireInHours))

	req := &scotia.RequestPaymentRequest{
		RequestData: &scotia.RequestPaymentRequestData{
			ProductCode:                    "DOMESTIC",
			MessageIdentification:          ciTx.EndToEndId,
			EndToEndIdentification:         ciTx.EndToEndId,
			CreditDebitIndicator:           "CRDT",
			CreationDatetime:               (*scotia.ScotiaDatetime)(ciTx.CreatedAt),
			PaymentExpiryDate:              (*scotia.ScotiaDatetime)(&expireDate),
			SuppressResponderNotifications: p.scotiaProfile.SupressResponderNotifications,
			ReturnUrl:                      "string",
			Language:                       "EN",
			InstructedAmtData: &scotia.ScotiaAmtData{
				Amount:   scotia.ScotiaAmount(ciTx.Amt.Amount),
				Currency: ciTx.Amt.Currency,
			},
			InitiatingParty: &scotia.InitiatingPartyData{
				Name: p.scotiaProfile.InitiatingPartyName,
				Identification: &scotia.IdentificationData{
					OrganisationIdentification: &scotia.OrganisationIdentificationData{
						Other: []scotia.OtherData{
							{
								Identification: p.scotiaProfile.InteracId,
								SchemeName: &scotia.SchemeNameData{
									Code: "BANK",
								},
							},
						},
					},
				},
			},
			Debtor: &scotia.DebtorData{
				Name:               ciTx.CashInAcc.GetLegalName(),
				CountryOfResidence: p.scotiaProfile.CountryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: ciTx.CashInAcc.Email,
				},
			},
			Creditor: &scotia.CreditorData{
				Name:               p.scotiaProfile.LegalName,
				CountryOfResidence: p.scotiaProfile.CountryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: p.scotiaProfile.Email,
				},
			},
			UltimateCreditor: &scotia.CreditorData{
				Name:               p.scotiaProfile.LegalName,
				CountryOfResidence: p.scotiaProfile.CountryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: p.scotiaProfile.Email,
				},
			},
			CreditorAccount: &scotia.CreditorAccountData{
				Identification: p.scotiaProfile.AccountNumber,
				Currency:       p.scotiaProfile.AccountCurrency,
				SchemeName:     "ALIAS_ACCT_NO",
			},
			FraudSupplementaryInfo: &scotia.FraudSupplementaryInfoData{
				CustomerAuthenticationMethod: "PASSWORD",
			},
			PaymentCondition: &scotia.PaymentConditionData{
				AmountModificationAllowed:  p.scotiaProfile.AmountModificationAllowed,
				EarlyPaymentAllowed:        p.scotiaProfile.EarlyPaymentAllowed,
				GuaranteedPaymentRequested: p.scotiaProfile.GuaranteedPaymentRequested,
			},
			PaymentTypeInformation: &scotia.PaymentTypeInformationData{
				CategoryPurpose: &scotia.CategoryPurposeData{
					Code: "CASH",
				},
			},
		},
	}
	return req
}
