package referral

import "time"

type Referral struct {
	ID           int64     `json:"id"`
	ReferralCode string    `json:"referralCode"`
	ReferrerId   int64     `json:"referrerId"`
	ReferreeId   int64     `json:"referreeId"`
	CreateAt     time.Time `json:"createAt"`
	UpdateAt     time.Time `json:"updateAt"`
}
