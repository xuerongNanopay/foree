package service

import (
	"database/sql"

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

func (p *CITxProcessor) generateValidateTransferReq(tx transaction.ForeeTx) (*idm.IDMRequest, error) {
	return &idm.IDMRequest{
		BillingFirstName: "",
	}, nil
}
