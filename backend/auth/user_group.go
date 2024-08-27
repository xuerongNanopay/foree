package auth

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLUserGroupInsert = `
		INSERT INTO user_group(
			role_group, transaction_limit_group, owner_id
		) VALUES(?,?,?)
	`
	sQLUserGoupUpdate = `
		UPDATE user_group SET
			role_group = ?, transaction_limit_group = ?
		WHERE owner_id = ?
	`
	sQLUserGroupGetUniqueByOwnerId = `
		SELECT
			u.id, u.role_group, u.transaction_limit_group,
			u.owner_id, u.create_at, u.update_at
		FROM user_group as u
		WHERE u.owner_id = ?
	`
)

type UserGroup struct {
	ID                    int64     `json:"id"`
	RoleGroup             string    `json:"roleGroup"`
	TransactionLimitGroup string    `json:"transactionLimitGroup"`
	OwnerId               int64     `json:"ownerId"`
	CreateAt              time.Time `json:"createAt"`
	UpdateAt              time.Time `json:"updateAt"`
}

func NewUserGroupRepo(db *sql.DB) *UserGroupRepo {
	return &UserGroupRepo{db: db}
}

type UserGroupRepo struct {
	db *sql.DB
}

func (repo *UserGroupRepo) InsertUserGroup(ctx context.Context, ug UserGroup) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLUserGroupInsert,
			ug.RoleGroup,
			ug.TransactionLimitGroup,
			ug.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLUserGroupInsert,
			ug.RoleGroup,
			ug.TransactionLimitGroup,
			ug.OwnerId,
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

func (repo *UserGroupRepo) UpdateUserById(ctx context.Context, ug UserGroup) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(
			sQLUserGoupUpdate,
			ug.RoleGroup,
			ug.TransactionLimitGroup,
			ug.OwnerId,
		)
	} else {
		_, err = repo.db.Exec(
			sQLUserGoupUpdate,
			ug.RoleGroup,
			ug.TransactionLimitGroup,
			ug.OwnerId,
		)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *UserGroupRepo) GetUniqueUserGroupByOwnerId(ownerId int64) (*UserGroup, error) {
	rows, err := repo.db.Query(sQLUserGroupGetUniqueByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u *UserGroup

	for rows.Next() {
		u, err = scanRowIntoUserGroup(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func scanRowIntoUserGroup(rows *sql.Rows) (*UserGroup, error) {
	u := new(UserGroup)
	err := rows.Scan(
		&u.ID,
		&u.RoleGroup,
		&u.TransactionLimitGroup,
		&u.OwnerId,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
