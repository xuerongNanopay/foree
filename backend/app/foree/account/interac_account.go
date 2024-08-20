package account

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/constant"
)

const (
	sQLInteracAccountInsert = `
		INSERT INTO interac_accounts
		(
			status, first_name, middle_name, last_name,
			address, phone_number, email, 
			institution_name, branch_number, account_number,
			owner_id, latest_acitvity_at
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)
	`
	sQLInteracAccountUpdateActiveByIdAndOwner = `
		UPDATE interac_accounts SET 
			status = ?, latest_acitvity_at = ?
		WHERE id = ? AND a.owner_id = ? AND a.status = ACTIVE
	`
	sQLInteracAccountGetUniqueById = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.address, a.phone_number, a.email, 
			a.institution_name, a.branch_number, a.account_number,
			a.owner_id, a.status, 
			a.latest_acitvity_at, a.create_at, a.update_at
		FROM interac_accounts a
		where a.id = ?
	`
	sQLInteracAccountGetUniqueActiveByOwnerAndId = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.address, a.phone_number, a.email, 
			a.institution_name, a.branch_number, a.account_number,
			a.owner_id, a.status, 
			a.latest_acitvity_at, a.create_at, a.update_at
		FROM interac_accounts a
		where a.owner_id = ? AND a.id = ? AND a.status = ACTIVE
	`
	sQLInteracAccountGetAllActiveByOwnerId = `
		SELECT 
			a.id, a.first_name, a.middle_name,
			a.last_name, a.address, a.phone_number, a.email, 
			a.institution_name, a.branch_number, a.account_number,
			a.owner_id, a.status, 
			a.create_at, a.update_at
		FROM interac_accounts a
		where a.owner_id = ? AND a.status = ACTIVE
		ORDER BY a.latest_acitvity_at DESC
	`
)

type InteracAccount struct {
	ID               int64         `json:"id"`
	FirstName        string        `json:"firstName"`
	MiddleName       string        `json:"middleName"`
	LastName         string        `json:"lastName"`
	Address          string        `json:"address"`
	PhoneNumber      string        `json:"phoneNumber"`
	Email            string        `json:"email"`
	InstitutionName  string        `json:"institutionName"`
	BranchNumber     string        `json:"branchNumber"`
	AccountNumber    string        `json:"accountNumber"`
	AccountHash      string        `json:"accountHash"`
	OwnerId          int64         `json:"ownerId"`
	Status           AccountStatus `json:"status"`
	LatestActivityAt time.Time     `json:"latestActivityAt"`
	CreateAt         time.Time     `json:"createAt"`
	UpdateAt         time.Time     `json:"updateAt"`
}

func (c *InteracAccount) GetLegalName() string {
	if c.MiddleName == "" {
		return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
	}
	return fmt.Sprintf("%s %s %s", c.FirstName, c.MiddleName, c.LastName)
}

func NewInteracAccountRepo(db *sql.DB) *InteracAccountRepo {
	return &InteracAccountRepo{db: db}
}

type InteracAccountRepo struct {
	db *sql.DB
}

func (repo *InteracAccountRepo) InsertInteracAccount(ctx context.Context, acc InteracAccount) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLInteracAccountInsert,
			acc.Status,
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
			acc.LatestActivityAt,
		)
	} else {
		result, err = repo.db.Exec(
			sQLInteracAccountInsert,
			acc.Status,
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
			acc.LatestActivityAt,
		)
	}

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *InteracAccountRepo) UpdateActiveInteracAccountByIdAndOwner(ctx context.Context, acc InteracAccount) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(
			sQLInteracAccountUpdateActiveByIdAndOwner,
			acc.Status,
			acc.LatestActivityAt,
			acc.OwnerId,
			acc.ID,
		)
	} else {
		_, err = repo.db.Exec(
			sQLInteracAccountUpdateActiveByIdAndOwner,
			acc.Status,
			acc.LatestActivityAt,
			acc.OwnerId,
			acc.ID,
		)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *InteracAccountRepo) GetUniqueActiveInteracAccountByOwnerAndId(ctx context.Context, ownerId, id int64) (*InteracAccount, error) {
	rows, err := repo.db.Query(sQLInteracAccountGetUniqueActiveByOwnerAndId, ownerId, id)

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

func (repo *InteracAccountRepo) GetUniqueInteracAccountById(ctx context.Context, id int64) (*InteracAccount, error) {
	rows, err := repo.db.Query(sQLInteracAccountGetUniqueById, id)

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

func (repo *InteracAccountRepo) GetAllActiveInteracAccountByOwnerId(ctx context.Context, ownerId int64) ([]*InteracAccount, error) {
	rows, err := repo.db.Query(sQLInteracAccountGetAllActiveByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := make([]*InteracAccount, 16)
	for rows.Next() {
		p, err := scanRowIntoInteracAccount(rows)
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
		&u.LatestActivityAt,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
