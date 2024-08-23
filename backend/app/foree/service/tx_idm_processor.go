package service

import (
	"database/sql"
	"fmt"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/partner/idm"
)

type IDMTxProcessor struct {
	db          *sql.DB
	txProcessor *TxProcessor
}

func (p *IDMTxProcessor) createAndProcessTx(tx transaction.ForeeTx) {

}

func (p *CITxProcessor) idmTransferValidate(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

// TODO
// RemitterPOB: "TODO add property",
// Identification
func (p *CITxProcessor) generateValidateTransferReq(tx transaction.ForeeTx) (*idm.IDMRequest, error) {
	IsCashPickup := false
	if tx.COUT.CashOutAcc.Type == account.ContactAccountTypeCash {
		IsCashPickup = true
	}

	beneBankName := ""
	if tx.COUT.CashOutAcc.Type != account.ContactAccountTypeCash {
		beneBankName = tx.COUT.CashOutAcc.InstitutionName
	}

	return &idm.IDMRequest{
		BillingFirstName:        tx.Owner.FirstName,
		BillingMiddleName:       tx.Owner.MiddleName,
		BillingLastName:         tx.Owner.LastName,
		RemitterOccupation:      "TODO",
		BillingStreet:           fmt.Sprintf("%s %s", tx.Owner.Address1, tx.Owner.Address2),
		BillingCity:             tx.Owner.City,
		BillingState:            tx.Owner.Province,
		BillingPostalCode:       tx.Owner.PostalCode,
		BillingCountry:          tx.Owner.Country,
		PhoneNumber:             tx.CI.CashInAcc.PhoneNumber,
		UserEmail:               tx.CI.CashInAcc.Email,
		Dob:                     (*idm.IDMDate)(&tx.Owner.Dob),
		Nationality:             tx.Owner.Nationality,
		SrcDigitalAccNOHash:     tx.CI.CashInAcc.AccountHash,
		ShippingFirstName:       tx.COUT.CashOutAcc.FirstName,
		ShippingMiddleName:      tx.COUT.CashOutAcc.MiddleName,
		ShippingLastName:        tx.COUT.CashOutAcc.LastName,
		IsCashPickup:            IsCashPickup,
		DestDigitalAccNOHash:    tx.COUT.CashOutAcc.AccountHash,
		BeneBankName:            beneBankName,
		DestPhoneNumber:         tx.COUT.CashOutAcc.PhoneNumber,
		SRRelationship:          tx.COUT.CashOutAcc.RelationshipToContact,
		PurposeOfTransfer:       tx.TransactionPurpose,
		TransactionCreationTime: tx.CreateAt.UnixMilli(),
		Amount:                  idm.IDMAmount(tx.SrcAmt.Amount),
		Currency:                tx.SrcAmt.Currency,
		PayoutAmount:            idm.IDMAmount(tx.DestAmt.Amount),
		PayoutCurrency:          tx.DestAmt.Currency,
		TransactionIdentifier:   fmt.Sprintf("%012d", tx.ID),
		TransactionRefId:        tx.Summary.NBPReference,
		Ip:                      tx.Ip,
		SrcAccountIdentifier:    fmt.Sprintf("%09d", tx.CI.CashInAccId),
		DestAccountIdentifier:   fmt.Sprintf("%09d", tx.COUT.CashOutAccId),
	}, nil
}
