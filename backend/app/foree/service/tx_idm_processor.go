package service

import (
	"context"
	"database/sql"
	"fmt"

	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/partner/idm"
)

type IDMTxProcessor struct {
	db          *sql.DB
	foreeTxRepo *transaction.ForeeTxRepo
	idmTxRepo   *transaction.IdmTxRepo
	txProcessor *TxProcessor
	idmClient   idm.IDMClient
}

// IDM API called
func (p *IDMTxProcessor) processTx(tx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	req, err := p.generateValidateTransferReq(tx)
	if err != nil {
		return nil, err
	}

	// Safe check
	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	nForeeTx, err := p.foreeTxRepo.GetUniqueForeeTxForUpdateById(ctx, tx.ID)
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}

	if nForeeTx.CurStage != transaction.TxStageIDM && nForeeTx.CurStageStatus != transaction.TxStatusInitial {
		dTx.Rollback()
		return nil, fmt.Errorf("IDM failed: transaction `%v` is in status `%s` at stage `%s`", nForeeTx.ID, nForeeTx.CurStageStatus, nForeeTx.Status)
	}

	resp, err := p.idmClient.Transfer(*req)
	if err != nil {
		dTx.Rollback()
		return nil, err
	}

	if resp.StatusCode/100 == 2 && resp.GetResultStatus() == "ACCEPT" {
		//TODO: log success
		tx.CurStageStatus = transaction.TxStatusCompleted
		tx.IDM.Status = transaction.TxStatusCompleted
		err = p.idmTxRepo.UpdateIDMTxById(ctx, *tx.IDM)
		if err != nil {
			dTx.Rollback()
			return nil, err
		}
		err = p.foreeTxRepo.UpdateForeeTxById(ctx, tx)
		if err != nil {
			dTx.Rollback()
			return nil, err
		}
	} else {
		//TODO: log fails
		tx.CurStageStatus = transaction.TxStatusSuspend
		tx.IDM.Status = transaction.TxStatusSuspend
		err = p.idmTxRepo.UpdateIDMTxById(ctx, *tx.IDM)
		if err != nil {
			dTx.Rollback()
			return nil, err
		}
		err = p.foreeTxRepo.UpdateForeeTxById(ctx, tx)
		if err != nil {
			dTx.Rollback()
			return nil, err
		}
		//TODO: generate approval.
	}

	if err = dTx.Commit(); err != nil {
		return nil, err
	}

	return &tx, nil
}

func (p *IDMTxProcessor) generateValidateTransferReq(tx transaction.ForeeTx) (*idm.IDMRequest, error) {
	IsCashPickup := false
	if tx.COUT.CashOutAcc.Type == foree_constant.ContactAccountTypeCash {
		IsCashPickup = true
	}

	beneBankName := ""
	if tx.COUT.CashOutAcc.Type != foree_constant.ContactAccountTypeCash {
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
		Nationality:             "TODO",
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
