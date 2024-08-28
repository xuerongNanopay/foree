package foree_auth

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLUserIdentificationInsert = `
		INSERT INTO user_identifications
		(	
			status, type, value, link1, link2, owner_id
		) VALUES (?,?,?,?,?,?)
	`
	sQLUserIdentificationUpdateStatusByOwnerId = `
		UPDATE user_identifications SET 
			status = ?, link1 = ? , link2 = ?
		WHERE owner_id = ?
	`
	sQLUserIdentificationGetAllByOwnerId = `
		SELECT 
			u.id, u.status, u.type, u.value, u.link1, u.link2,
			u.owner_id, u.created_at, u.updated_at
		FROM user_identifications as u 
		WHERE u.owner_id = ?
	`
)

type UserIdentificationStatus string

const (
	IdentificationStatusAwaitApprove UserIdentificationStatus = "AWAIT_APPROVE"
	IdentificationStatusActive       UserIdentificationStatus = "ACTIVE"
	IdentificationStatusExpired      UserIdentificationStatus = "EXPIRED"
	IdentificationStatusDisabled     UserIdentificationStatus = "Disabled"
)

type IdentificationType string

const (
	IDTypePassport      IdentificationType = "PASSPORT"
	IDTypeDriverLicense IdentificationType = "DRIVER_LICENSE"
	IDTypeProvincalId   IdentificationType = "PROVINCIAL_ID"
	IDTypeNationId      IdentificationType = "NATIONAL_ID"
)

type UserIdentification struct {
	ID        int64                    `json:"id"`
	Status    UserIdentificationStatus `json:"status"`
	Type      IdentificationType       `json:"type"`
	Value     string                   `json:"value"`
	Link1     string                   `json:"link1"`
	Link2     string                   `json:"link2"`
	ExpiredAt time.Time                `json:"expiredAt"`
	OwnerId   int64                    `json:"ownerId"`
	CreatedAt time.Time                `json:"createdAt"`
	UpdateAt  time.Time                `json:"updatedAt"`
}

func NewUserIdentificationRepo(db *sql.DB) *UserIdentificationRepo {
	return &UserIdentificationRepo{db: db}
}

type UserIdentificationRepo struct {
	db *sql.DB
}

func (repo *UserIdentificationRepo) InsertUserIdentification(ctx context.Context, uid UserIdentification) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLUserIdentificationInsert,
			uid.Status,
			uid.Type,
			uid.Value,
			uid.Link1,
			uid.Link2,
			uid.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLUserIdentificationInsert,
			uid.Status,
			uid.Type,
			uid.Value,
			uid.Link1,
			uid.Link2,
			uid.OwnerId,
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

func (repo *UserIdentificationRepo) UpdateUserStatusByOwnerId(uid UserIdentification) error {
	_, err := repo.db.Exec(
		sQLUserIdentificationUpdateStatusByOwnerId,
		uid.Status,
		uid.Link1,
		uid.Link2,
		uid.OwnerId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserIdentificationRepo) GetAllUserIdentificationByOwnerId(ownerId int64) ([]*UserIdentification, error) {
	rows, err := repo.db.Query(sQLUserIdentificationGetAllByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uids []*UserIdentification

	for rows.Next() {
		u, err := scanRowIntoUserIdentification(rows)
		if err != nil {
			return nil, err
		}
		uids = append(uids, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uids, nil
}

func scanRowIntoUserIdentification(rows *sql.Rows) (*UserIdentification, error) {
	u := new(UserIdentification)
	err := rows.Scan(
		&u.ID,
		&u.Status,
		&u.Type,
		&u.Value,
		&u.Link1,
		&u.Link2,
		&u.OwnerId,
		&u.CreatedAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
