package transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type RewardStatus string

const (
	RewardStatusEnable  = "ENABLE"
	RewardStatusDisable = "DISABLE"
)

type Reward struct {
	ID                 string
	Description        string
	Amt                types.AmountData
	Status             RewardStatus
	IsRedeemed         bool
	OwnerId            int64
	ForeeTransactionId int64
	ExpireAt           time.Time `json:"expireAt"`
	CreateAt           time.Time `json:"createAt"`
	UpdateAt           time.Time `json:"updateAt"`
}
