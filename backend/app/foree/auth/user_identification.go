package auth

import (
	"database/sql"
	"time"
)

const (
	sQLUserIdentificationInsert = `
		INSERT INTO user_identifications
		(	
			status, type, value, owner_id
		) VALUES (?,?,?,?)
	`
	sQLUserIdentificationUpdateStatusById = `
		UPDATE user_identifications SET 
			status = ?
		WHERE id = ?
	`
	sQLUserIdentificationGetAllByUserId = `
		SELECT 
			u.id, u.status, u.type, u.value, u.owner_id,
			u.create_at, u.update_at
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
	ExpiredAt time.Time                `json:"expiredAt"`
	OwnerId   int64                    `json:"ownerId"`
	CreateAt  time.Time                `json:"createAt"`
	UpdateAt  time.Time                `json:"updateAt"`
}

func NewUserIdentificationRepo(db *sql.DB) *UserIdentificationRepo {
	return &UserIdentificationRepo{db: db}
}

type UserIdentificationRepo struct {
	db *sql.DB
}

func (repo *UserIdentificationRepo) InsertUserIdentification(uid UserIdentification) (int64, error) {
	result, err := repo.db.Exec(
		sQLUserIdentificationInsert,
		uid.Status,
		uid.Type,
		uid.Value,
		uid.OwnerId,
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

func (repo *UserIdentificationRepo) UpdateUserStatusById(uid UserIdentification) error {
	_, err := repo.db.Exec(
		sQLUserIdentificationUpdateStatusById,
		uid.Status,
		uid.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserIdentificationRepo) GetAllUserIdentificationByUserId(userId int64) ([]*UserIdentification, error) {
	rows, err := repo.db.Query(sQLUserIdentificationGetAllByUserId)

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
		&u.OwnerId,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
