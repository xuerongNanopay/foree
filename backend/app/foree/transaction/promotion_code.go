package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

type PromoCode struct {
	Code        string           `json:"code"`
	Description string           `json:"description"`
	MinAmt      types.AmountData `json:"minAmt"`
	LimitPerAcc int              `json:"limit_per_acc"`
	StartTime   time.Time        `json:"startTime"`
	EndTime     time.Time        `json:"endTime"`
	CreateAt    time.Time        `json:"createAt"`
	UpdateAt    time.Time        `json:"updateAt"`
}

func NewPromoCodeRepo(db *sql.DB) *PromoCodeRepo {
	return &PromoCodeRepo{db: db}
}

type PromoCodeRepo struct {
	db *sql.DB
}
