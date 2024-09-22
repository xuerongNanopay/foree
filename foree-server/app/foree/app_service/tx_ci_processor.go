package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/partner/scotia"
)

const DefaultScotiaProfileName = "FOREE"
const waitRecheckInterval = 5 * time.Minute

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
	ret := &CITxProcessor{
		db:                  db,
		scotiaProfile:       scotiaProfile,
		scotiaClient:        scotiaClient,
		interacTxRepo:       interacTxRepo,
		foreeTxRepo:         foreeTxRepo,
		txSummaryRepo:       txSummaryRepo,
		txProcessor:         txProcessor,
		statusPullingChan:   make(chan transaction.InteracCITx, 64),
		refreshStatusChan:   make(chan string, 64),
		statusRecheckticker: time.NewTicker(5 * time.Minute),
		statusRefreshticker: time.NewTicker(5 * time.Minute),
	}
	ret.start()
	return ret
}

type waitWrapper struct {
	ciTx      transaction.InteracCITx
	recheckAt time.Time
}

type CITxProcessor struct {
	db                  *sql.DB
	scotiaProfile       ScotiaProfile
	scotiaClient        scotia.ScotiaClient
	interacTxRepo       *transaction.InteracCITxRepo
	foreeTxRepo         *transaction.ForeeTxRepo
	txSummaryRepo       *transaction.TxSummaryRepo
	txProcessor         *TxProcessor
	waits               sync.Map
	statusPullingChan   chan transaction.InteracCITx
	refreshStatusChan   chan string
	statusRecheckticker *time.Ticker
	statusRefreshticker *time.Ticker
}

