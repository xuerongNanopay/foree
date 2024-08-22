package service

import (
	"database/sql"
	"fmt"

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
	return &idm.IDMRequest{
		BillingFirstName:   tx.Owner.FirstName,
		BillingMiddleName:  tx.Owner.MiddleName,
		BillingLastname:    tx.Owner.LastName,
		RemitterOccupation: "TODO",
		BillingStreet:      fmt.Sprintf("%s %s", tx.Owner.Address1, tx.Owner.Address2),
		BillingCity:        tx.Owner.City,
		BillingState:       tx.Owner.Province,
		BillingPostalCode:  tx.Owner.PostalCode,
		BillingCountry:     tx.Owner.Country,
		PhoneNumber:        tx.CI.SrcInteracAcc.PhoneNumber,
		UserEmail:          tx.CI.SrcInteracAcc.Email,
		Dob:                (*idm.IDMDate)(&tx.Owner.Dob),
		Nationality:        tx.Owner.Nationality,
	}, nil
}
