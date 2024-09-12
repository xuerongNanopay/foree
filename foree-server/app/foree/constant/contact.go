package foree_constant

import "xue.io/go-pay/app/foree/account"

const (
	ContactAccountTypeCash         account.ContactAccountType = "CASH_PICKUP"
	ContactAccountTypeBankAccount  account.ContactAccountType = "BANK_ACCOUNT"
	ContactAccountTypeMobileWallet account.ContactAccountType = "MOBILE_WALLET"
	ContactAccountTypeRDA          account.ContactAccountType = "ROSHAN_DIGITAL_ACCOUNT"
)

var AllowContactAccountType = map[account.ContactAccountType]bool{
	ContactAccountTypeCash:         true,
	ContactAccountTypeBankAccount:  true,
	ContactAccountTypeMobileWallet: true,
	ContactAccountTypeRDA:          true,
}

var AllowRelationshipToContactTypes = map[string]bool{
	"Extended Family":  true,
	"Friend":           true,
	"Immediate Family": true,
	"Other":            true,
	"Self":             true,
}