func (p *CITxProcessor) start() {
	for {
		select {
		// Push into map.
		case ciTx := <-p.statusPullingChan:
			_, ok := p.waits.Load(ciTx.ScotiaPaymentId)
			if ok {
				continue
			}
			p.waits.Store(ciTx.ScotiaPaymentId, waitWrapper{
				ciTx:      ciTx,
				recheckAt: time.Now().Add(waitRecheckInterval),
			})
		case paymentId := <-p.refreshStatusChan:
			v, ok := p.waits.Load(paymentId)
			w, _ := v.(waitWrapper)
			if !ok {
				ciTx, err := p.interacTxRepo.GetUniqueInteracCITxByScotiaPaymentId(context.TODO(), paymentId)
				if err != nil {
					foree_logger.Logger.Error("CITxprocessor-scotiaWebhoodChan_getInteracCiTxByScotiaPaymentId_FAIL",
						"socitaPaymentId", paymentId,
						"cause", err.Error(),
					)
					continue
				}
				if ciTx == nil {
					foree_logger.Logger.Error("CITxprocessor-scotiaWebhoodChan_getInteracCiTxByScotiaPaymentId_FAIL",
						"socitaPaymentId", paymentId,
						"cause", "InteracTx no found",
					)
					continue
				}
				if ciTx.Status != transaction.TxStatusSent {
					foree_logger.Logger.Warn("CITxprocessor-scotiaWebhoodChan_FAIL",
						"socitaPaymentId", paymentId,
						"interacCITxId", ciTx.ID,
						"interacCITxStatus", ciTx.Status,
						"cause", "get scotia webhook but interacCITx.status is incorrect",
					)
					continue
				}
				p.waits.Store(paymentId, waitWrapper{
					ciTx:      *ciTx,
					recheckAt: time.Now().Add(waitRecheckInterval),
				})
				v, _ := p.waits.Load(paymentId)
				w, _ = v.(waitWrapper)
				foree_logger.Logger.Info("CITxprocessor-scotiaWebhoodChan",
					"socitaPaymentId", paymentId,
					"interacCITxId", ciTx.ID,
					"interacCITxStatus", ciTx.Status,
					"msg", "interacCITx miss in wait map, add it back",
				)
			}

			newStatus, newScotiaStatus, err := p.refreshScotiaStatus(w.ciTx)
			if err != nil {
				foree_logger.Logger.Error("CITxprocessor-scotiaWebhoodChan_FAIL",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.ciTx.ID,
					"interacCITxStatus", w.ciTx.Status,
					"cause", err.Error(),
				)
				continue
			}
			if newStatus == w.ciTx.Status && newScotiaStatus == w.ciTx.ScotiaStatus {
				foree_logger.Logger.Debug("CITxprocessor-scotiaWebhoodChan_still_in_waiting",
					"socitaPaymentId", paymentId,
					"foreeTxId", w.ciTx.ParentTxId,
					"interacCITxId", w.ciTx.ID,
					"newInteracCITxStatus", newStatus,
					"newScotiaStatus", newScotiaStatus,
				)
				continue
			}
			curCiTx, err := p.interacTxRepo.GetUniqueInteracCITxByScotiaPaymentId(context.TODO(), paymentId)
			if err != nil {
				foree_logger.Logger.Error("CITxprocessor-scotiaWebhoodChan_FAIL",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.ciTx.ID,
					"interacCITxStatus", w.ciTx.Status,
					"cause", err.Error(),
				)
				continue
			}
			if curCiTx.Status != transaction.TxStatusSent {
				foree_logger.Logger.Warn("CITxprocessor-scotiaWebhoodChan_FAIL",
					"socitaPaymentId", paymentId,
					"foreeTxId", curCiTx.ParentTxId,
					"interacCITxId", curCiTx.ID,
					"curInteracCITxStatus", curCiTx.Status,
					"cause", "interacCITx is not in SENT status",
				)
				p.waits.Delete(paymentId)
				continue
			}

			curCiTx.Status = newStatus
			curCiTx.ScotiaStatus = newScotiaStatus

			err = p.interacTxRepo.UpdateInteracCITxById(context.TODO(), *curCiTx)
			if err != nil {
				foree_logger.Logger.Error("CITxprocessor-scotiaWebhoodChan_FAIL",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.ciTx.ID,
					"newInteracCITxStatus", newStatus,
					"newScotiaStatus", newScotiaStatus,
					"cause", err.Error(),
				)
				continue
			}
			// Still in send
			if curCiTx.Status == w.ciTx.Status {
				//TODO: Update status
				go p.txProcessor.onStatusUpdate(curCiTx.ParentTxId)
				p.waits.Swap(paymentId, waitWrapper{
					ciTx:      *curCiTx,
					recheckAt: time.Now().Add(waitRecheckInterval),
				})
			} else {
				p.waits.Delete(paymentId)
				go p.process(curCiTx.ParentTxId)
				foree_logger.Logger.Info("CITxprocessor-scotiaWebhoodChan_SUCCESS",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.ciTx.ID,
					"newInteracCITxStatus", newStatus,
					"newScotiaStatus", newScotiaStatus,
				)
			}
		case <-p.statusRefreshticker.C:
			func() {
				p.waits.Range(func(k, _ interface{}) bool {
					paymentId, _ := k.(string)
					p.refreshStatusChan <- paymentId
					return true
				})
			}()
		case <-p.statusRecheckticker.C:
			func() {

			}()
		}
	}
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
		p.requestPayment(*ciTx)
	case transaction.TxStatusSent:
		p.statusPullingChan <- *ciTx
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
		foree_logger.Logger.Error("CITxProcessor-requestPayment_FAIL", "interacTxId", ciTx.ID, "cause", err.Error())
	}
	if resp.StatusCode/100 != 2 {
		foree_logger.Logger.Warn("CITxProcessor-requestPayment_FAIL",
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
			foree_logger.Logger.Error("CITxProcessor-requestPayment_FAIL", "interacTxId", ciTx.ID, "cause", err.Error())
		}
		p.process(ciTx.ParentTxId)
		return
	}

	ciTx.ScotiaPaymentId = resp.Data.PaymentId
	statusResp, err := p.scotiaClient.PaymentStatus(scotia.PaymentStatusRequest{
		PaymentId:  ciTx.ScotiaPaymentId,
		EndToEndId: ciTx.EndToEndId,
	})

	if err != nil {
		foree_logger.Logger.Error("CITxProcessor-requestPayment_FAIL", "interacTxId", ciTx.ID, "cause", err.Error())
	}
	if statusResp.StatusCode/100 != 2 {
		foree_logger.Logger.Warn("CITxProcessor-requestPayment_FAIL",
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
			foree_logger.Logger.Error("CITxProcessor-requestPayment_FAIL", "interacTxId", ciTx.ID, "cause", err.Error())
		}
		p.process(ciTx.ParentTxId)
		return
	}

	//Success
	ciTx.PaymentUrl = statusResp.PaymentStatuses[0].GatewayUrl
	ciTx.ScotiaPaymentId = resp.Data.PaymentId
	ciTx.Status = transaction.TxStatusSent
	ciTx.ScotiaClearingReference = resp.Data.ClearingSystemReference

	err = p.interacTxRepo.UpdateInteracCITxById(context.TODO(), ciTx)
	if err != nil {
		foree_logger.Logger.Error("CITxProcessor-requestPayment_FAIL", "interacTxId", ciTx.ID, "cause", err.Error())
	}
	foree_logger.Logger.Info("CITxProcessor-requestPayment_SUCCESS",
		"interacTxId", ciTx.ID,
	)

	p.process(ciTx.ID)
}

// Scotia APi Call
func (p *CITxProcessor) processTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

// The function normally in a goroutine
func (p *CITxProcessor) refreshScotiaStatus(ciTx transaction.InteracCITx) (transaction.TxStatus, string, error) {

	detailResp, err := p.scotiaClient.PaymentDetail(scotia.PaymentDetailRequest{
		PaymentId:  ciTx.ScotiaPaymentId,
		EndToEndId: ciTx.EndToEndId,
	})
	if err != nil {
		return "", "", err
	}

	if detailResp.StatusCode/100 != 2 {
		return "", "", fmt.Errorf("CITxProcessor -- refreshScotiaStatusAndProcess -- scotia paymentdetail error: (httpCode: `%v`, request: `%s`, response: `%s`)", detailResp.StatusCode, detailResp.RawRequest, detailResp.RawResponse)
	}

	scotiaStatus := detailResp.PaymentDetail.TransactionStatus
	if scotiaStatus == "" {
		scotiaStatus = detailResp.PaymentDetail.RequestForPaymentStatus
	}

	newStatus := scotiaToInternalStatusMapper(scotiaStatus)
	return newStatus, scotiaStatus, nil
}

func (p *CITxProcessor) Webhook(paymentId string) {
	p.refreshStatusChan <- paymentId
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
