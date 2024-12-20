package foree_account_service

import (
	"fmt"

	"xue.io/go-pay/app/foree/account"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/server/transport"
	"xue.io/go-pay/validation"
)

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
	RelationshipToContact string `json:"relationshipToContact" validate:"required"`
	TransferMethod        string `json:"transferMethod" validate:"required"`
	BankName              string `json:"bankName"`
	AccountNoOrIBAN       string `json:"accountNoOrIBAN"`
}

func (q CreateContactReq) Validate() *transport.BadRequestError {
	ret := validation.ValidateStruct(q, "Invalid create contact request")

	// Check relationship
	_, ok := foree_constant.AllowRelationshipToContactTypes[q.RelationshipToContact]
	if !ok {
		ret.AddDetails("relationshipToContact", fmt.Sprintf("invalid relationshipToContact `%v`", q.RelationshipToContact))
	}

	// Check transferMethod
	_, ok = foree_constant.AllowContactAccountType[account.ContactAccountType(q.TransferMethod)]
	if !ok {
		ret.AddDetails("transferMethod", fmt.Sprintf("invalid transferMethod `%v`", q.TransferMethod))
	}

	if account.ContactAccountType(q.TransferMethod) != foree_constant.ContactAccountTypeCash {
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

func (q DeleteContactReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid delete contact request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type GetContactReq struct {
	transport.SessionReq
	ContactId int64 `json:"contactId" validate:"required,gt=0"`
}

func (q GetContactReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid get contact request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type QueryContactReq struct {
	transport.SessionReq
	Offset int `json:"offset" validate:"required,gte=0"`
	Limit  int `json:"limit" validate:"required,gt=0"`
}

func (q QueryContactReq) Validate() *transport.BadRequestError {
	if ret := validation.ValidateStruct(q, "Invalid query contact request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

// ----------   Response --------------

type ContactAccountSummaryDTO struct {
	ID              int64                      `json:"id,omitempty"`
	Status          account.AccountStatus      `json:"status,omitempty"`
	FirstName       string                     `json:"firstName,omitempty"`
	MiddleName      string                     `json:"middleName,omitempty"`
	LastName        string                     `json:"lastName,omitempty"`
	TransferMethod  account.ContactAccountType `json:"transferMethod,omitempty"`
	BankName        string                     `json:"bankName,omitempty"`
	AccountNoOrIBAN string                     `json:"accountNoOrIBAN,omitempty"`
	LatestActiveAt  int64                      `json:"latestActiveAt,omitempty"`
	CreateAt        int64                      `json:"createAt,omitempty"`
}

func NewContactAccountSummaryDTO(account *account.ContactAccount) *ContactAccountSummaryDTO {
	ret := &ContactAccountSummaryDTO{
		ID:              account.ID,
		Status:          account.Status,
		FirstName:       account.FirstName,
		MiddleName:      account.MiddleName,
		LastName:        account.LastName,
		TransferMethod:  account.Type,
		BankName:        account.InstitutionName,
		AccountNoOrIBAN: account.AccountNumber,
	}
	if account.CreatedAt != nil {
		ret.CreateAt = account.CreatedAt.UnixMilli()
	}
	if account.LatestActivityAt != nil {
		ret.LatestActiveAt = account.LatestActivityAt.UnixMilli()
	}
	return ret
}

type ContactAccountDetailDTO struct {
	ID                    int64                      `json:"id,omitempty"`
	Status                account.AccountStatus      `json:"status,omitempty"`
	FirstName             string                     `json:"firstName,omitempty"`
	MiddleName            string                     `json:"middleName,omitempty"`
	LastName              string                     `json:"lastName,omitempty"`
	Address1              string                     `json:"address1,omitempty"`
	Address2              string                     `json:"address2,omitempty"`
	City                  string                     `json:"city,omitempty"`
	Province              string                     `json:"province,omitempty"`
	Country               string                     `json:"country,omitempty"`
	PostalCode            string                     `json:"postalCode,omitempty"`
	PhoneNumber           string                     `json:"phoneNumber,omitempty"`
	RelationshipToContact string                     `json:"relationshipToContact,omitempty"`
	TransferMethod        account.ContactAccountType `json:"transferMethod,omitempty"`
	BankName              string                     `json:"bankName,omitempty"`
	AccountNoOrIBAN       string                     `json:"accountNoOrIBAN,omitempty"`
	LatestActiveAt        int64                      `json:"latestActiveAt,omitempty"`
	CreateAt              int64                      `json:"createAt,omitempty"`
}

func NewContactAccountDetailDTO(account *account.ContactAccount) *ContactAccountDetailDTO {
	ret := &ContactAccountDetailDTO{
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
		AccountNoOrIBAN:       account.AccountNumber,
	}
	if account.CreatedAt != nil {
		ret.CreateAt = account.CreatedAt.UnixMilli()
	}
	if account.LatestActivityAt != nil {
		ret.LatestActiveAt = account.LatestActivityAt.UnixMilli()
	}
	return ret
}

type InteracAccountSummaryDTO struct {
	ID         int64                 `json:"id,omitempty"`
	Status     account.AccountStatus `json:"status,omitempty"`
	FirstName  string                `json:"firstName,omitempty"`
	MiddleName string                `json:"middleName,omitempty"`
	LastName   string                `json:"lastName,omitempty"`
	Email      string                `json:"email,omitempty"`
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

type VerifyContactDTO struct {
	AccountNo     string `json:"accountNo,omitempty"`
	Iban          string `json:"iban,omitempty"`
	BankName      string `json:"bankName,omitempty"`
	BranchCode    string `json:"branchCode,omitempty"`
	AccountTitle  string `json:"accountTitle,omitempty"`
	AccountStatus string `json:"accountStatus,omitempty"`
}
