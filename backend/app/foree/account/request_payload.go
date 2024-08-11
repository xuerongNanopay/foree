package account

import "xue.io/go-pay/app/foree/transport"

type CreateNewContactReq struct {
	transport.SessionReq
	FirstName             string             `json:"firstName" validate:"required"`
	MiddleName            string             `json:"middleName"`
	LastName              string             `json:"lastName" validate:"required"`
	Address1              string             `json:"address1" validate:"required"`
	Address2              string             `json:"address2"`
	City                  string             `json:"city" validate:"required"`
	Province              string             `json:"province" validate:"required"`
	Country               string             `json:"country" validate:"required"`
	PostalCode            string             `json:"postalCode"`
	PhoneNumber           string             `json:"phoneNumber"`
	RelationshipToContact string             `json:"relationshipToContact"`
	TransferMethod        ContactAccountType `json:"transferMethod"`
	BankName              string             `json:"bankName"`
	AccountNoOrIBAN       string             `json:"accountNoOrIBAN"`
}

type DeleteContactReq struct {
	transport.SessionReq
	ContactId int `json:"contactId" validate:"required,gte=0"`
}
