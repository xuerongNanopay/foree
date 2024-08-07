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
		Where r.transaction_id = ?
	`
	SQLRewardGetAllByOwnerId = `
		SELECT
			r.id, r.type, r.description, r.amount, r.currency,
			r.status, r.is_redeemed, r.owner_id, r.transaction_id,
			r.expire_at, f.create_at, f.update_at
		FROM rewards as r
		Where r.owner_id = ?
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

func (repo *FeeRepo) InsertReward(reward Reward) (int64, error) {
	result, err := repo.db.Exec(
		SQLRewardInsert,
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

func (repo *FeeRepo) GetAllReward() ([]*Reward, error) {
	rows, err := repo.db.Query(SQLRewardGetAll)

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

func (repo *FeeRepo) GetAllRewardByTransactionId(transactionId int64) ([]*Reward, error) {
	rows, err := repo.db.Query(SQLRewardGetAllByTransactionId, transactionId)

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

func (repo *FeeRepo) GetAllRewardByOwnerId(owenerId int64) ([]*Reward, error) {
	rows, err := repo.db.Query(SQLRewardGetAllByOwnerId, owenerId)

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
