package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/partner/scotia"
)

const DefaultScotiaProfileName = "FOREE"

type ScotiaProfile struct {
	Name                          string
	legalName                     string
	companyName                   string
	initiatingPartyName           string
	expireInHours                 int
	supressResponderNotifications bool
	interacId                     string
	countryOfResidence            string
	profileName                   string
	email                         string
	accountNumber                 string
	accountCurrency               string
	amountModificationAllowed     bool
	earlyPaymentAllowed           bool
	guaranteedPaymentRequested    bool
}

type CITxProcessor struct {
	db            *sql.DB
	scotiaProfile ScotiaProfile
	scotiaClient  scotia.ScotiaClient
	interacTxRepo *transaction.InteracCITxRepo
	foreeTxRepo   *transaction.ForeeTxRepo
	txSummaryRepo *transaction.TxSummaryRepo
	txProcessor   *TxProcessor
	fTxs          map[int64]*transaction.ForeeTx
	webhookChan   chan int64
	clearChan     chan int64
	forwardChan   chan transaction.ForeeTx
	startChan     chan transaction.ForeeTx
	ticker        time.Ticker
}

// Loading from DB at beginning. OR, let foree processor do it.

func (p *CITxProcessor) start() error {
	go p.startProcessor()
	return nil
}

func (p *CITxProcessor) startProcessor() error {
	for {
		select {
		case fTx := <-p.startChan:
			_, ok := p.fTxs[fTx.ID]
			if ok {
				//Log duplicate
			} else {
				p.fTxs[fTx.ID] = &fTx
			}
		case fTxId := <-p.clearChan:
			delete(p.fTxs, fTxId)
		case fTx := <-p.forwardChan:
			_, ok := p.fTxs[fTx.ID]
			if !ok {
				//Log miss
			} else {
				delete(p.fTxs, fTx.ID)
			}
			go func() {
				_, err := p.txProcessor.processTx(fTx)
				if err != nil {
					//log err
				}
			}()
		case foreeTxId := <-p.webhookChan:
			v, ok := p.fTxs[foreeTxId]
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
			for _, tx := range p.fTxs {
				func() {
					nTx, err := p.refreshScotiaStatus(*tx)
					if err != nil {
						//Log error
						return
					}
					if nTx.CurStageStatus != tx.CurStageStatus {
						p.forwardChan <- *nTx
					}
				}()
			}
		}
	}
}

func (p *CITxProcessor) startFTx(fTx transaction.ForeeTx) {
	p.startChan <- fTx
}

func (p *CITxProcessor) forwardFTx(fTx transaction.ForeeTx) {
	p.forwardChan <- fTx
}

func (p *CITxProcessor) processTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	t, err := p.requestPayment(fTx)
	if err != nil {
		return nil, err
	}

	p.startChan <- *t

	return t, nil
}

