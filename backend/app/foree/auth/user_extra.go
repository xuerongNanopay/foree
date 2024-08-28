package foree_auth

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLUserExtraInsert = `
		INSERT INTO user_extra
		(
			pob, cor, nationality, occupation_category, 
			occupation_name, ower_id
		) VALUES(?,?,?,?,?,?)
	`
	sQLUserExtraUpdate = `
		UPDATE user_extra SET
			pob = ?, cor = ? nationality = ?,
			occupation_category = ?,  occupation_name = ?
		WHERE owner_id = ?
	`
	sQLUserExtraGetUniqueByOwnerId = `
		SELECT
			u.id, u.pob, u.cor, u.nationality, u.occupation_category,
			u.occupation_name, u.owner_id, u.created_at, u.updated_at
		FROM user_extra as u
		WHERE u.owner_id = ?
	`
)

type UserExtra struct {
	ID                 int64     `json:"id"`
	Pob                string    `json:"pob"`
	Cor                string    `json:"cor"`
	Nationality        string    `json:"nationality"`
	OccupationCategory string    `json:"occupationCategory"`
	OccupationName     string    `json:"occupationName"`
	OwnerId            int64     `json:"ownerId"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdateAt           time.Time `json:"updatedAt"`
}

func NewUserExtraRepo(db *sql.DB) *UserExtraRepo {
	return &UserExtraRepo{db: db}
}

type UserExtraRepo struct {
	db *sql.DB
}

func (repo *UserExtraRepo) InsertUserExtra(ctx context.Context, ue UserExtra) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLUserExtraInsert,
			ue.Pob,
			ue.Cor,
			ue.Nationality,
			ue.OccupationCategory,
			ue.OccupationName,
			ue.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLUserExtraInsert,
			ue.Pob,
			ue.Cor,
			ue.Nationality,
			ue.OccupationCategory,
			ue.OccupationName,
			ue.OwnerId,
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

func (repo *UserExtraRepo) GetUniqueUserExtraByOwnerId(ownerId int64) (*UserExtra, error) {
	rows, err := repo.db.Query(sQLUserExtraGetUniqueByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var u *UserExtra

	for rows.Next() {
		u, err = scanRowIntoUserExtra(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, nil
	}

	return u, nil
}

func (repo *UserExtraRepo) UpdateUserExtrayOwnerId(ctx context.Context, ue UserExtra) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(
			sQLUserExtraUpdate,
			ue.Pob,
			ue.Cor,
			ue.Nationality,
			ue.OccupationCategory,
			ue.OccupationName,
			ue.OwnerId,
		)
	} else {
		_, err = repo.db.Exec(
			sQLUserExtraUpdate,
			ue.Pob,
			ue.Cor,
			ue.Nationality,
			ue.OccupationCategory,
			ue.OccupationName,
			ue.OwnerId,
		)
	}

	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoUserExtra(rows *sql.Rows) (*UserExtra, error) {
	u := new(UserExtra)
	err := rows.Scan(
		&u.ID,
		&u.Pob,
		&u.Cor,
		&u.Nationality,
		&u.OccupationCategory,
		&u.OccupationName,
		&u.OwnerId,
		&u.CreatedAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
