package promotion

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/constant"
)

const (
	sQLPromotionRewardJointInsert = `
		INSERT INTO promotion_reward_joint(	 
			promotion_id, promotion_version, reward_id
		) VALUES (?,?,?)
	`
	sQLPromotionRewardJointCountByPromotionIdAndPromotionVersion = `
		SELECT
			COUNT(*)
		FROM promotion_reward_joint p
		WHERE p.promotion_id = ? AND p.promotion_version = ?
	`
)

type PromotionRewardJoint struct {
	ID               int64      `json:"id"`
	PromotionId      int64      `json:"promotionId"`
	PromotionVersion int64      `json:"promotionVersion"`
	RewardId         int64      `json:"rewardId"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}

func NewPromotionRewardJointRepo(db *sql.DB) *PromotionRewardJointRepo {
	return &PromotionRewardJointRepo{db: db}
}

type PromotionRewardJointRepo struct {
	db *sql.DB
}

func (repo *PromotionRewardJointRepo) InsertPromotionRewardJoint(ctx context.Context, j PromotionRewardJoint) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result
	if ok {
		result, err = dTx.Exec(
			sQLPromotionRewardJointInsert,
			j.PromotionId,
			j.PromotionVersion,
			j.RewardId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLPromotionRewardJointInsert,
			j.PromotionId,
			j.PromotionVersion,
			j.RewardId,
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

func (repo *PromotionRewardJointRepo) CountPromotionRewardJointByPromotionIdAndPromotionVersion(ctx context.Context, pi int64, pv int) (int, error) {
	var count int
	if err := repo.db.QueryRow(sQLPromotionRewardJointCountByPromotionIdAndPromotionVersion, pi, pv).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
