package account

import (
	"time"
)

type ForeeContactType string

const (
	ForeeContactTypeCash               ForeeContactType = "CASH"
	ForeeContactTypeAccountTransfers   ForeeContactType = "ACCOUNT_TRANSFERS"
	ForeeContactTypeThirdPartyPayments ForeeContactType = "THIRD_PARTY_PAYMENTS"
)

type ForeeContactAccount struct {
	ID                    int64            `json:"id"`
	Status                AccountStatus    `json:"status"`
	Type                  ForeeContactType `json:"type"`
	FirstName             string           `json:"firstName"`
	MiddleName            string           `json:"middleName"`
	LastName              string           `json:"lastName"`
	Address1              string           `json:"address1"`
	Address2              string           `json:"address2"`
	City                  string           `json:"city"`
	Province              string           `json:"province"`
	Country               string           `json:"country"`
	PhoneNumber           string           `json:"phoneNumber"`
	InstitutionName       string           `json:"institutionName"`
	AccountNumber         string           `json:"accountNumber"`
	AccountHash           string           `json:"accountHash"`
	RelationshipToContact string           `json:"relationshipToContact"`
	OwnerId               int64            `json:"owerId"`
	CreateAt              time.Time        `json:"createAt"`
	UpdateAt              time.Time        `json:"updateAt"`
}
