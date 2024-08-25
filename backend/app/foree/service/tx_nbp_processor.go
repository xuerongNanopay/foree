package service

import (
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/partner/nbp"
)

type NBPTxProcessor struct {
	db                     *sql.DB
	foreeTxRepo            *transaction.ForeeTxRepo
	txProcessor            *TxProcessor
	idmTxRepo              *transaction.IdmTxRepo
	nbpClient              nbp.NBPClient
	userExtraRepo          *foree_auth.UserExtraRepo
	userIdentificationRepo *foree_auth.UserIdentificationRepo
}

func (p *NBPTxProcessor) processTx(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {
	return nil, nil
}

func (p *NBPTxProcessor) sendPayment(fTx transaction.ForeeTx) (*transaction.ForeeTx, error) {

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
