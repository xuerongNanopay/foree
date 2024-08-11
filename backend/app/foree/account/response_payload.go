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
