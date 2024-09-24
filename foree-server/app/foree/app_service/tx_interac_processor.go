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
const interacTxRecheckInterval = 10 * time.Minute

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

func NewInteracTxProcessor(
	db *sql.DB,
	scotiaProfile ScotiaProfile,
	scotiaClient scotia.ScotiaClient,
	interacTxRepo *transaction.InteracCITxRepo,
	foreeTxRepo *transaction.ForeeTxRepo,
	txSummaryRepo *transaction.TxSummaryRepo,
	txProcessor *TxProcessor,
) *InteracTxProcessor {
	ret := &InteracTxProcessor{
		db:                  db,
		scotiaProfile:       scotiaProfile,
		scotiaClient:        scotiaClient,
		interacTxRepo:       interacTxRepo,
		foreeTxRepo:         foreeTxRepo,
		txSummaryRepo:       txSummaryRepo,
		txProcessor:         txProcessor,
		statusRefreshChan:   make(chan transaction.InteracCITx, 64),
		refreshStatusChan:   make(chan string, 64),
		statusRecheckticker: time.NewTicker(5 * time.Minute),
		statusRefreshTicker: time.NewTicker(5 * time.Minute),
	}
	go ret.start()
	//TODO: load all wait transaction from DB
	return ret
}

type interacTxWrapper struct {
	interacTx transaction.InteracCITx
	recheckAt time.Time
}

type InteracTxProcessor struct {
	db                  *sql.DB
	scotiaProfile       ScotiaProfile
	scotiaClient        scotia.ScotiaClient
	interacTxRepo       *transaction.InteracCITxRepo
	foreeTxRepo         *transaction.ForeeTxRepo
	txSummaryRepo       *transaction.TxSummaryRepo
	txProcessor         *TxProcessor
	waits               sync.Map
	statusRefreshChan   chan transaction.InteracCITx
	refreshStatusChan   chan string
	statusRecheckticker *time.Ticker
	statusRefreshTicker *time.Ticker
}

