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
			email, owner_id, status
		) VALUES(?,?,?,?,?,?)
	`
	sQLInteractAccountGetAll = `
	
	`
	sQLInteractAccountGetUniqueByOwnerId = `
	
	`
	sQLInteractAccountGetUniqueById = `
	
	`
)

type InteracAccount struct {
	ID         int64
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	OwnerId    int64
	Status     AccountStatus
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
}

func NewInteracAccountRepo(db *sql.DB) *InteracAccountRepo {
	return &InteracAccountRepo{db: db}
}

type InteracAccountRepo struct {
	db *sql.DB
}

func (repo *InteracAccountRepo) InsertReferral(acc InteracAccount) (int64, error) {
	result, err := repo.db.Exec(
		sQLInteracAccountInsert,
		acc.FirstName,
		acc.MiddleName,
		acc.LastName,
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
