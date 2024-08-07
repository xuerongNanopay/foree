package account

import (
	"database/sql"
	"time"
)

const (
	sQLContactAccountInsert = `
		INSERT INTO contact_accounts
		(
			status, type, first_name, middle_name,
			last_name, address1, address2, city, province,
			country, phone_number, institution_name, account_number,
			account_hash, relationship_to_contact, owner_id
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	sQLContactAccountGetByOwnerId = `
		SELECT 
			a.id, a.status, a.type, a.first_name, a.middle_name,
			a.last_name, a.address1, a.address2, a.city, a.province,
			a.country, a.phone_number, a.institution_name, a.account_number,
			a.account_hash, a.relationship_to_contact, a.owner_id
			a.create_at, a.update_at
		FROM contact_accounts a
		where a.owner_id = ?
	`
	sQLContactAccountGetUniqueById = `
		SELECT 
			a.id, a.status, a.type, a.first_name, a.middle_name,
			a.last_name, a.address1, a.address2, a.city, a.province,
			a.country, a.phone_number, a.institution_name, a.account_number,
			a.account_hash, a.relationship_to_contact, a.owner_id
			a.create_at, a.update_at
		FROM contact_accounts a
		where a.owner_id = ?
	`
)

type ContactAccountType string

const (
	ContactAccountTypeCash               ContactAccountType = "CASH"
	ContactAccountTypeAccountTransfers   ContactAccountType = "ACCOUNT_TRANSFERS"
	ContactAccountTypeThirdPartyPayments ContactAccountType = "THIRD_PARTY_PAYMENTS"
)

type ContactAccount struct {
	ID                    int64              `json:"id"`
	Status                AccountStatus      `json:"status"`
	Type                  ContactAccountType `json:"type"`
	FirstName             string             `json:"firstName"`
	MiddleName            string             `json:"middleName"`
	LastName              string             `json:"lastName"`
	Address1              string             `json:"address1"`
	Address2              string             `json:"address2"`
	City                  string             `json:"city"`
	Province              string             `json:"province"`
	Country               string             `json:"country"`
	PhoneNumber           string             `json:"phoneNumber"`
	InstitutionName       string             `json:"institutionName"`
	AccountNumber         string             `json:"accountNumber"`
	AccountHash           string             `json:"accountHash"`
	RelationshipToContact string             `json:"relationshipToContact"`
	OwnerId               int64              `json:"owerId"`
	CreateAt              time.Time          `json:"createAt"`
	UpdateAt              time.Time          `json:"updateAt"`
}

func NewContactAccountRepo(db *sql.DB) *ContactAccountRepo {
	return &ContactAccountRepo{db: db}
}

type ContactAccountRepo struct {
	db *sql.DB
}

func (repo *ContactAccountRepo) InsertContactAccount(acc ContactAccount) (int64, error) {
	result, err := repo.db.Exec(
		sQLContactAccountInsert,
		acc.Status,
		acc.Type,
		acc.FirstName,
		acc.MiddleName,
		acc.LastName,
		acc.Address1,
		acc.Address2,
		acc.City,
		acc.Province,
		acc.Country,
		acc.PhoneNumber,
		acc.InstitutionName,
		acc.AccountHash,
		acc.AccountHash,
		acc.RelationshipToContact,
		acc.OwnerId,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