func (p *InteracTxProcessor) start() {
	for {
		select {
		// Push into map.
		case interacTx := <-p.statusRefreshChan:
			_, ok := p.waits.Load(interacTx.ScotiaPaymentId)
			if ok {
				foree_logger.Logger.Warn("InteracTxProcessor--statusRefreshChan",
					"interacTxId", interacTx.ID,
					"msg", "interacTx is in waiting aleardy",
				)
				continue
			}
			p.waits.Store(interacTx.ScotiaPaymentId, interacTxWrapper{
				interacTx: interacTx,
				recheckAt: time.Now().Add(interacTxRecheckInterval),
			})
		case paymentId := <-p.refreshStatusChan:
			v, ok := p.waits.Load(paymentId)
			w, _ := v.(interacTxWrapper)
			if !ok {
				interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByScotiaPaymentId(context.TODO(), paymentId)
				if err != nil {
					foree_logger.Logger.Error("InteracTxProcessor--scotiaWebhoodChan_getInteracCiTxByScotiaPaymentId_FAIL",
						"socitaPaymentId", paymentId,
						"cause", err.Error(),
					)
					continue
				}
				if interacTx == nil {
					foree_logger.Logger.Error("InteracTxProcessor--scotiaWebhoodChan_getInteracCiTxByScotiaPaymentId_FAIL",
						"socitaPaymentId", paymentId,
						"cause", "InteracTx no found",
					)
					continue
				}
				if interacTx.Status != transaction.TxStatusSent {
					foree_logger.Logger.Warn("InteracTxProcessor--scotiaWebhoodChan_FAIL",
						"socitaPaymentId", paymentId,
						"interacCITxId", interacTx.ID,
						"interacCITxStatus", interacTx.Status,
						"cause", "get scotia webhook but interacCITx.status is incorrect",
					)
					continue
				}
				p.waits.Store(paymentId, interacTxWrapper{
					interacTx: *interacTx,
					recheckAt: time.Now().Add(interacTxRecheckInterval),
				})
				v, _ := p.waits.Load(paymentId)
				w, _ = v.(interacTxWrapper)
				foree_logger.Logger.Info("InteracTxProcessor--scotiaWebhoodChan",
					"socitaPaymentId", paymentId,
					"interacCITxId", interacTx.ID,
					"interacCITxStatus", interacTx.Status,
					"msg", "interacCITx miss in wait map, add it back",
				)
			}

			newStatus, newScotiaStatus, err := p.refreshScotiaStatus(w.interacTx)
			if err != nil {
				foree_logger.Logger.Error("InteracTxProcessor--scotiaWebhoodChan_FAIL",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.interacTx.ID,
					"interacCITxStatus", w.interacTx.Status,
					"cause", err.Error(),
				)
				continue
			}
			//No status change
			if newStatus == w.interacTx.Status && newScotiaStatus == w.interacTx.ScotiaStatus {
				foree_logger.Logger.Debug("InteracTxProcessor--scotiaWebhoodChan_still_in_waiting",
					"socitaPaymentId", paymentId,
					"foreeTxId", w.interacTx.ParentTxId,
					"interacCITxId", w.interacTx.ID,
					"newInteracCITxStatus", newStatus,
					"newScotiaStatus", newScotiaStatus,
				)
				continue
			}
			curCiTx, err := p.interacTxRepo.GetUniqueInteracCITxByScotiaPaymentId(context.TODO(), paymentId)
			if err != nil {
				foree_logger.Logger.Error("InteracTxProcessor--scotiaWebhoodChan_FAIL",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.interacTx.ID,
					"interacCITxStatus", w.interacTx.Status,
					"cause", err.Error(),
				)
				continue
			}
			//status laging.
			if curCiTx.Status != transaction.TxStatusSent {
				foree_logger.Logger.Warn("InteracTxProcessor--scotiaWebhoodChan_FAIL",
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
				foree_logger.Logger.Error("InteracTxProcessor--scotiaWebhoodChan_FAIL",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.interacTx.ID,
					"newInteracCITxStatus", newStatus,
					"newScotiaStatus", newScotiaStatus,
					"cause", err.Error(),
				)
				continue
			}
			// Still in send.
			if curCiTx.Status == transaction.TxStatusSent {
				go p.txProcessor.onStatusUpdate(curCiTx.ParentTxId)
				p.waits.Swap(paymentId, interacTxWrapper{
					interacTx: *curCiTx,
					recheckAt: time.Now().Add(interacTxRecheckInterval),
				})
			} else {
				p.waits.Delete(paymentId)
				//moving forward
				go p.process(curCiTx.ParentTxId)
				foree_logger.Logger.Info("InteracTxProcessor--scotiaWebhoodChan",
					"socitaPaymentId", paymentId,
					"interacCITxId", w.interacTx.ID,
					"newInteracCITxStatus", newStatus,
					"newScotiaStatus", newScotiaStatus,
				)
			}
		case <-p.statusRefreshTicker.C:
			waitPaymentIds := make([]string, 0)
			p.waits.Range(func(k, _ interface{}) bool {
				paymentId, _ := k.(string)
				waitPaymentIds = append(waitPaymentIds, paymentId)
				return true
			})
			go func() {
				for _, paymentId := range waitPaymentIds {
					p.refreshStatusChan <- paymentId
				}
			}()
		}
	}
}

func (p *InteracTxProcessor) process(parentTxId int64) {
	ctx := context.TODO()
	interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, parentTxId)
	if err != nil {
		foree_logger.Logger.Error("CITxProcessor-process", "parentTxId", parentTxId, "cause", err.Error())
		return
	}
	if interacTx == nil {
		foree_logger.Logger.Error("CITxProcessor-process", "parentTxId", parentTxId, "cause", "interacTx no found")
		return
	}
	switch interacTx.Status {
	case transaction.TxStatusInitial:
		p.requestPayment(*interacTx)
	case transaction.TxStatusSent:
		p.statusRefreshChan <- *interacTx
	case transaction.TxStatusCompleted:
		p.txProcessor.next(interacTx.ParentTxId)
	case transaction.TxStatusRejected:
		p.txProcessor.rollback(interacTx.ParentTxId)
	case transaction.TxStatusCancelled:
		p.txProcessor.rollback(interacTx.ParentTxId)
	default:
		foree_logger.Logger.Error(
			"CITxProcessor-process",
			"parentTxId", parentTxId,
			"interacCITxId", interacTx.ID,
			"interacCITxStatus", interacTx.Status,
			"cause", "unsupport status",
		)
	}
}

