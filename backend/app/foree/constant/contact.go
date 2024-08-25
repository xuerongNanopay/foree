package foree_constant

import "xue.io/go-pay/app/foree/account"

const (
	ContactAccountTypeCash               account.ContactAccountType = "CASH"
	ContactAccountTypeAccountTransfers   account.ContactAccountType = "ACCOUNT_TRANSFERS"
	ContactAccountTypeThirdPartyPayments account.ContactAccountType = "THIRD_PARTY_PAYMENTS"
)

var AllowContactAccountType = map[account.ContactAccountType]bool{
	ContactAccountTypeCash:               true,
	ContactAccountTypeAccountTransfers:   true,
	ContactAccountTypeThirdPartyPayments: true,
}

var AllowRelationshipToContactTypes = map[string]bool{
	"Extended Family":  true,
	"Friend":           true,
	"Immediate Family": true,
	"Other":            true,
	"Self":             true,
}
