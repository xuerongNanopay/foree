package service

import (
	"fmt"

	"xue.io/go-pay/app/foree/account"
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
	OwnerId     int64  `json:"ownerId"`
}

func generateInteracAddressFromUser(user *auth.User) string {
	if user.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", user.Address1, user.City, user.Province, user.PostalCode, user.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", user.Address1, user.Address2, user.City, user.Province, user.PostalCode, user.Country)
}
func NewDefaultInteracReqFromSession(session *auth.Session) *DefaultInteracReq {
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
		OwnerId:     session.User.ID,
	}
}

type CreateNewContactReq struct {
	transport.SessionReq
	FirstName             string                     `json:"firstName" validate:"required"`
	MiddleName            string                     `json:"middleName"`
	LastName              string                     `json:"lastName" validate:"required"`
	Address1              string                     `json:"address1" validate:"required"`
	Address2              string                     `json:"address2"`
	City                  string                     `json:"city" validate:"required"`
	Province              string                     `json:"province" validate:"required"`
	Country               string                     `json:"country" validate:"required"`
	PostalCode            string                     `json:"postalCode"`
	PhoneNumber           string                     `json:"phoneNumber"`
	RelationshipToContact string                     `json:"relationshipToContact"`
	TransferMethod        account.ContactAccountType `json:"transferMethod"`
	BankName              string                     `json:"bankName"`
	AccountNoOrIBAN       string                     `json:"accountNoOrIBAN"`
}

type DeleteContactReq struct {
	transport.SessionReq
	ContactId int `json:"contactId" validate:"required,gte=0"`
}

// ----------   Response --------------

type ContactAccountSummaryDTO struct {
	ID              int64                      `json:"id"`
	Status          account.AccountStatus      `json:"status"`
	FirstName       string                     `json:"firstName"`
	MiddleName      string                     `json:"middleName"`
	LastName        string                     `json:"lastName"`
	TransferMethod  account.ContactAccountType `json:"transferMethod"`
	BankName        string                     `json:"bankName"`
	AccountNoOrIBAN string                     `json:"accountNoOrIBAN"`
}

func NewContactAccountSummaryDTO(account *account.ContactAccount) *ContactAccountSummaryDTO {
	return &ContactAccountSummaryDTO{
		ID:              account.ID,
		Status:          account.Status,
		FirstName:       account.FirstName,
		MiddleName:      account.MiddleName,
		LastName:        account.LastName,
		TransferMethod:  account.Type,
		BankName:        account.InstitutionName,
		AccountNoOrIBAN: account.BranchNumber,
	}
}

type ContactAccountDetailDTO struct {
	ID                    int64                      `json:"id"`
	Status                account.AccountStatus      `json:"status"`
	FirstName             string                     `json:"firstName"`
	MiddleName            string                     `json:"middleName"`
	LastName              string                     `json:"lastName"`
	Address1              string                     `json:"address1"`
	Address2              string                     `json:"address2"`
	City                  string                     `json:"city"`
	Province              string                     `json:"province"`
	Country               string                     `json:"country"`
	PostalCode            string                     `json:"postalCode"`
	PhoneNumber           string                     `json:"phoneNumber"`
	RelationshipToContact string                     `json:"relationshipToContact"`
	TransferMethod        account.ContactAccountType `json:"transferMethod"`
	BankName              string                     `json:"bankName"`
	AccountNoOrIBAN       string                     `json:"accountNoOrIBAN"`
}

func NewContactAccountDetailDTO(account *account.ContactAccount) *ContactAccountDetailDTO {
	return &ContactAccountDetailDTO{
		ID:                    account.ID,
		Status:                account.Status,
		FirstName:             account.FirstName,
		MiddleName:            account.MiddleName,
		LastName:              account.LastName,
		Address1:              account.Address1,
		Address2:              account.Address2,
		City:                  account.City,
		Province:              account.Province,
		Country:               account.Country,
		PostalCode:            account.PostalCode,
		PhoneNumber:           account.PhoneNumber,
		RelationshipToContact: account.RelationshipToContact,
		TransferMethod:        account.Type,
		BankName:              account.InstitutionName,
		AccountNoOrIBAN:       account.BranchNumber,
	}
}

type InteracAccountSummaryDTO struct {
	ID         int64                 `json:"id"`
	Status     account.AccountStatus `json:"status"`
	FirstName  string                `json:"firstName"`
	MiddleName string                `json:"middleName"`
	LastName   string                `json:"lastName"`
	Email      string                `json:"email"`
}

func NewInteracAccountSummaryDTO(account *account.InteracAccount) *InteracAccountSummaryDTO {
	return &InteracAccountSummaryDTO{
		ID:         account.ID,
		Status:     account.Status,
		FirstName:  account.FirstName,
		MiddleName: account.MiddleName,
		LastName:   account.LastName,
		Email:      account.Email,
	}
}
