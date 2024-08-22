package service

import (
	"database/sql"

	"xue.io/go-pay/app/foree/transaction"
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

func (p *CITxProcessor) generateValidateTransferReq(tx transaction.ForeeTx) {

}
