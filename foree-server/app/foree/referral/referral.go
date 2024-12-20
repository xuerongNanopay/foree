package referral

import (
	"database/sql"
	"fmt"
	"time"

	uuid_util "xue.io/go-pay/util/uuid"
)

type ReferralType string

const (
	ReferralTypeEmail ReferralType = "EMAIL"
)

const (
	sQLReferralInsert = `
		INSERT INTO referral
		(	
			referrer_id, referee_id, accept_at
		) VALUES (?,?,?)
	`
	// sQLReferralUpdateByReferralCode = `
	// 	UPDATE referral SET
	// 		referee_id = ?, accept_at = ?
	// 	WHERE referral_code = ?
	// `
	sQLReferralGetUniqueByReferralCode = `
		SELECT 
			r.id, r.referrer_id, r.referee_id, r.accept_at,
			r.created_at, r.updated_at
		FROM referral as r
		WHERE r.referral_code = ?
	`
	sQLReferralGetUniqueByRefereeId = `
		SELECT 
			r.id, r.referrer_id, r.referee_id, r.accept_at,
			r.created_at, r.updated_at
		FROM referral as r
		WHERE r.referee_id = ?
	`
)

type Referral struct {
	ID         int64      `json:"id"`
	ReferrerId int64      `json:"referrerId"`
	RefereeId  int64      `json:"refereeId"`
	AcceptAt   *time.Time `json:"acceptAt"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

func NewReferralRepo(db *sql.DB) *ReferralRepo {
	return &ReferralRepo{db: db}
}

type ReferralRepo struct {
	db *sql.DB
}

func (repo *ReferralRepo) InsertReferral(r Referral) (int64, error) {
	result, err := repo.db.Exec(
		sQLReferralInsert,
		r.ReferrerId,
		r.RefereeId,
		r.AcceptAt,
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

func (repo *ReferralRepo) GetUniqueReferralByReferralCode(referralCode string) (*Referral, error) {
	rows, err := repo.db.Query(sQLReferralGetUniqueByReferralCode, referralCode)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Referral

	for rows.Next() {
		f, err = scanRowIntoReferral(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *ReferralRepo) GetUniqueReferralByRefereeId(refereeId int64) (*Referral, error) {
	rows, err := repo.db.Query(sQLReferralGetUniqueByRefereeId, refereeId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Referral

	for rows.Next() {
		f, err = scanRowIntoReferral(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoReferral(rows *sql.Rows) (*Referral, error) {
	u := new(Referral)
	err := rows.Scan(
		&u.ID,
		&u.ReferrerId,
		&u.RefereeId,
		&u.AcceptAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GenerateReferralCode(id int64) string {
	return fmt.Sprintf("%09d-%s", id, uuid_util.UUID())
}
