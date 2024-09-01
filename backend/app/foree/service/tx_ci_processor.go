package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	ticker        time.Ticker
}

// Loading from DB at beginning. OR, let foree processor do it.

func (p *CITxProcessor) start() error {
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
		return nil, fmt.Errorf("CITxProcessor -- transaction `%v` is in status `%s` at stage `%s`", nForeeTx.ID, nForeeTx.CurStageStatus, nForeeTx.Status)
	}

	resp, err := p.scotiaClient.RequestPayment(*p.createRequestPaymentReq(fTx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		//TODO: logging?
		dTx.Rollback()
		return nil, fmt.Errorf("CITxProcessor -- scotial requestPayment error: (httpCode: `%v`, request: `%s`, response: `%s`)", resp.StatusCode, resp.RawRequest, resp.RawResponse)
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

func (p *CITxProcessor) createRequestPaymentReq(tx transaction.ForeeTx) *scotia.RequestPaymentRequest {
	expireDate := time.Now().Add(time.Hour * time.Duration(p.scotiaProfile.ExpireInHours))

	req := &scotia.RequestPaymentRequest{
		RequestData: &scotia.RequestPaymentRequestData{
			ProductCode:                    "DOMESTIC",
			MessageIdentification:          tx.Summary.NBPReference,
			EndToEndIdentification:         tx.Summary.NBPReference,
			CreditDebitIndicator:           "CRDT",
			CreationDatetime:               (*scotia.ScotiaDatetime)(&tx.CreatedAt),
			PaymentExpiryDate:              (*scotia.ScotiaDatetime)(&expireDate),
			SuppressResponderNotifications: p.scotiaProfile.SupressResponderNotifications,
			ReturnUrl:                      "string",
			Language:                       "EN",
			InstructedAmtData: &scotia.ScotiaAmtData{
				Amount:   scotia.ScotiaAmount(tx.CI.Amt.Amount),
				Currency: tx.CI.Amt.Currency,
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
				Name:               tx.CI.CashInAcc.GetLegalName(),
				CountryOfResidence: p.scotiaProfile.CountryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: tx.CI.CashInAcc.Email,
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
