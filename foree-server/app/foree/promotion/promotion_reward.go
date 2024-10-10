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
	sQLPromotionRewardJointGetUniqueByPromotionIdAndOwnerId = `
		SELECT
			p.id, p.promotion_id, p.promotion_version, p.reward_id, p.referrer_id, p.referee_id
			p.owner_id, p.created_at, p.updated_at
		FROM FROM promotion_reward_joint p
		WHERE p.promotion_id = ? AND p.owner_id = ?
	`
	sQLPromotionRewardJointGetUniqueByPromotionIdAndReferrerIdAndRefereeId = `
		SELECT
			p.id, p.promotion_id, p.promotion_version, p.reward_id, p.referrer_id, p.referee_id
			p.owner_id, p.created_at, p.updated_at
		FROM FROM promotion_reward_joint p
		WHERE p.promotion_id = ? AND p.referrer_id = ? AND p.referee_id = ?
	`
)

type PromotionRewardJoint struct {
	ID               int64      `json:"id"`
	PromotionId      int64      `json:"promotionId"`
	PromotionVersion int        `json:"promotionVersion"`
	RewardId         int64      `json:"rewardId"`
	ReferrerId       int64      `json:"referrerId"`
	RefereeId        int64      `json:"refereeId"`
	OwnerId          int64      `json:"ownerId"`
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

func (repo *PromotionRewardJointRepo) GetUniquePromotionRewardJointByPromotionIdAndOwnerId(promotionId int64, ownerId int64) (*PromotionRewardJoint, error) {
	rows, err := repo.db.Query(sQLPromotionRewardJointGetUniqueByPromotionIdAndOwnerId, promotionId, ownerId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *PromotionRewardJoint

	for rows.Next() {
		f, err = scanRowIntoPromotionRewardJoint(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoPromotionRewardJoint(rows *sql.Rows) (*PromotionRewardJoint, error) {
	prj := new(PromotionRewardJoint)
	err := rows.Scan(
		&prj.ID,
		&prj.PromotionId,
		&prj.PromotionVersion,
		&prj.RewardId,
		&prj.ReferrerId,
		&prj.RefereeId,
		&prj.OwnerId,
		&prj.CreatedAt,
		&prj.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return prj, nil
}
