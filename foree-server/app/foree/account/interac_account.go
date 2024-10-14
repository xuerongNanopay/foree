package account

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLInteracAccountInsert = `
		INSERT INTO interac_accounts
		(
			status, first_name, middle_name, last_name,
			email, owner_id, latest_activity_at
		) VALUES(?,?,?,?,?,?,?)
	`
	sQLInteracAccountUpdateActiveByIdAndOwner = `
		UPDATE interac_accounts SET 
			status = ?, latest_activity_at = ?
		WHERE id = ? AND a.owner_id = ? AND a.status = 'ACTIVE'
	`
	sQLInteracAccountGetUniqueById = `
		SELECT 
			a.id, a.status, a.first_name, a.middle_name, a.last_name, a.email,
			a.owner_id, a.latest_activity_at, a.created_at, a.updated_at
		FROM interac_accounts a
		where a.id = ?
	`
	sQLInteracAccountGetUniqueActiveByOwnerAndId = `
		SELECT 
			a.id, a.status, a.first_name, a.middle_name, a.last_name, a.email,
			a.owner_id, a.latest_activity_at, a.created_at, a.updated_at
		FROM interac_accounts a
		where a.owner_id = ? AND a.id = ? AND a.status = 'ACTIVE'
	`
	sQLInteracAccountGetUniqueActiveForUPdateByOwnerAndId = `
		SELECT 
			a.id, a.status, a.first_name, a.middle_name, a.last_name, a.email,
			a.owner_id, a.latest_activity_at, a.created_at, a.updated_at
		FROM interac_accounts a
		where a.owner_id = ? AND a.id = ? AND a.status = 'ACTIVE'
		FOR UPDATE
	`
	sQLInteracAccountGetAllActiveByOwnerId = `
		SELECT 
			a.id, a.status, a.first_name, a.middle_name, a.last_name, a.email,
			a.owner_id, a.latest_activity_at, a.created_at, a.updated_at
		FROM interac_accounts a
		where a.owner_id = ? AND a.status = 'ACTIVE'
		ORDER BY a.latest_activity_at DESC
	`
)

type InteracAccount struct {
	ID               int64         `json:"id"`
	Status           AccountStatus `json:"status"`
	FirstName        string        `json:"firstName"`
	MiddleName       string        `json:"middleName"`
	LastName         string        `json:"lastName"`
	Email            string        `json:"email"`
	OwnerId          int64         `json:"ownerId"`
	LatestActivityAt *time.Time    `json:"latestActivityAt"`
	CreatedAt        *time.Time    `json:"createdAt"`
	UpdatedAt        *time.Time    `json:"updatedAt"`
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
			acc.Email,
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
			acc.Email,
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

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *InteracAccountRepo) GetUniqueActiveInteracAccountForUpdateByOwnerAndId(ctx context.Context, ownerId, id int64) (*InteracAccount, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var rows *sql.Rows

	if ok {
		rows, err = dTx.Query(sQLInteracAccountGetUniqueActiveForUPdateByOwnerAndId, ownerId, id)

	} else {
		rows, err = repo.db.Query(sQLInteracAccountGetUniqueActiveForUPdateByOwnerAndId, ownerId, id)
	}

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

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *InteracAccountRepo) GetUniqueInteracAccountById(ctx context.Context, id int64) (*InteracAccount, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var rows *sql.Rows
	var err error

	if ok {
		fmt.Println("aaaaa")
		rows, err = dTx.Query(sQLInteracAccountGetUniqueById, id)
	} else {
		rows, err = repo.db.Query(sQLInteracAccountGetUniqueById, id)
	}

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

	if f == nil || f.ID == 0 {
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

	accounts := make([]*InteracAccount, 0)
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
		&u.Status,
		&u.FirstName,
		&u.MiddleName,
		&u.LastName,
		&u.Email,
		&u.OwnerId,
		&u.LatestActivityAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
