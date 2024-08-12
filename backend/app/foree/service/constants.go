package service

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/auth"
)

const UserGroup = "foree-person"

// 3600 * 24 * 365 *19
const Second_In_Year = 31536000

// letters, spaces, number and extended latin
const NameReg = `^[a-zA-Z_0-9\u00C0-\u017F][a-zA-Z_0-9\u00C0-\u017F\s]*$`
const NineDigitReg = `^\\d{9}$`

var phoneNumberReplayer = strings.NewReplacer(" ", "", "(", "", ")", "", "-", "", "+", "")
var validate = validator.New()

var allowIdentificationTypes = map[auth.IdentificationType]bool{
	auth.IDTypePassport:      true,
	auth.IDTypeDriverLicense: true,
	auth.IDTypeProvincalId:   true,
	auth.IDTypeNationId:      true,
}

var allowRelationshipToContactTypes = map[string]bool{
	"Extended Family":  true,
	"Friend":           true,
	"Immediate Family": true,
	"Other":            true,
	"Self":             true,
}

var allowContactAccountType = map[account.ContactAccountType]bool{
	account.ContactAccountTypeCash:               true,
	account.ContactAccountTypeAccountTransfers:   true,
	account.ContactAccountTypeThirdPartyPayments: true,
}
