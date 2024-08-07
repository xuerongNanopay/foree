package foree_auth

import "time"

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
	OwerId        string
	RefereeId     string
	IsRedeemed    bool
	ExpireAt      time.Time `json:"expireAt"`
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
}
