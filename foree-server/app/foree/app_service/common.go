package foree_service

import (
	"fmt"

	"xue.io/go-pay/app/foree/account"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/partner/nbp"
)

func mapNBPMode(contactAccount *account.ContactAccount) (nbp.PMTMode, error) {
	switch contactAccount.Type {
	case foree_constant.ContactAccountTypeCash:
		return nbp.PMTModeCash, nil
	case foree_constant.ContactAccountTypeBankAccount:
		if contactAccount.InstitutionName == "NBP" {
			return nbp.PMTModeAccountTransfers, nil
		} else {
			return nbp.PMTModeThirdPartyPayments, nil
		}
	case foree_constant.ContactAccountTypeRDA:
		fallthrough
	case foree_constant.ContactAccountTypeMobileWallet:
		return nbp.PMTModeThirdPartyPayments, nil
	default:
		return "", fmt.Errorf("NBPTxProcessor -- unknown contact account type `%s`", contactAccount.Type)
	}
}
