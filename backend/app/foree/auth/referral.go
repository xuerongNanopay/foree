package foree_auth

import (
	"database/sql"
	"time"
)

const (
	SQLReferralInsert = `
		INSERT INTO referrals
		(
			code, referral_type, referral_value,
			status, referrer_id, is_redeemed, expire_at
		) VALUES(?,?,?,?,?,?,?)
	`
	SQLReferralGetUniqueByCode = `
		SELECT 
			r.id, r.code, r.referral_type, r.referral_value,
			r.status, r.referrer_id, r.referree_id, r.referree_hash
			r.is_redeemed, r.expire_at, r.create_at, r.update_at
		FROM referrals r
		where r.code = ?
	`
	SQLReferralGetByReferrerId = `
		SELECT 
			r.id, r.code, r.referral_type, r.referral_value,
			r.status, r.referrer_id, r.referree_id, r.referree_hash
			r.is_redeemed, r.expire_at, r.create_at, r.update_at
		FROM referrals r
		where r.referrer_id = ?
	`
	SQLReferralUpdateReferreeByCode = `
		UPDATE referrals SET referree_id = ?, referree_hash = ?  WHERE code = ?
	`
)

type ReferralStatus string

const (
	ReferralStatusEnable  = "ENABLE"
	ReferralStatusDisable = "DISABLE"
)

type ReferralType string

const (
	ReferralTypeEmail ReferralType = "EMAIL"
	ReferralTypePhone ReferralType = "PHONE"
)

type Referral struct {
	ID            int64
	Code          string
	ReferralType  ReferralType
	ReferralValue string
	Status        ReferralStatus
	ReferrerId    int64
	ReferreeId    int64
	ReferreeHash  string
	IsRedeemed    bool
	ExpireAt      time.Time `json:"expireAt"`
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
}

func NewReferralRepo(db *sql.DB) *ReferralRepo {
	return &ReferralRepo{db: db}
}

type ReferralRepo struct {
	db *sql.DB
}

func scanRowIntoReferral(rows *sql.Rows) (*Referral, error) {
	u := new(Referral)
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
