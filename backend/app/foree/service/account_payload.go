package service

import (
	"fmt"
	"strings"

	"xue.io/go-pay/app/foree/account"
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/auth"
)

func generateInteracAddressFromUser(user *auth.User) string {
	if user.Address2 == "" {
		return fmt.Sprintf("%s,%s,%s,%s,%s", user.Address1, user.City, user.Province, user.PostalCode, user.Country)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", user.Address1, user.Address2, user.City, user.Province, user.PostalCode, user.Country)
}

type CreateContactReq struct {
	transport.SessionReq
	FirstName             string `json:"firstName" validate:"required"`
	MiddleName            string `json:"middleName"`
	LastName              string `json:"lastName" validate:"required"`
	Address1              string `json:"address1" validate:"required"`
	Address2              string `json:"address2"`
	City                  string `json:"city" validate:"required"`
	Province              string `json:"province" validate:"required"`
	Country               string `json:"country" validate:"required"`
	PostalCode            string `json:"postalCode"`
	PhoneNumber           string `json:"phoneNumber"`
	RelationshipToContact string `json:"relationshipToContact"`
	TransferMethod        string `json:"transferMethod"`
	BankName              string `json:"bankName"`
	AccountNoOrIBAN       string `json:"accountNoOrIBAN"`
}

func (q *CreateContactReq) TrimSpace() {
	q.FirstName = strings.TrimSpace(q.FirstName)
	q.MiddleName = strings.TrimSpace(q.MiddleName)
	q.LastName = strings.TrimSpace(q.LastName)
	q.Address1 = strings.TrimSpace(q.Address1)
	q.Address2 = strings.TrimSpace(q.Address2)
	q.City = strings.TrimSpace(q.City)
	q.Province = strings.TrimSpace(q.Province)
	q.Country = strings.TrimSpace(q.Country)
	q.PostalCode = strings.TrimSpace(q.PostalCode)
	q.PhoneNumber = strings.TrimSpace(q.PhoneNumber)
	q.RelationshipToContact = strings.TrimSpace(q.RelationshipToContact)
	q.TransferMethod = strings.TrimSpace(q.TransferMethod)
	q.BankName = strings.TrimSpace(q.BankName)
	q.AccountNoOrIBAN = strings.TrimSpace(q.AccountNoOrIBAN)
}

func (q *CreateContactReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid create contact request")

	// Check relationship
	_, ok := allowRelationshipToContactTypes[q.RelationshipToContact]
	if !ok {
		ret.AddDetails("relationshipToContact", fmt.Sprintf("invalid relationshipToContact `%v`", q.RelationshipToContact))
	}

	// Check transferMethod
	_, ok = allowContactAccountType[account.ContactAccountType(q.TransferMethod)]
	if !ok {
		ret.AddDetails("transferMethod", fmt.Sprintf("invalid transferMethod `%v`", q.TransferMethod))
	}

	if account.ContactAccountType(q.TransferMethod) != account.ContactAccountTypeCash {
		if q.BankName == "" {
			ret.AddDetails("bankName", fmt.Sprintf("invalid bankName `%v`", q.BankName))
		}
		if q.AccountNoOrIBAN == "" {
			ret.AddDetails("accountNoOrIBAN", fmt.Sprintf("invalid accountNoOrIBAN `%v`", q.AccountNoOrIBAN))
		}
	}

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type DeleteContactReq struct {
	transport.SessionReq
	ContactId int64 `json:"contactId" validate:"required,gte=0"`
}

func (q *DeleteContactReq) TrimSpace() {
}

func (q *DeleteContactReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid delete contact request")

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type GetContactReq struct {
	transport.SessionReq
	ContactId int64 `json:"contactId" validate:"required,gte=0"`
}

func (q *GetContactReq) TrimSpace() {
}

func (q *GetContactReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid get contact request")

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type QueryContactReq struct {
	transport.SessionReq
	Offset int `json:"offset" validate:"required,gte=0"`
	Limit  int `json:"limit" validate:"required,gte=0"`
}

func (q *QueryContactReq) TrimSpace() {
}

func (q *QueryContactReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid query contact request")

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
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
