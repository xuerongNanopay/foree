package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	SQLRewardInsert = `
		INSERT INTO rewards
		(
			type, description, amount, currency,
			status, is_redeemed, owner_id, transaction_id
		) VALUES(?,?,?,?,?,?,?,?)
	`
	SQLRewardGetAll = `
		SELECT
			r.id, r.type, r.description, r.amount, r.currency,
			r.status, r.is_redeemed, r.owner_id, r.transaction_id,
			r.expire_at, f.create_at, f.update_at
		FROM rewards as r
	`
	SQLRewardGetAllByTransactionId = `
		SELECT
			r.id, r.type, r.description, r.amount, r.currency,
			r.status, r.is_redeemed, r.owner_id, r.transaction_id,
			r.expire_at, f.create_at, f.update_at
		FROM rewards as r
		Where r.transaction_id=?
	`
	SQLRewardGetAllByOwnerId = `
		SELECT
			r.id, r.type, r.description, r.amount, r.currency,
			r.status, r.is_redeemed, r.owner_id, r.transaction_id,
			r.expire_at, f.create_at, f.update_at
		FROM rewards as r
		Where r.owner_id=?
	`
)

type RewardType string

const (
	RewardTypeSignUp  RewardType = "SIGN_UP_REWARD"
	RewardTypeReferal RewardType = "REFERAL_REWARD"
)

type RewardStatus string

const (
	RewardStatusEnable  = "ENABLE"
	RewardStatusDisable = "DISABLE"
)

type Reward struct {
	ID            string
	Type          string
	Description   string
	Amt           types.AmountData
	Status        RewardStatus
	IsRedeemed    bool
	OwnerId       int64
	TransactionId int64
	ExpireAt      time.Time `json:"expireAt"`
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
}

func NewRewardRepo(db *sql.DB) *RewardRepo {
	return &RewardRepo{db: db}
}

type RewardRepo struct {
	db *sql.DB
}
