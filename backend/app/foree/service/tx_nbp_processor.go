package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/partner/nbp"
)

type NBPTxProcessor struct {
	db                     *sql.DB
	foreeTxRepo            *transaction.ForeeTxRepo
	txProcessor            *TxProcessor
	nbpTxRepo              *transaction.NBPCOTxRepo
	nbpClient              nbp.NBPClient
	userExtraRepo          *foree_auth.UserExtraRepo
	userIdentificationRepo *foree_auth.UserIdentificationRepo
	retryFTxs              map[int64]*transaction.ForeeTx
	waitFTxs               map[int64]*transaction.ForeeTx
	retryChan              chan transaction.ForeeTx
	waitChan               chan transaction.ForeeTx
	forwardChan            chan transaction.ForeeTx
	retryTicker            time.Ticker
	checkStatusTicker      time.Ticker
	clearChan              chan int64
}

func (p *NBPTxProcessor) start() error {
	go p.startProcessor()
	return nil
}

func (p *NBPTxProcessor) startProcessor() {
	for {
		select {
		case fTx := <-p.waitChan:
			_, ok := p.waitFTxs[fTx.ID]
			if ok {
				//Log duplicate
			} else {
				p.waitFTxs[fTx.ID] = &fTx
			}
		case fTx := <-p.forwardChan:
			_, ok := p.waitFTxs[fTx.ID]
			if !ok {
				//Log miss
			} else {
				delete(p.waitFTxs, fTx.ID)
			}
			go func() {
				_, err := p.txProcessor.processTx(fTx)
				if err != nil {
					//log err
				}
			}()
		case foreeTxId := <-p.clearChan:
			delete(p.waitFTxs, foreeTxId)
		case <-p.checkStatusTicker.C:
			req := nbp.TransactionStatusByIdsRequest{
				Ids: nbp.NBPIds{},
			}
			for _, tx := range p.waitFTxs {
				req.Ids = append(req.Ids, tx.COUT.NBPReference)
			}

			resp, err := p.nbpClient.TransactionStatusByIds(req)
			if err != nil {
				//TODO: Log
				return
			}
			m := make(map[string]nbp.TransactionStatus, len(resp.TransactionStatuses))
			for _, v := range resp.TransactionStatuses {
				m[v.GlobalId] = v
			}

			for _, fTx := range p.waitFTxs {
				func() {
					s, ok := m[fTx.COUT.NBPReference]
					if !ok {
						//log: error
						return
					}
					nTx, err := p.refreshNBPStatus(*fTx, s.Status)
					if err != nil {
						//Log error
						return
					}

					if nTx.CurStageStatus != fTx.CurStageStatus {
						p.forwardChan <- *nTx
					}
				}()
			}
		}
	}
}

// We don't use transaction here, case NBP can check duplicate.
func (p *NBPTxProcessor) processTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	// Safe check.
	if fTx.CurStage != transaction.TxStageNBPCO && fTx.CurStageStatus != transaction.TxStatusInitial {
		return nil, fmt.Errorf("NBPTxProcessor -- transaction `%v` is in status `%s` at stage `%s`", fTx.ID, fTx.CurStageStatus, fTx.Status)
	}

	req, err := p.buildLoadRemittanceRequest(fTx)
	if err != nil {
		return nil, err
	}
	mode, err := mapNBPMode(fTx.COUT.CashOutAcc.Type)
	if err != nil {
		return nil, err
	}

	var resp *nbp.LoadRemittanceResponse

	// Retry 5 times with 15 second interval.
	for i := 0; i < 5; i++ {
		resp, err = p.sendPaymentWithMode(*req, mode)
		if err != nil {
			return nil, err
		}
		//Retry case: 5xx, 401, 403
		if resp.StatusCode/100 == 5 || resp.ResponseCode == "401" || resp.ResponseCode == "403" || resp.ResponseCode == "406" {
			time.Sleep(15 * time.Second)
		} else {
			break
		}

	}

	dTx, err := p.db.Begin()
	if err != nil {
		dTx.Rollback()
		//TODO: log err
		return nil, err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	if resp.StatusCode/100 == 5 || resp.ResponseCode == "401" || resp.ResponseCode == "403" {
		fTx.COUT.Status = transaction.TxStatusPending
		fTx.CurStageStatus = transaction.TxStatusPending
	} else if resp.StatusCode/100 == 2 || resp.ResponseCode == "405" {
		fTx.COUT.Status = transaction.TxStatusSent
		fTx.CurStageStatus = transaction.TxStatusSent
	} else {
		fTx.COUT.Status = transaction.TxStatusRejected
		fTx.CurStageStatus = transaction.TxStatusRejected
	}

	err = p.nbpTxRepo.UpdateNBPCOTxById(ctx, *fTx.COUT)
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

	if resp.StatusCode/100 == 5 || resp.ResponseCode == "401" || resp.ResponseCode == "403" {
		p.retryChan <- fTx
	} else if resp.StatusCode/100 == 2 || resp.ResponseCode == "405" {
		p.waitChan <- fTx
	}

	return &fTx, nil
}

