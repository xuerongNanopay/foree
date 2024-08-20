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
	scotiaProfile ScotiaProfile
	scotiaClient  *scotia.ScotiaClient
	foreeTxRepo   *transaction.ForeeTxRepo
	db            *sql.DB
}

func (p *CITxProcessor) requestPayment(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	// Lock transaction and safety check.
	nForeeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, tx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	if nForeeTx.CurStage != transaction.TxStageInteracCI && nForeeTx.CurStageStatus != transaction.TxStatusInitial {
		return nil, fmt.Errorf("transaction `%v` is in status `%s` at stage `%s`", nForeeTx.ID, nForeeTx.CurStageStatus, nForeeTx.Status)
	}

	// API call
	//
	return nil, nil
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

func (p *CITxProcessor) waitPaymentReceive(ctx context.Context, tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *CITxProcessor) doWait(ctx context.Context, tx transaction.ForeeTx) {
	return
}

func generateScotiaIdentification()
