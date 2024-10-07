package foree_service

import (
	"context"
	"database/sql"
	"fmt"

	foree_constant "xue.io/go-pay/app/foree/constant"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/partner/idm"
)

func NewIDMTxProcessor(
	db *sql.DB,
	foreeTxRepo *transaction.ForeeTxRepo,
	idmTxRepo *transaction.IdmTxRepo,
	idmClient idm.IDMClient,
	txProcessor *TxProcessor,
) *IDMTxProcessor {
	return &IDMTxProcessor{
		db:          db,
		foreeTxRepo: foreeTxRepo,
		idmTxRepo:   idmTxRepo,
		idmClient:   idmClient,
		txProcessor: txProcessor,
	}
}

type IDMTxProcessor struct {
	db          *sql.DB
	foreeTxRepo *transaction.ForeeTxRepo
	idmTxRepo   *transaction.IdmTxRepo
	idmClient   idm.IDMClient
	txProcessor *TxProcessor
}

func (p *IDMTxProcessor) process(parentTxId int64) {
	ctx := context.TODO()
	idmTx, err := p.idmTxRepo.GetUniqueIDMTxByParentTxId(ctx, parentTxId)
	if err != nil {
		foree_logger.Logger.Error("IDMTxProcessor--process_FAIL", "parentTxId", parentTxId, "cause", err.Error())
		return
	}
	if idmTx == nil {
		foree_logger.Logger.Error("IDMTxProcessor--process_FAIL", "parentTxId", parentTxId, "cause", "idmTx no found")
		return
	}
	switch idmTx.Status {
	case transaction.TxStatusInitial:
		p.idmTransferVeirfy(parentTxId)
	case transaction.TxStatusSuspend:
		foree_logger.Logger.Debug("IDMTxProcessor--process", "parentTxId", parentTxId, "idmTxId", idmTx.ID, "idmTxStatus", idmTx.Status, "msg", "waiting for action")
	case transaction.TxStatusCompleted:
		p.txProcessor.next(idmTx.ParentTxId)
	case transaction.TxStatusRejected:
		p.txProcessor.rollback(idmTx.ParentTxId)
	default:
		foree_logger.Logger.Error(
			"IDMTxProcessor--process_FAIL",
			"parentTxId", parentTxId,
			"interacCITxId", idmTx.ID,
			"interacCITxStatus", idmTx.Status,
			"cause", "unsupport status",
		)
	}
}

func (p *IDMTxProcessor) idmTransferVeirfy(parentTxId int64) {
	fTx, err := p.txProcessor.loadTx(parentTxId, true)
	if err != nil {
		foree_logger.Logger.Error("IDMTxProcessor--idmTransferVeirfy_FAIL", "parentTxId", parentTxId, "cause", err.Error())
	}

	req := p.generateValidateTransferReq(fTx)
	resp, err := p.idmClient.Transfer(*req)
	// Treat err and err response as Suspend.
	if err != nil {
		foree_logger.Logger.Error("IDMTxProcessor--idmTransferVeirfy_FAIL",
			"parentTxId", parentTxId,
			"cause", err.Error(),
		)
	}
	if resp.StatusCode/100 != 2 || resp.GetResultStatus() != "ACCEPT" {
		foree_logger.Logger.Warn("InteracTxProcessor-idmTransferVeirfy_FAIL",
			"idmTxId", fTx.IDM.ID,
			"httpResponseStatus", resp.StatusCode,
			"httpRequest", resp.RawRequest,
			"httpResponseBody", resp.RawResponse,
			"cause", "idm response error",
		)
	}
	if err != nil || resp.StatusCode/100 != 2 || resp.GetResultStatus() != "ACCEPT" {
		idm := *fTx.IDM
		idm.Status = transaction.TxStatusSuspend
		err = p.idmTxRepo.UpdateIDMTxById(context.TODO(), idm)
		if err != nil {
			foree_logger.Logger.Error("IDMTxProcessor--idmTransferVeirfy_FAIL",
				"parentTxId", parentTxId,
				"cause", err.Error(),
			)
			return
		}
		foree_logger.Logger.Warn("IDMTxProcessor--idmTransferVeirfy",
			"parentTxId", parentTxId,
			"idmTxId", fTx.IDM.ID,
			"msg", "idm did not approve",
		)
		return
	}

	idm := *fTx.IDM
	idm.Status = transaction.TxStatusCompleted
	err = p.idmTxRepo.UpdateIDMTxById(context.TODO(), idm)
	if err != nil {
		foree_logger.Logger.Error("IDMTxProcessor--idmTransferVeirfy_FAIL",
			"parentTxId", parentTxId,
			"cause", err.Error(),
		)
		return
	}
	foree_logger.Logger.Info("IDMTxProcessor--idmTransferVeirfy",
		"parentTxId", parentTxId,
		"idmTxId", fTx.IDM.ID,
		"msg", "IDM approve",
	)

	p.process(idm.ParentTxId)
}

