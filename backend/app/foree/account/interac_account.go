package account

import (
	"database/sql"
	"time"
)

const (
	sQLInteracAccountInsert = `
		INSERT INTO interac_accounts
		(
			first_name, middle_name, last_name,
			address, phone_number, email, 
			institution_name, branch_number, account_number,
			owner_id, status
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)
	`
	// sQLInteractAccountGetAll = `

	// `
	sQLInteractAccountGetUniqueByOwnerId = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.address, a.phone_number, a.email, 
			a.institution_name, a.branch_number, a.account_number,
			a.owner_id, a.status, 
			a.create_at, a.update_at
		FROM interac_accounts a
		where a.owner_id = ? AND a.status != DELETE
	`
	sQLInteractAccountGetUniqueById = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.address, a.phone_number, a.email, 
			a.institution_name, a.branch_number, a.account_number,
			a.owner_id, a.status, 
			a.create_at, a.update_at
		FROM interac_accounts a
		where a.id = ? AND a.status != DELETE
	`
)

type InteracAccount struct {
	ID              int64         `json:"id"`
	FirstName       string        `json:"firstName"`
	MiddleName      string        `json:"middleName"`
	LastName        string        `json:"lastName"`
	Address         string        `json:"address"`
	PhoneNumber     string        `json:"phoneNumber"`
	Email           string        `json:"email"`
	InstitutionName string        `json:"institutionName"`
	BranchNumber    string        `json:"branchNumber"`
	AccountNumber   string        `json:"accountNumber"`
	AccountHash     string        `json:"accountHash"`
	OwnerId         int64         `json:"ownerId"`
	Status          AccountStatus `json:"status"`
	CreateAt        time.Time     `json:"createAt"`
	UpdateAt        time.Time     `json:"updateAt"`
}

func NewInteracAccountRepo(db *sql.DB) *InteracAccountRepo {
	return &InteracAccountRepo{db: db}
}

type InteracAccountRepo struct {
	db *sql.DB
}

func (repo *InteracAccountRepo) InsertInteracAccount(acc InteracAccount) (int64, error) {
	result, err := repo.db.Exec(
		sQLInteracAccountInsert,
		acc.FirstName,
		acc.MiddleName,
		acc.LastName,
		acc.Address,
		acc.PhoneNumber,
		acc.Email,
		acc.InstitutionName,
		acc.BranchNumber,
		acc.AccountNumber,
		acc.OwnerId,
		acc.Status,
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

func (repo *InteracAccountRepo) GetUniqueInteractAccountByOwnerId(ownerId int64) (*InteracAccount, error) {
	rows, err := repo.db.Query(sQLInteractAccountGetUniqueByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *InteracAccount

	for rows.Next() {
		f, err = scanRowIntoInteracAccount(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *InteracAccountRepo) GetUniqueInteractAccountById(id int64) (*InteracAccount, error) {
	rows, err := repo.db.Query(sQLInteractAccountGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *InteracAccount

	for rows.Next() {
		f, err = scanRowIntoInteracAccount(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoInteracAccount(rows *sql.Rows) (*InteracAccount, error) {
	u := new(InteracAccount)
	err := rows.Scan(
		&u.ID,
		&u.FirstName,
		&u.MiddleName,
		&u.LastName,
		&u.Address,
		&u.PhoneNumber,
		&u.Email,
		&u.InstitutionName,
		&u.BranchNumber,
		&u.AccountNumber,
		&u.OwnerId,
		&u.Status,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
