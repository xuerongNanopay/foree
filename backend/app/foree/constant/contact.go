package foree_constant

import "xue.io/go-pay/app/foree/account"

var AllowContactAccountType = map[account.ContactAccountType]bool{
	account.ContactAccountTypeCash:               true,
	account.ContactAccountTypeAccountTransfers:   true,
	account.ContactAccountTypeThirdPartyPayments: true,
}

var AllowRelationshipToContactTypes = map[string]bool{
	"Extended Family":  true,
	"Friend":           true,
	"Immediate Family": true,
	"Other":            true,
	"Self":             true,
}
