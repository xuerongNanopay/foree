package referral

import "time"

type ReferralRewardJoint struct {
	ID         int64      `json:"id"`
	ReferralId int64      `json:"referralId"`
	ReferrerId int64      `json:"referrerId"`
	RefereeId  int64      `json:"refereeId"`
	RewardId   int64      `json:"rewardId"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
