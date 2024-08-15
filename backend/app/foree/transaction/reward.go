package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLRewardInsert = `
		INSERT INTO rewards
		(
			type, description, amount, currency,
			status, is_redeemed, owner_id, transaction_id
		) VALUES(?,?,?,?,?,?,?,?)
	`
	sQLRewardGetAll = `
		SELECT
			r.id, r.type, r.description, r.amount, r.currency,
			r.status, r.is_redeemed, r.owner_id, r.transaction_id,
			r.expire_at, f.create_at, f.update_at
		FROM rewards as r
	`
	sQLRewardGetUniqueByTransactionId = `
	SELECT
		r.id, r.type, r.description, r.amount, r.currency,
		r.status, r.is_redeemed, r.owner_id, r.transaction_id,
		r.expire_at, f.create_at, f.update_at
	FROM rewards as r
	Where r.id = ?
`
	sQLRewardGetAllByTransactionId = `
		SELECT
			r.id, r.type, r.description, r.amount, r.currency,
			r.status, r.is_redeemed, r.owner_id, r.transaction_id,
			r.expire_at, f.create_at, f.update_at
		FROM rewards as r
		Where r.transaction_id = ?
	`
	sQLRewardGetAllUnredeemByOwnerId = `
		SELECT
			r.id, r.type, r.description, r.amount, r.currency,
			r.status, r.is_redeemed, r.owner_id, r.transaction_id,
			r.expire_at, f.create_at, f.update_at
		FROM rewards as r
		Where r.owner_id = ? AND r.is_redeemed = FALSE
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
	ID            int64            `json:"id"`
	Type          string           `json:"type"`
	Description   string           `json:"description"`
	Amt           types.AmountData `json:"amt"`
	Status        RewardStatus     `json:"status"`
	IsRedeemed    bool             `json:"isRedeemed"`
	OwnerId       int64            `json:"ownerId"`
	TransactionId int64            `json:"transactionId"`
	ExpireAt      time.Time        `json:"expireAt"`
	CreateAt      time.Time        `json:"createAt"`
	UpdateAt      time.Time        `json:"updateAt"`
}

func NewRewardRepo(db *sql.DB) *RewardRepo {
	return &RewardRepo{db: db}
}

type RewardRepo struct {
	db *sql.DB
}

func (repo *FeeRepo) InsertReward(ctx context.Context, reward Reward) (int64, error) {
	result, err := repo.db.Exec(
		sQLRewardInsert,
		reward.Type,
		reward.Description,
		reward.Amt.Amount,
		reward.Amt.Curreny,
		reward.Status,
		reward.IsRedeemed,
		reward.OwnerId,
		reward.TransactionId,
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

func (repo *FeeRepo) GetAllReward(ctx context.Context) ([]*Reward, error) {
	rows, err := repo.db.Query(sQLRewardGetAll)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := make([]*Reward, 16)
	for rows.Next() {
		p, err := scanRowIntoReward(rows)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rewards, nil
}

func (repo *FeeRepo) GetAllRewardByTransactionId(ctx context.Context, transactionId int64) ([]*Reward, error) {
	rows, err := repo.db.Query(sQLRewardGetAllByTransactionId, transactionId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := make([]*Reward, 16)
	for rows.Next() {
		p, err := scanRowIntoReward(rows)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rewards, nil
}

func (repo *FeeRepo) GetAllUnredeemRewardByOwnerId(ctx context.Context, ownerId int64) ([]*Reward, error) {
	rows, err := repo.db.Query(sQLRewardGetAllUnredeemByOwnerId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := make([]*Reward, 16)
	for rows.Next() {
		p, err := scanRowIntoReward(rows)
		if err != nil {
			return nil, err
		}
		rewards = append(rewards, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rewards, nil
}

func scanRowIntoReward(rows *sql.Rows) (*Reward, error) {
	u := new(Reward)
	err := rows.Scan(
		&u.ID,
		&u.Type,
		&u.Description,
		&u.Amt.Amount,
		&u.Amt.Curreny,
		&u.Status,
		&u.IsRedeemed,
		&u.OwnerId,
		&u.TransactionId,
		&u.ExpireAt,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
