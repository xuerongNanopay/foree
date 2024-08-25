package service

import (
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/partner/nbp"
)

type NBPTxProcessor struct {
	db          *sql.DB
	foreeTxRepo *transaction.ForeeTxRepo
	txProcessor *TxProcessor
	idmTxRepo   *transaction.IdmTxRepo
	nbpClient   nbp.NBPClient
}

func (p *NBPTxProcessor) processTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *NBPTxProcessor) sendPayment(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {

	return nil, nil
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

func (p *NBPTxProcessor) buildLoadRemittanceRequest(tx transaction.ForeeTx) (*nbp.LoadRemittanceRequest, error) {
	transactionDate := time.Now()
	return &nbp.LoadRemittanceRequest{
		Amount:             nbp.NBPAmount(tx.COUT.Amt.Amount),
		Currency:           tx.COUT.Amt.Currency,
		TransactionDate:    (*nbp.NBPDate)(&transactionDate),
		OriginatingCountry: "Canada",
		PurposeRemittance:  tx.TransactionPurpose,
		RemitterName:       tx.CI.CashInAcc.GetLegalName(),
		RemitterEmail:      tx.CI.CashInAcc.Email,
		RemitterContact:    tx.CI.CashInAcc.PhoneNumber,
		RemitterDOB:        (*nbp.NBPDate)(&tx.Owner.Dob),
		RemitterAddress:    generateLoadRemittanceFromInteracAccount(tx.CI.CashInAcc),
		RemitterIdType:     nbp.RemitterIdTypeOther,
		// RemitterPOB: tx.Owner,
	}, nil
}

func generateLoadRemittanceFromInteracAccount(acc *account.InteracAccount) string {
	if acc.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", acc.Address1, acc.City, acc.Province, acc.PostalCode, acc.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", acc.Address1, acc.Address2, acc.City, acc.Province, acc.PostalCode, acc.Country)
}
