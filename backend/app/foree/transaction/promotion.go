package foree_transaction

import (
	"time"

	"xue.io/go-pay/app/foree/types"
)

type PromotionStatus string

const (
	PromotionStatusEnable  = "ENABLE"
	PromotionStatusDisable = "DISABLE"
)

type Promotion struct {
	ID                 string
	Description        string
	Amt                types.AmountData
	Status             PromotionStatus
	IsRedeemed         bool
	OwnerId            int64
	ForeeTransactionId int64
	ExpireAt           time.Time `json:"expireAt"`
	CreateAt           time.Time `json:"createAt"`
	UpdateAt           time.Time `json:"updateAt"`
}
