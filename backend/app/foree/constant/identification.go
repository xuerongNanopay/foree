package foree_constant

import foree_auth "xue.io/go-pay/app/foree/auth"

const (
	IDTypePassport      foree_auth.IdentificationType = "PASSPORT"
	IDTypeDriverLicense foree_auth.IdentificationType = "DRIVER_LICENSE"
	IDTypeProvincalId   foree_auth.IdentificationType = "PROVINCIAL_ID"
	IDTypeNationId      foree_auth.IdentificationType = "NATIONAL_ID"
)

var AllowIdentificationTypes = map[foree_auth.IdentificationType]bool{
	foree_auth.IDTypePassport:      true,
	foree_auth.IDTypeDriverLicense: true,
	foree_auth.IDTypeProvincalId:   true,
	foree_auth.IDTypeNationId:      true,
}