func (p *InteracTxProcessor) requestPayment(interacTx transaction.InteracCITx) {
	resp, err := p.scotiaClient.RequestPayment(*p.createRequestPaymentReq(&interacTx))

	if err != nil {
		foree_logger.Logger.Error("InteracTxProcessor--requestPayment_FAIL", "interacTxId", interacTx.ID, "cause", err.Error())
	}
	if resp.StatusCode/100 != 2 {
		foree_logger.Logger.Warn("InteracTxProcessor--requestPayment_FAIL",
			"interacTxId", interacTx.ID,
			"httpResponseStatus", resp.StatusCode,
			"httpRequest", resp.RawRequest,
			"httpResponseBody", resp.RawResponse,
			"cause", "scotia response error",
		)
	}

	if err != nil && resp.StatusCode/100 != 2 {
		interacTx.Status = transaction.TxStatusRejected
		err := p.interacTxRepo.UpdateInteracCITxById(context.TODO(), interacTx)
		if err != nil {
			foree_logger.Logger.Error("InteracTxProcessor--requestPayment_FAIL", "interacTxId", interacTx.ID, "cause", err.Error())
		}
		p.process(interacTx.ParentTxId)
		return
	}

	interacTx.ScotiaPaymentId = resp.Data.PaymentId
	statusResp, err := p.scotiaClient.PaymentStatus(scotia.PaymentStatusRequest{
		PaymentId:  interacTx.ScotiaPaymentId,
		EndToEndId: interacTx.EndToEndId,
	})

	if err != nil {
		foree_logger.Logger.Error("InteracTxProcessor--requestPayment_FAIL", "interacTxId", interacTx.ID, "cause", err.Error())
	}
	if statusResp.StatusCode/100 != 2 {
		foree_logger.Logger.Warn("InteracTxProcessor--requestPayment_FAIL",
			"interacTxId", interacTx.ID,
			"httpResponseStatus", statusResp.StatusCode,
			"httpRequest", statusResp.RawRequest,
			"httpResponseBody", statusResp.RawResponse,
			"cause", "scotia status response error",
		)
	}

	if err != nil && statusResp.StatusCode/100 != 2 {
		interacTx.Status = transaction.TxStatusRejected
		err := p.interacTxRepo.UpdateInteracCITxById(context.TODO(), interacTx)
		if err != nil {
			foree_logger.Logger.Error("InteracTxProcessor--requestPayment_FAIL", "interacTxId", interacTx.ID, "cause", err.Error())
		}
		p.process(interacTx.ParentTxId)
		return
	}

	//Success
	interacTx.PaymentUrl = statusResp.PaymentStatuses[0].GatewayUrl
	interacTx.ScotiaPaymentId = resp.Data.PaymentId
	interacTx.Status = transaction.TxStatusSent
	interacTx.ScotiaClearingReference = resp.Data.ClearingSystemReference

	err = p.interacTxRepo.UpdateInteracCITxById(context.TODO(), interacTx)
	if err != nil {
		foree_logger.Logger.Error("InteracTxProcessor--requestPayment_FAIL", "interacTxId", interacTx.ID, "cause", err.Error())
		return
	}
	foree_logger.Logger.Info("InteracTxProcessor--requestPayment_SUCCESS",
		"interacTxId", interacTx.ID,
	)

	p.process(interacTx.ParentTxId)
}

