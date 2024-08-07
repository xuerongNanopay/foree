package foree_auth

import (
	"database/sql"
	"time"
)

const (
	sQLReferralInsert = `
		INSERT INTO referrals
		(
			code, referral_type, referral_value,
			status, referrer_id, is_redeemed, expire_at
		) VALUES(?,?,?,?,?,?,?)
	`
	sQLReferralGetUniqueByCode = `
		SELECT 
			r.id, r.code, r.referral_type, r.referral_value,
			r.status, r.referrer_id, r.referree_id, r.referree_hash
			r.is_redeemed, r.expire_at, r.create_at, r.update_at
		FROM referrals r
		where r.code = ?
	`
	sQLReferralGetAllByReferrerId = `
		SELECT 
			r.id, r.code, r.referral_type, r.referral_value,
			r.status, r.referrer_id, r.referree_id, r.referree_hash
			r.is_redeemed, r.expire_at, r.create_at, r.update_at
		FROM referrals r
		where r.referrer_id = ?
	`
	sQLReferralGetUniqueByReferreeHash = `
		SELECT 
			r.id, r.code, r.referral_type, r.referral_value,
			r.status, r.referrer_id, r.referree_id, r.referree_hash
			r.is_redeemed, r.expire_at, r.create_at, r.update_at
		FROM referrals r
		where r.referree_hash = ?
	`
	sQLReferralUpdateReferralByCode = `
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

func (repo *ReferralRepo) InsertReferral(referal Referral) (int64, error) {
	result, err := repo.db.Exec(
		sQLReferralInsert,
		referal.Code,
		referal.ReferralType,
		referal.ReferralValue,
		referal.Status,
		referal.ReferrerId,
		referal.IsRedeemed,
		referal.ExpireAt,
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

func (repo *ReferralRepo) UpdateReferralByCode(refereeId string, referreeHash string, code string) error {
	_, err := repo.db.Exec(sQLReferralUpdateReferralByCode, refereeId, referreeHash, code)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ReferralRepo) GetAllReferralByReferrerId(referrerId int64) ([]*Referral, error) {
	rows, err := repo.db.Query(sQLReferralGetAllByReferrerId, referrerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	referrals := make([]*Referral, 16)
	for rows.Next() {
		p, err := scanRowIntoReferral(rows)
		if err != nil {
			return nil, err
		}
		referrals = append(referrals, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return referrals, nil
}

func (repo *ReferralRepo) GetUniqueReferralByCode(code string) (*Referral, error) {
	rows, err := repo.db.Query(sQLReferralGetUniqueByCode, code)

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

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *ReferralRepo) GetUniqueReferralByReferreeHash(referreeHash string) (*Referral, error) {
	rows, err := repo.db.Query(sQLReferralGetUniqueByReferreeHash, referreeHash)

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

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
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
