package promotion

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/constant"
	uuid_util "xue.io/go-pay/util/uuid"
)

const (
	sQLRewardInsert = `
		INSERT INTO rewards
		(
			s_id, type, description, amount, currency,
			status, owner_id, expire_at
		) VALUES(?,?,?,?,?,?,?,?)
	`
	sQLRewardUpdateById = `
		UPDATE rewards SET
			status = ?, applied_transaction_id = ?
		WHERE id = ?
	`
	sQLRewardGetUniqueById = `
		SELECT
			r.id, r.s_id, r.type, r.description, r.amount, r.currency,
			r.status, r.owner_id, r.applied_transaction_id,
			r.expire_at, r.created_at, r.updated_at
		FROM rewards as r
		Where r.id = ?
	`
	sQLRewardGetAllByOwnerIdAndSids = `
		SELECT
			r.id, r.s_id, r.type, r.description, r.amount, r.currency,
			r.status, r.owner_id, r.applied_transaction_id,
			r.expire_at, r.created_at, r.updated_at
		FROM rewards as r
		Where r.owner_id = ? AND r.s_id in (%v)
	`
	sQLRewardGetAllByOwnerIdAndIds = `
		SELECT
			r.id, r.s_id, r.type, r.description, r.amount, r.currency,
			r.status, r.owner_id, r.applied_transaction_id,
			r.expire_at, r.created_at, r.updated_at
		FROM rewards as r
		Where r.owner_id = ? AND r.id in (%v)
	`
	sQLRewardGetAllByAppliedTransactionId = `
		SELECT
			r.id, r.s_id, r.type, r.description, r.amount, r.currency,
			r.status, r.owner_id, r.applied_transaction_id,
			r.expire_at, r.created_at, r.updated_at
		FROM rewards as r
		Where r.applied_transaction_id = ?
	`
	sQLRewardGetAllActiveByOwnerId = `
		SELECT
			r.id, r.s_id, r.type, r.description, r.amount, r.currency,
			r.status, r.owner_id, r.applied_transaction_id,
			r.expire_at, r.created_at, r.updated_at
		FROM rewards as r
		Where r.owner_id = ? AND r.status = 'ACTIVE'
	`
)

type RewardType string

const (
	RewardTypeOnboard   string = "ONBOARD_REWARD"
	RewardTypeReferrer  string = "REFERRER_REWARD"
	RewardTypeReferee   string = "REFEREE_REWARD"
	RewardTypeTx        string = "TX_REWARD"
	RewardTypePromoCode string = "PROMO_CODE_REWARD"
)

type RewardStatus string

const (
	RewardStatusInitial  = "INITIAL"
	RewardStatusActive   = "ACTIVE"
	RewardStatusPending  = "PENDING"
	RewardStatusRedeemed = "REDEEMED"
	RewardStatusDelete   = "DELETE"
)

type Reward struct {
	ID                   int64            `json:"id"`
	SID                  string           `json:"sId"`
	Type                 string           `json:"type"`
	Description          string           `json:"description"`
	Amt                  types.AmountData `json:"amt"`
	Status               RewardStatus     `json:"status"`
	OwnerId              int64            `json:"ownerId"`
	AppliedTransactionId int64            `json:"appliedTransactionId"`
	ExpireAt             *time.Time       `json:"expireAt"`
	CreatedAt            *time.Time       `json:"createdAt"`
	UpdatedAt            *time.Time       `json:"updatedAt"`
}

func NewRewardRepo(db *sql.DB) *RewardRepo {
	return &RewardRepo{db: db}
}

type RewardRepo struct {
	db *sql.DB
}

func (repo *RewardRepo) InsertReward(ctx context.Context, reward Reward) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLRewardInsert,
			uuid_util.UUID(),
			reward.Type,
			reward.Description,
			reward.Amt.Amount,
			reward.Amt.Currency,
			reward.Status,
			reward.OwnerId,
			reward.ExpireAt,
		)
	} else {
		result, err = repo.db.Exec(
			sQLRewardInsert,
			uuid_util.UUID(),
			reward.Type,
			reward.Description,
			reward.Amt.Amount,
			reward.Amt.Currency,
			reward.Status,
			reward.OwnerId,
			reward.ExpireAt,
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

func (repo *RewardRepo) UpdateRewardTxById(ctx context.Context, reward Reward) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	if ok {
		_, err = dTx.Exec(sQLRewardUpdateById, reward.Status, reward.AppliedTransactionId, reward.ID)
	} else {
		_, err = repo.db.Exec(sQLRewardUpdateById, reward.Status, reward.AppliedTransactionId, reward.ID)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo *RewardRepo) GetUniqueRewardById(ctx context.Context, id int64) (*Reward, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(sQLRewardGetUniqueById, id)
	} else {
		rows, err = repo.db.Query(sQLRewardGetUniqueById, id)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Reward

	for rows.Next() {
		f, err = scanRowIntoReward(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func (repo *RewardRepo) GetAllRewardByOwnerIdAndIds(ctx context.Context, ownerId int64, ids ...int64) ([]*Reward, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	args := make([]interface{}, len(ids)+1)
	args[0] = ownerId
	p := make([]string, len(ids))
	for i, id := range ids {
		args[i+1] = id
		p[i] = "?"
	}

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(fmt.Sprintf(sQLRewardGetAllByOwnerIdAndIds, strings.Join(p, ",")), args...)

	} else {
		rows, err = repo.db.Query(fmt.Sprintf(sQLRewardGetAllByOwnerIdAndIds, strings.Join(p, ",")), args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := make([]*Reward, 0)
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

func (repo *RewardRepo) GetAllRewardByOwnerIdAndSids(ctx context.Context, ownerId int64, sids ...string) ([]*Reward, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	args := make([]interface{}, len(sids)+1)
	args[0] = ownerId
	p := make([]string, len(sids))
	for i, ref := range sids {
		args[i+1] = ref
		p[i] = "?"
	}

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(fmt.Sprintf(sQLRewardGetAllByOwnerIdAndSids, strings.Join(p, ",")), args...)

	} else {
		rows, err = repo.db.Query(fmt.Sprintf(sQLRewardGetAllByOwnerIdAndSids, strings.Join(p, ",")), args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := make([]*Reward, 0)
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

func (repo *RewardRepo) GetAllRewardByAppliedTransactionId(ctx context.Context, appliedTransactionId int64) ([]*Reward, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(sQLRewardGetAllByAppliedTransactionId, appliedTransactionId)

	} else {
		rows, err = repo.db.Query(sQLRewardGetAllByAppliedTransactionId, appliedTransactionId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := make([]*Reward, 0)
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

func (repo *RewardRepo) GetAllActiveRewardByOwnerId(ctx context.Context, ownerId int64) ([]*Reward, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(sQLRewardGetAllActiveByOwnerId, ownerId)

	} else {
		rows, err = repo.db.Query(sQLRewardGetAllActiveByOwnerId, ownerId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rewards := make([]*Reward, 0)
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
		&u.SID,
		&u.Type,
		&u.Description,
		&u.Amt.Amount,
		&u.Amt.Currency,
		&u.Status,
		&u.OwnerId,
		&u.AppliedTransactionId,
		&u.ExpireAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
