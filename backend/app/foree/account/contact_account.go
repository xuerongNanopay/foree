package account

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	sQLContactAccountInsert = `
		INSERT INTO contact_accounts
		(
			hash_id, status, type, first_name, middle_name,
			last_name, address1, address2, city, province,
			country, postal_code, phone_number, institution_name, branch_number, account_number,
			account_hash, relationship_to_contact, owner_id
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	sQLContactAccountUpdateById = `
		UPDATE contact_accounts SET 
			status = ?, latest_acitvity_at = ?
		WHERE id = ? AND a.owner_id = ?
	`
	sQLContactAccountGetUniqueById = `
		SELECT 
			a.id, a.hash_id, a.status, a.type, a.first_name, a.middle_name,
			a.last_name, a.address1, a.address2, a.city, a.province,
			a.country, a.postal_code, a.phone_number, a.institution_name, a.branch_number, a.account_number,
			a.account_hash, a.relationship_to_contact, a.owner_id
			a.latest_acitvity_at, a.create_at, a.update_at
		FROM contact_accounts a
		where a.owner_id = ? AND a.id = ? AND a.status != DELETE
	`
	sQLContactAccountGetAllByOwnerId = `
		SELECT 
			a.id, a.hash_id, a.status, a.type, a.first_name, a.middle_name,
			a.last_name, a.address1, a.address2, a.city, a.province,
			a.country, a.postal_code, a.phone_number, a.institution_name, a.branch_number, a.account_number,
			a.account_hash, a.relationship_to_contact, a.owner_id
			a.latest_acitvity_at, a.create_at, a.update_at
		FROM contact_accounts a
		where a.owner_id = ? AND a.status != DELETE
	`
	sQLContactAccountQueryByOwnerId = `
		SELECT 
			a.id, a.hash_id, a.status, a.type, a.first_name, a.middle_name,
			a.last_name, a.address1, a.address2, a.city, a.province,
			a.country, a.postal_code, a.phone_number, a.institution_name, a.branch_number, a.account_number,
			a.account_hash, a.relationship_to_contact, a.owner_id
			a.latest_acitvity_at, a.create_at, a.update_at
		FROM contact_accounts a
		where a.owner_id = ? AND a.status != DELETE
		LIMIT ? OFFSET ?
	`
)

type ContactAccountType string

const (
	ContactAccountTypeCash               ContactAccountType = "CASH"
	ContactAccountTypeAccountTransfers   ContactAccountType = "ACCOUNT_TRANSFERS"
	ContactAccountTypeThirdPartyPayments ContactAccountType = "THIRD_PARTY_PAYMENTS"
)

// TODO: improve security by using hashId.
type ContactAccount struct {
	ID                    int64              `json:"id"`
	HashId                string             `json:"hashId"`
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
	PostalCode            string             `json:"postalCode"`
	PhoneNumber           string             `json:"phoneNumber"`
	InstitutionName       string             `json:"institutionName"`
	BranchNumber          string             `json:"branchNumber"`
	AccountNumber         string             `json:"accountNumber"`
	AccountHash           string             `json:"accountHash"`
	RelationshipToContact string             `json:"relationshipToContact"`
	OwnerId               int64              `json:"owerId"`
	LatestActivityAt      time.Time          `json:"latestActivityAt"`
	CreateAt              time.Time          `json:"createAt"`
	UpdateAt              time.Time          `json:"updateAt"`
}

func (c *ContactAccount) HashMyself() {
	s := fmt.Sprintf(
		"%s%s%s%s%s%s%s%s%s%s%s",
		c.FirstName,
		c.MiddleName,
		c.LastName,
		c.Address1,
		c.Address2,
		c.City,
		c.Province,
		c.Country,
		c.InstitutionName,
		c.BranchNumber,
		c.AccountNumber,
	)
	hash := md5.Sum([]byte(s))
	c.AccountHash = hex.EncodeToString(hash[:])
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
		acc.HashId,
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
		acc.PostalCode,
		acc.PhoneNumber,
		acc.InstitutionName,
		acc.BranchNumber,
		acc.AccountNumber,
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

func (repo *ContactAccountRepo) UpdateContactAccountById(acc ContactAccount) error {
	_, err := repo.db.Exec(
		sQLContactAccountUpdateById,
		acc.Status,
		acc.LatestActivityAt,
		acc.OwnerId,
		acc.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ContactAccountRepo) GetUniqueContactAccountById(ownerid, id int64) (*ContactAccount, error) {
	rows, err := repo.db.Query(sQLContactAccountGetUniqueById, ownerid, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *ContactAccount

	for rows.Next() {
		f, err = scanRowIntoContactAccount(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *ContactAccountRepo) GetAllContactAccountByOwnerId(ownerId int64) ([]*ContactAccount, error) {
	rows, err := repo.db.Query(sQLContactAccountGetAllByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*ContactAccount, 16)
	for rows.Next() {
		p, err := scanRowIntoContactAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func scanRowIntoContactAccount(rows *sql.Rows) (*ContactAccount, error) {
	u := new(ContactAccount)
	err := rows.Scan(
		&u.ID,
		&u.HashId,
		&u.Status,
		&u.Type,
		&u.FirstName,
		&u.MiddleName,
		&u.LastName,
		&u.Address1,
		&u.Address2,
		&u.City,
		&u.Province,
		&u.Country,
		&u.PostalCode,
		&u.PhoneNumber,
		&u.InstitutionName,
		&u.BranchNumber,
		&u.AccountNumber,
		&u.AccountHash,
		&u.RelationshipToContact,
		&u.OwnerId,
		&u.LatestActivityAt,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
