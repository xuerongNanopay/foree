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
	// sQLInteractAccountGetAll = `

	// `
	sQLInteractAccountGetUniqueByOwnerId = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.email, a.owner_id,
			a.status, a.create_at, a.update_at
		FROM interac_accounts a
		where a.owner_id = ?
	`
	sQLInteractAccountGetUniqueById = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.email, a.owner_id,
			a.status, a.create_at, a.update_at
		FROM interac_accounts a
		where a.id = ?
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

func (repo *InteracAccountRepo) InsertInteracAccount(acc InteracAccount) (int64, error) {
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

func (repo *InteracAccountRepo) GetUniqueInteractAccountByOwnerId(ownerId int64) {

}

func scanRowIntoInteracAccount(rows *sql.Rows) (*InteracAccount, error) {
	u := new(InteracAccount)
	err := rows.Scan(
		&u.ID,
		&u.Code,
		&u.ReferralType,
		&u.ReferralValue,
		&u.Status,
		&u.ReferrerId,
		&u.ReferreeId,
		&u.IsRedeemed,
		&u.ExpireAt,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
