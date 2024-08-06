package foree_transaction

import "xue.io/go-pay/app/foree/types"

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
}

type FeeJoint struct {
	ID            int64
	Amt           types.AmountData
	Description   string
	TransactionId int64
	OwnerId       int64
}
