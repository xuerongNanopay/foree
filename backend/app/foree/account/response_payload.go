package account

type ContactAccountSummaryDTO struct {
	ID              int64              `json:"id"`
	Status          AccountStatus      `json:"status"`
	FirstName       string             `json:"firstName"`
	MiddleName      string             `json:"middleName"`
	LastName        string             `json:"lastName"`
	TransferMethod  ContactAccountType `json:"transferMethod"`
	BankName        string             `json:"bankName"`
	AccountNoOrIBAN string             `json:"accountNoOrIBAN"`
}

func NewContactAccountSummaryDTO(account *ContactAccount) *ContactAccountSummaryDTO {
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
	ID                    int64              `json:"id"`
	Status                AccountStatus      `json:"status"`
	FirstName             string             `json:"firstName"`
	MiddleName            string             `json:"middleName"`
	LastName              string             `json:"lastName"`
	Address1              string             `json:"address1"`
	Address2              string             `json:"address2"`
	City                  string             `json:"city"`
	Province              string             `json:"province"`
	Country               string             `json:"country"`
	PostalCode            string             `json:"postalCode"`
	PhoneNumber           string             `json:"phoneNumber"`
	RelationshipToContact string             `json:"relationshipToContact"`
	TransferMethod        ContactAccountType `json:"transferMethod"`
	BankName              string             `json:"bankName"`
	AccountNoOrIBAN       string             `json:"accountNoOrIBAN"`
}

func NewContactAccountDetailDTO(account *ContactAccount) *ContactAccountDetailDTO {
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
	ID         int64         `json:"id"`
	Status     AccountStatus `json:"status"`
	FirstName  string        `json:"firstName"`
	MiddleName string        `json:"middleName"`
	LastName   string        `json:"lastName"`
	Email      string        `json:"email"`
}

func NewInteracAccountSummaryDTO(account *InteracAccount) *InteracAccountSummaryDTO {
	return &InteracAccountSummaryDTO{
		ID:         account.ID,
		Status:     account.Status,
		FirstName:  account.FirstName,
		MiddleName: account.MiddleName,
		LastName:   account.LastName,
		Email:      account.Email,
	}
}