func (p *IDMTxProcessor) ManualUpdate(parentTxId int64, newTxStatus transaction.TxStatus) (bool, error) {
	if newTxStatus != transaction.TxStatusRejected && newTxStatus != transaction.TxStatusCompleted {
		return false, fmt.Errorf("unsupport transaction status `%v`", newTxStatus)
	}

	ctx := context.TODO()
	idmTx, err := p.idmTxRepo.GetUniqueIDMTxByParentTxId(ctx, parentTxId)
	if err != nil {
		return false, err
	}
	if idmTx == nil {
		return false, fmt.Errorf("idmTx no found with parentTxId `%v`", parentTxId)
	}
	if idmTx.Status != transaction.TxStatusSent {
		return false, fmt.Errorf("expect idmTx in `%v`, but got `%v`", transaction.TxStatusSent, idmTx.Status)
	}

	idmTx.Status = transaction.TxStatusCompleted
	err = p.idmTxRepo.UpdateIDMTxById(context.TODO(), *idmTx)
	if err != nil {
		return false, err
	}
	go p.txProcessor.next(idmTx.ParentTxId)
	return true, nil
}

func (p *IDMTxProcessor) generateValidateTransferReq(tx *transaction.ForeeTx) *idm.IDMRequest {
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
		BillingStreet:           fmt.Sprintf("%s %s", tx.Owner.Address1, tx.Owner.Address2),
		BillingCity:             tx.Owner.City,
		BillingState:            tx.Owner.Province,
		BillingPostalCode:       tx.Owner.PostalCode,
		BillingCountry:          tx.Owner.Country,
		RemitterOccupation:      "TODO",
		PhoneNumber:             tx.Owner.PhoneNumber,
		UserEmail:               tx.CI.CashInAcc.Email,
		Dob:                     (*idm.IDMDate)(tx.Owner.Dob),
		Nationality:             "TODO",
		SrcDigitalAccNOHash:     tx.CI.CashInAcc.Email,
		ShippingFirstName:       tx.COUT.CashOutAcc.FirstName,
		ShippingMiddleName:      tx.COUT.CashOutAcc.MiddleName,
		ShippingLastName:        tx.COUT.CashOutAcc.LastName,
		IsCashPickup:            IsCashPickup,
		DestDigitalAccNOHash:    tx.COUT.CashOutAcc.AccountHash,
		BeneBankName:            beneBankName,
		DestPhoneNumber:         tx.COUT.CashOutAcc.PhoneNumber,
		SRRelationship:          tx.COUT.CashOutAcc.RelationshipToContact,
		PurposeOfTransfer:       tx.TransactionPurpose,
		TransactionCreationTime: tx.CreatedAt.UnixMilli(),
		Amount:                  idm.IDMAmount(tx.SrcAmt.Amount),
		Currency:                tx.SrcAmt.Currency,
		PayoutAmount:            idm.IDMAmount(tx.DestAmt.Amount),
		PayoutCurrency:          tx.DestAmt.Currency,
		TransactionIdentifier:   fmt.Sprintf("%012d", tx.ID),
		TransactionRefId:        tx.COUT.NBPReference,
		Ip:                      tx.Ip,
		SrcAccountIdentifier:    fmt.Sprintf("%09d", tx.CI.CashInAccId),
		DestAccountIdentifier:   fmt.Sprintf("%09d", tx.COUT.CashOutAccId),
	}
}