func (p *InteracTxProcessor) refreshScotiaStatus(interacTx transaction.InteracCITx) (transaction.TxStatus, string, error) {

	detailResp, err := p.scotiaClient.PaymentDetail(scotia.PaymentDetailRequest{
		PaymentId:  interacTx.ScotiaPaymentId,
		EndToEndId: interacTx.EndToEndId,
	})
	if err != nil {
		return "", "", err
	}

	if detailResp.StatusCode/100 != 2 {
		return "", "", fmt.Errorf("InteracTxProcessor -- refreshScotiaStatusAndProcess -- scotia paymentdetail error: (httpCode: `%v`, request: `%s`, response: `%s`)", detailResp.StatusCode, detailResp.RawRequest, detailResp.RawResponse)
	}

	scotiaStatus := detailResp.PaymentDetail.TransactionStatus
	if scotiaStatus == "" {
		scotiaStatus = detailResp.PaymentDetail.RequestForPaymentStatus
	}

	newStatus := scotiaToInternalStatusMapper(scotiaStatus)
	return newStatus, scotiaStatus, nil
}

func (p *InteracTxProcessor) Webhook(paymentId string) {
	p.refreshStatusChan <- paymentId
}

func (p *InteracTxProcessor) ManualUpdate(parentTxId int64, newTxStatus transaction.TxStatus) error {
	if newTxStatus != transaction.TxStatusRejected && newTxStatus != transaction.TxStatusCompleted {
		return fmt.Errorf("unsupport transaction status `%v`", newTxStatus)
	}

	ctx := context.TODO()
	interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, parentTxId)
	if err != nil {
		return err
	}
	if interacTx == nil {
		return fmt.Errorf("InteracCITx no found with parentTxId `%v`", parentTxId)
	}
	if interacTx.Status != transaction.TxStatusSent {
		return fmt.Errorf("expect InteracCITx in `%v`, but got `%v`", transaction.TxStatusSent, interacTx.Status)
	}

	interacTx.Status = transaction.TxStatusCompleted
	err = p.interacTxRepo.UpdateInteracCITxById(context.TODO(), *interacTx)
	if err != nil {
		return err
	}

	p.waits.Delete(interacTx.ScotiaPaymentId)
	go p.txProcessor.next(interacTx.ParentTxId)
	return nil
}

// TODO: call scotial cancel api
func (p *InteracTxProcessor) Cancel(parentTxId int64) error {
	ctx := context.TODO()
	interacTx, err := p.interacTxRepo.GetUniqueInteracCITxByParentTxId(ctx, parentTxId)
	if err != nil {
		return err
	}
	if interacTx == nil {
		return fmt.Errorf("InteracCITx no found with parentTxId `%v`", parentTxId)
	}
	if interacTx.Status != transaction.TxStatusSent {
		return fmt.Errorf("expect InteracCITx in `%v`, but got `%v`", transaction.TxStatusSent, interacTx.Status)
	}

	//TODO: call scotial cancel api
	//TODO: if error return.
	interacTx.Status = transaction.TxStatusCancelled
	err = p.interacTxRepo.UpdateInteracCITxById(context.TODO(), *interacTx)
	if err != nil {
		return err
	}
	p.waits.Delete(interacTx.ScotiaPaymentId)
	go p.txProcessor.rollback(interacTx.ParentTxId)
	return nil
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

func (p *InteracTxProcessor) createRequestPaymentReq(interacTx *transaction.InteracCITx) *scotia.RequestPaymentRequest {
	expireDate := time.Now().Add(time.Hour * time.Duration(p.scotiaProfile.ExpireInHours))

	req := &scotia.RequestPaymentRequest{
		RequestData: &scotia.RequestPaymentRequestData{
			ProductCode:                    "DOMESTIC",
			MessageIdentification:          interacTx.EndToEndId,
			EndToEndIdentification:         interacTx.EndToEndId,
			CreditDebitIndicator:           "CRDT",
			CreationDatetime:               (*scotia.ScotiaDatetime)(interacTx.CreatedAt),
			PaymentExpiryDate:              (*scotia.ScotiaDatetime)(&expireDate),
			SuppressResponderNotifications: p.scotiaProfile.SupressResponderNotifications,
			ReturnUrl:                      "string",
			Language:                       "EN",
			InstructedAmtData: &scotia.ScotiaAmtData{
				Amount:   scotia.ScotiaAmount(interacTx.Amt.Amount),
				Currency: interacTx.Amt.Currency,
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
				Name:               interacTx.CashInAcc.GetLegalName(),
				CountryOfResidence: p.scotiaProfile.CountryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: interacTx.CashInAcc.Email,
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
