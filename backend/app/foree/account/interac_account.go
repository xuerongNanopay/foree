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
			address1, address2, city, province, country,
			phone_number, email, owner_id, status
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)
	`
	// sQLInteractAccountGetAll = `

	// `
	sQLInteractAccountGetUniqueByOwnerId = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.address1, a.address2, a.city, 
			a.province, a.country, a.phone_number,
			a.email, a.owner_id, a.status, 
			a.create_at, a.update_at
		FROM interac_accounts a
		where a.owner_id = ?
	`
	sQLInteractAccountGetUniqueById = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.address1, a.address2, a.city, 
			a.province, a.country, a.phone_number,
			a.email, a.owner_id, a.status, 
			a.create_at, a.update_at
		FROM interac_accounts a
		where a.id = ?
	`
)

type InteracAccount struct {
	ID          int64
	FirstName   string
	MiddleName  string
	LastName    string
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	City        string `json:"city"`
	Province    string `json:"province"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string
	AccountHash string `json:"accountHash"`
	OwnerId     int64
	Status      AccountStatus
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`
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
		acc.Address1,
		acc.Address2,
		acc.City,
		acc.Province,
		acc.Country,
		acc.PhoneNumber,
		acc.Email,
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
		&u.Address1,
		&u.Address2,
		&u.City,
		&u.Province,
		&u.Country,
		&u.PhoneNumber,
		&u.Email,
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
