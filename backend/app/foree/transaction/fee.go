package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	SQLFeeGetAll = `
	
	`
	SQLFeeGetUniqueById = `
	
	`
)

type FeeType string
type FeeOperator string

const (
	FeeTypeFixCost      = "FIX_COST"
	FeeTypeVariableCost = "VARIABLE_COST"
)

const (
	FeeOperatorLTE FeeOperator = "LTE"
	FeeOperatorLT  FeeOperator = "LT"
	FeeOperatorGTE FeeOperator = "GTE"
	FeeOperatorGT  FeeOperator = "GT"
)

type Fee struct {
	ID                        string
	Description               string
	Type                      FeeType
	Condition                 FeeOperator
	ConditionAmt              types.AmountData
	FeeAmt                    types.AmountData
	IsApplyInConditionAmtOnly bool
	IsEnable                  bool
	CreateAt                  time.Time `json:"createAt"`
	UpdateAt                  time.Time `json:"updateAt"`
}

type FeeJoint struct {
	ID             int64
	FeeId          string
	FeeDescription string
	Amt            types.AmountData
	TransactionId  int64
	OwnerId        int64
	CreateAt       time.Time `json:"createAt"`
	UpdateAt       time.Time `json:"updateAt"`
}

func NewFee(db *sql.DB) *FeeRepo {
	return &FeeRepo{db: db}
}

type FeeRepo struct {
	db *sql.DB
}
