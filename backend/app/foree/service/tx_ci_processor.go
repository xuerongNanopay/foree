package service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
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
	scotiaProfile  ScotiaProfile
	scotiaClient   scotia.ScotiaClient
	interacTxRepo  *transaction.InteracCITxRepo
	foreeTxRepo    *transaction.ForeeTxRepo
	txSummaryRepo  *transaction.TxSummaryRepo
	db             *sql.DB
	waitChan       chan transaction.ForeeTx
	processingMap  []map[int64]*transaction.ForeeTx // Avoid duplicate process
	processingLock sync.RWMutex
}

func (p *CITxProcessor) start() error {
	go p.startProcessLoop()
	return nil
}

// timer, checker,
func (p *CITxProcessor) startProcessLoop() {
	for {
		select {
		case _ = <-p.waitChan:
			//if tx in processinMap

		}
	}

	//TODO: log fata if code reach here.
}

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
		return nil, fmt.Errorf("transaction `%v` is in status `%s` at stage `%s`", nForeeTx.ID, nForeeTx.CurStageStatus, nForeeTx.Status)
	}

	resp, err := p.scotiaClient.RequestPayment(*p.createRequestPaymentReq(tx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		//TODO: logging?
		dTx.Rollback()
		return nil, fmt.Errorf("scotial request payment error: (httpCode: `%v`, request: `%s`, response: `%s`)", resp.StatusCode, resp.RawRequest, resp.RawResponse)
	}

	//TODO: log success

	tx.CI.ScotiaPaymentId = resp.Data.PaymentId

	// Get url payment link
	statusResp, err := p.scotiaClient.PaymentStatus(scotia.PaymentStatusRequest{
		PaymentId:  tx.CI.ScotiaPaymentId,
		EndToEndId: tx.Summary.NBPReference,
	})
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if statusResp.StatusCode/100 != 2 {
		//TODO: logging?
		dTx.Rollback()
		return nil, fmt.Errorf("scotial payment status error: (httpCode: `%v`, request: `%s`, response: `%s`)", statusResp.StatusCode, statusResp.RawRequest, statusResp.RawResponse)
	}

	if len(statusResp.PaymentStatuses) != 1 {
		dTx.Rollback()
		return nil, fmt.Errorf("scotial payment status error: (httpCode: `%v`, request: `%s`, response: `%s`)", statusResp.StatusCode, statusResp.RawRequest, statusResp.RawResponse)
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
				Name:               tx.CI.SrcInteracAcc.GetLegalName(),
				CountryOfResidence: p.scotiaProfile.countryOfResidence,
				ContactDetails: &scotia.ContactDetailsData{
					EmailAddress: tx.CI.SrcInteracAcc.Email,
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

func (p *CITxProcessor) waitPaymentConfirm(tx transaction.ForeeTx) {
	go func() {
		p.waitChan <- tx
	}()
}
