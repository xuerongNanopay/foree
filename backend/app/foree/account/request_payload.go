package account

import (
	"fmt"

	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/auth"
)

type DefaultInteracReq struct {
	transport.SessionReq
	FirstName   string `json:"firstName"`
	MiddleName  string `json:"middleName"`
	LastName    string `json:"lastName"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

func generateInteracAddressFromUser(user *auth.User) string {
	if user.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", user.Address1, user.City, user.Province, user.PostalCode, user.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", user.Address1, user.Address2, user.City, user.Province, user.PostalCode, user.Country)
}
func NewDefaultInteracReqFromSession(session auth.Session) *DefaultInteracReq {
	return &DefaultInteracReq{
		SessionReq: transport.SessionReq{
			SessionId: session.ID,
		},
		FirstName:   session.User.FirstName,
		MiddleName:  session.User.MiddleName,
		LastName:    session.User.LastName,
		Address:     session.User.FirstName,
		PhoneNumber: generateInteracAddressFromUser(session.User),
		Email:       session.User.Email,
	}
}

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