// Scotia APi Call
func (p *CITxProcessor) requestPayment(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	// Lock transaction and safety check.
	nForeeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, tx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	if nForeeTx.CurStage != transaction.TxStageInteracCI && nForeeTx.CurStageStatus != transaction.TxStatusInitial {
		dTx.Rollback()
		return nil, fmt.Errorf("CI failed -- transaction `%v` is in status `%s` at stage `%s`", nForeeTx.ID, nForeeTx.CurStageStatus, nForeeTx.Status)
	}

	resp, err := p.scotiaClient.RequestPayment(*p.createRequestPaymentReq(tx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		//TODO: logging?
		dTx.Rollback()
		return nil, fmt.Errorf("CI failed -- scotial requestPayment error: (httpCode: `%v`, request: `%s`, response: `%s`)", resp.StatusCode, resp.RawRequest, resp.RawResponse)
	}

	//TODO: log success

	tx.CI.ScotiaPaymentId = resp.Data.PaymentId

	// Get url payment link
	statusResp, err := p.scotiaClient.PaymentStatus(scotia.PaymentStatusRequest{
		PaymentId:  tx.CI.ScotiaPaymentId,
		EndToEndId: tx.CI.EndToEndId,
	})
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if statusResp.StatusCode/100 != 2 {
		//TODO: logging?
		dTx.Rollback()
		return nil, fmt.Errorf("CI failed -- scotial paymentstatus error: (httpCode: `%v`, request: `%s`, response: `%s`)", statusResp.StatusCode, statusResp.RawRequest, statusResp.RawResponse)
	}

	if len(statusResp.PaymentStatuses) != 1 {
		dTx.Rollback()
		return nil, fmt.Errorf("CI failed -- scotial paymentstatus error: (httpCode: `%v`, request: `%s`, response: `%s`)", statusResp.StatusCode, statusResp.RawRequest, statusResp.RawResponse)
	}

	// Update CI
	tx.CI.PaymentUrl = statusResp.PaymentStatuses[0].GatewayUrl
	tx.CI.ScotiaPaymentId = resp.Data.PaymentId
	tx.CI.Status = transaction.TxStatusSent
	tx.CI.ScotiaClearingReference = resp.Data.ClearingSystemReference

	// Update Foree
	tx.CurStageStatus = transaction.TxStatusSent
	tx.CI.ScotiaPaymentId = resp.Data.PaymentId

	// Update summary
	tx.Summary.Status = transaction.TxSummaryStatusAwaitPayment
	tx.Summary.PaymentUrl = statusResp.PaymentStatuses[0].GatewayUrl

	err = p.interacTxRepo.UpdateInteracCITxById(ctx, *tx.CI)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	err = p.foreeTxRepo.UpdateForeeTxById(ctx, tx)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	err = p.txSummaryRepo.UpdateTxSummaryById(ctx, *tx.Summary)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if err = dTx.Commit(); err != nil {
		return nil, err
	}

	return &tx, nil
}

// The function normally in a goroutine
func (p *CITxProcessor) refreshScotiaStatus(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	if fTx.CI.Status != transaction.TxStatusSent {
		p.clearChan <- fTx.ID
		return nil, fmt.Errorf("CI failed -- refreshScotiaStatusAndProcess: InteracCITx `%v` is in `%s`", fTx.CI.ID, fTx.CI.Status)
	}

	detailResp, err := p.scotiaClient.PaymentDetail(scotia.PaymentDetailRequest{
		PaymentId:  fTx.CI.ScotiaPaymentId,
		EndToEndId: fTx.CI.EndToEndId,
	})
	if err != nil {
		return nil, err
	}

	if detailResp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("CI failed -- refreshScotiaStatusAndProcess: scotia paymentdetail error: (httpCode: `%v`, request: `%s`, response: `%s`)", detailResp.StatusCode, detailResp.RawRequest, detailResp.RawResponse)
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
		p.clearChan <- fTx.ID
		dTx.Rollback()
		return nil, fmt.Errorf("CI failed -- refreshScotiaStatusAndProcess: ForeeTx `%v` is in stage `%s` at status `%s`", curFTx.ID, curFTx.CurStage, curFTx.CurStageStatus)
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

func (p *CITxProcessor) cleanTx(foreeTxId int64) {
	p.clearChan <- foreeTxId
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
	expireDate := time.Now().Add(time.Hour * time.Duration(p.scotiaProfile.expireInHours))

	req := &scotia.RequestPaymentRequest{
		RequestData: &scotia.RequestPaymentRequestData{
			ProductCode:                    "DOMESTIC",
			MessageIdentification:          tx.Summary.NBPReference,
			EndToEndIdentification:         tx.Summary.NBPReference,
			CreditDebitIndicator:           "CRDT",
			CreationDatetime:               (*scotia.ScotiaDatetime)(&tx.CreateAt),
			PaymentExpiryDate:              (*scotia.ScotiaDatetime)(&expireDate),
			SuppressResponderNotifications: p.scotiaProfile.supressResponderNotifications,
			ReturnUrl:                      "string",
			Language:                       "EN",
			InstructedAmtData: &scotia.ScotiaAmtData{
				Amount:   scotia.ScotiaAmount(tx.CI.Amt.Amount),
				Currency: tx.CI.Amt.Currency,
			},
			InitiatingParty: &scotia.InitiatingPartyData{
				Name: p.scotiaProfile.initiatingPartyName,
				Identification: &scotia.IdentificationData{
					OrganisationIdentification: &scotia.OrganisationIdentificationData{
						Other: []scotia.OtherData{
							{
								Identification: p.scotiaProfile.interacId,
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
				CountryOfResidence: p.scotiaProfile.countryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: tx.CI.CashInAcc.Email,
				},
			},
			Creditor: &scotia.CreditorData{
				Name:               p.scotiaProfile.legalName,
				CountryOfResidence: p.scotiaProfile.countryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: p.scotiaProfile.email,
				},
			},
			UltimateCreditor: &scotia.CreditorData{
				Name:               p.scotiaProfile.legalName,
				CountryOfResidence: p.scotiaProfile.countryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: p.scotiaProfile.email,
				},
			},
			CreditorAccount: &scotia.CreditorAccountData{
				Identification: p.scotiaProfile.accountNumber,
				Currency:       p.scotiaProfile.accountCurrency,
				SchemeName:     "ALIAS_ACCT_NO",
			},
			FraudSupplementaryInfo: &scotia.FraudSupplementaryInfoData{
				CustomerAuthenticationMethod: "PASSWORD",
			},
			PaymentCondition: &scotia.PaymentConditionData{
				AmountModificationAllowed:  p.scotiaProfile.amountModificationAllowed,
				EarlyPaymentAllowed:        p.scotiaProfile.earlyPaymentAllowed,
				GuaranteedPaymentRequested: p.scotiaProfile.guaranteedPaymentRequested,
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