func (p *NBPTxProcessor) refreshNBPStatus(fTx transaction.ForeeTx, nbpStatus string) (*transaction.ForeeTx, error) {
	newStatus := nbpToInternalStatusMapper(nbpStatus)
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
		dTx.Rollback()
		p.clearChan <- fTx.ID
		return nil, fmt.Errorf("NBPTxProcessor -- refreshNBPStatus -- ForeeTx `%v` is in stage `%s` at status `%s`", curFTx.ID, curFTx.CurStage, curFTx.CurStageStatus)
	}

	fTx.COUT.Status = newStatus
	fTx.CurStageStatus = newStatus

	err = p.nbpTxRepo.UpdateNBPCOTxById(ctx, *fTx.COUT)
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

func (p *NBPTxProcessor) buildLoadRemittanceRequest(fTx transaction.ForeeTx) (*nbp.LoadRemittanceRequest, error) {
	userExtra, err := p.userExtraRepo.GetUniqueUserExtraByOwnerId(fTx.OwnerId)
	if err != nil {
		return nil, err
	}
	identifications, err := p.userIdentificationRepo.GetAllUserIdentificationByOwnerId(fTx.OwnerId)
	if err != nil {
		return nil, err
	}

	var identification *foree_auth.UserIdentification = nil

	for _, v := range identifications {
		if v.Status == foree_auth.IdentificationStatusActive {
			identification = v
		}
	}

	if identification == nil {
		return nil, fmt.Errorf("NBPTxProcessor -- user `%v` do not find a proper identification", fTx.OwnerId)
	}

	transactionDate := time.Now()
	lrr := &nbp.LoadRemittanceRequest{
		GlobalId:                        fTx.COUT.NBPReference,
		Amount:                          nbp.NBPAmount(fTx.COUT.Amt.Amount),
		Currency:                        fTx.COUT.Amt.Currency,
		TransactionDate:                 (*nbp.NBPDate)(&transactionDate),
		OriginatingCountry:              "Canada",
		PurposeRemittance:               fTx.TransactionPurpose,
		RemitterName:                    fTx.CI.CashInAcc.GetLegalName(),
		RemitterEmail:                   fTx.CI.CashInAcc.Email,
		RemitterContact:                 fTx.CI.CashInAcc.PhoneNumber,
		RemitterDOB:                     (*nbp.NBPDate)(&fTx.Owner.Dob),
		RemitterAddress:                 generateLoadRemittanceFromInteracAccount(fTx.CI.CashInAcc),
		RemitterIdType:                  mapNBPRemitterIdType(identification.Type),
		RemitterId:                      identification.Value,
		RemitterPOB:                     userExtra.Pob,
		BeneficiaryName:                 fTx.COUT.CashOutAcc.GetLegalName(),
		BeneficiaryAddress:              generateLoadRemittanceFromContactAccount(fTx.COUT.CashOutAcc),
		BeneficiaryCity:                 fTx.COUT.CashOutAcc.City,
		RemitterBeneficiaryRelationship: fTx.COUT.CashOutAcc.RelationshipToContact,
	}

	if fTx.COUT.CashOutAcc.Type != foree_constant.ContactAccountTypeCash {
		lrr.BeneficiaryBank = fTx.COUT.CashOutAcc.InstitutionName
		lrr.BeneficiaryAccount = fTx.COUT.CashOutAcc.AccountNumber
	}

	return lrr, nil

}

func nbpToInternalStatusMapper(nbpStatus string) transaction.TxStatus {
	switch nbpStatus {
	case "REJECTED":
		fallthrough
	case "ERROR":
		return transaction.TxStatusRejected
	case "PAID":
		return transaction.TxStatusCompleted
	case "CANCELLED":
		return transaction.TxStatusCancelled
	default:
		return transaction.TxStatusSent
	}
}

func mapNBPMode(accType account.ContactAccountType) (nbp.PMTMode, error) {
	switch accType {
	case foree_constant.ContactAccountTypeCash:
		return nbp.PMTModeCash, nil
	case foree_constant.ContactAccountTypeThirdPartyPayments:
		return nbp.PMTModeThirdPartyPayments, nil
	case foree_constant.ContactAccountTypeAccountTransfers:
		return nbp.PMTModeAccountTransfers, nil
	default:
		return "", fmt.Errorf("NBPTxProcessor -- unknown contact account type `%s`", accType)
	}
}

func mapNBPRemitterIdType(idType foree_auth.IdentificationType) nbp.RemitterIdType {
	switch idType {
	case foree_constant.IDTypePassport:
		return nbp.RemitterIdTypePassport
	case foree_constant.IDTypeDriverLicense:
		return nbp.RemitterIdTypeDrivinglicense
	default:
		return nbp.RemitterIdTypeOther
	}
}

func generateLoadRemittanceFromInteracAccount(acc *account.InteracAccount) string {
	if acc.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", acc.Address1, acc.City, acc.Province, acc.PostalCode, acc.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", acc.Address1, acc.Address2, acc.City, acc.Province, acc.PostalCode, acc.Country)
}

func generateLoadRemittanceFromContactAccount(acc *account.ContactAccount) string {
	if acc.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", acc.Address1, acc.City, acc.Province, acc.PostalCode, acc.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", acc.Address1, acc.Address2, acc.City, acc.Province, acc.PostalCode, acc.Country)
}
