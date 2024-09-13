package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/constant"
)

const (
	sQLFeeGetAll = `
		SELECT
			f.name, f.description, f.fee_group, f.type, f.condition,
			f.condition_amount, f.condition_currency,
			f.ratio, f.is_apply_in_condition_amount_only
			f.is_enable, f.created_at, f.updated_at
		FROM fees as f
	`
	sQLFeeGetUniqueByFeeGroup = `
		SELECT
			f.name, f.description, f.fee_group, f.type, f.condition,
			f.condition_amount, f.condition_currency,
			f.ratio, f.is_apply_in_condition_amount_only
			f.is_enable, f.created_at, f.updated_at
		FROM fees as f
		Where f.name = ?
	`
	sQLFeeJointInsert = `
		INSERT INTO fees
		(
			fee_name, description, amount, currency,
			parent_tx_id, owner_id
		) VALUES(?,?,?,?,?,?)
	`
	sQLFeeJointGetByParentTxId = `
		SELECT
			f.id, f.fee_name, f.description, f.amount, f.currency,
			f.parent_tx_id, f.owner_id, f.created_at, f.updated_at
		FROM fee_joint as f
		Where f.parent_tx_id = ?
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
	Name                      string           `json:"name"`
	Description               string           `json:"description"`
	FeeGroup                  string           `json:"feeGroup"`
	Type                      FeeType          `json:"type"`
	Condition                 FeeOperator      `json:"condition"`
	ConditionAmt              types.AmountData `json:"conditionAmt"`
	Ratio                     types.Amount     `json:"ratio"`
	IsApplyInConditionAmtOnly bool             `json:"isApplyInConditionAmtOnly"` //TODO: support in future.
	IsEnable                  bool             `json:"isEnable"`
	CreatedAt                 time.Time        `json:"createdAt"`
	UpdatedAt                 time.Time        `json:"updatedAt"`
}

func (f *Fee) MaybeApplyFee(amt types.AmountData) (*FeeJoint, error) {
	if !f.IsEnable {
		return nil, nil
	}

	if amt.Currency != f.ConditionAmt.Currency {
		return nil, fmt.Errorf("Fee should apply in same currency, expect `%s` but ``%s", f.ConditionAmt.Currency, amt.Currency)
	}
	cond, err := applyFeeOperator(f.Condition)
	if err != nil {
		return nil, err
	}
	if !cond(f.ConditionAmt.Amount, amt.Amount) {
		return nil, nil
	}

	switch f.Type {
	case FeeTypeFixCost:
		return &FeeJoint{
			FeeName: f.Name,
			Amt: types.AmountData{
				Amount:   f.Ratio,
				Currency: f.ConditionAmt.Currency,
			},
		}, nil
	case FeeTypeVariableCost:
		return &FeeJoint{
			FeeName: f.Name,
			Amt: types.AmountData{
				Amount:   types.Amount(math.Round(float64(f.Ratio*amt.Amount)*100.0) / 100.0),
				Currency: f.ConditionAmt.Currency,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown fee type `%s`", string(f.Type))
	}
}

func applyFeeOperator(op FeeOperator) (func(l, r types.Amount) bool, error) {
	switch op {
	case FeeOperatorLTE:
		return func(l, r types.Amount) bool {
			return float64(l) <= float64(r)
		}, nil
	case FeeOperatorLT:
		return func(l, r types.Amount) bool {
			return float64(l) < float64(r)
		}, nil
	case FeeOperatorGTE:
		return func(l, r types.Amount) bool {
			return float64(l) >= float64(r)
		}, nil
	case FeeOperatorGT:
		return func(l, r types.Amount) bool {
			return float64(l) > float64(r)
		}, nil
	default:
		return nil, fmt.Errorf("unknown fee operator `%s`", string(op))
	}
}

type FeeJoint struct {
	ID          int64            `json:"id"`
	FeeName     string           `json:"feeName"`
	Description string           `json:"description"`
	Amt         types.AmountData `json:"amt"`
	ParentTxId  int64            `json:"parentTxId"`
	OwnerId     int64            `json:"ownerId"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

func NewFeeRepo(db *sql.DB) *FeeRepo {
	return &FeeRepo{db: db}
}

type FeeRepo struct {
	db *sql.DB
}

func NewFeeJointRepo(db *sql.DB) *FeeJointRepo {
	return &FeeJointRepo{db: db}
}

type FeeJointRepo struct {
	db *sql.DB
}

func (repo *FeeRepo) GetUniqueFeeByName(ctx context.Context, name string) (*Fee, error) {
	rows, err := repo.db.Query(sQLFeeGetUniqueByFeeGroup, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Fee

	for rows.Next() {
		f, err = scanRowIntoFee(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.Name == "" {
		return nil, nil
	}

	return f, nil
}

func (repo *FeeRepo) GetAllFee() ([]*Fee, error) {
	rows, err := repo.db.Query(sQLFeeGetAll)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fees := make([]*Fee, 0)
	for rows.Next() {
		p, err := scanRowIntoFee(rows)
		if err != nil {
			return nil, err
		}
		fees = append(fees, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fees, nil
}

func (repo *FeeJointRepo) InsertFeeJoint(ctx context.Context, feeJoint FeeJoint) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLFeeJointInsert,
			feeJoint.FeeName,
			feeJoint.Description,
			feeJoint.Amt.Amount,
			feeJoint.Amt.Currency,
			feeJoint.ParentTxId,
			feeJoint.OwnerId,
		)
	} else {
		result, err = repo.db.Exec(
			sQLFeeJointInsert,
			feeJoint.FeeName,
			feeJoint.Description,
			feeJoint.Amt.Amount,
			feeJoint.Amt.Currency,
			feeJoint.ParentTxId,
			feeJoint.OwnerId,
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

func (repo *FeeJointRepo) GetAllFeeJointbyParentTxId(transactionId int64) ([]*FeeJoint, error) {
	rows, err := repo.db.Query(sQLFeeJointGetByParentTxId, transactionId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeJoints := make([]*FeeJoint, 0)
	for rows.Next() {
		p, err := scanRowIntoFeeJoint(rows)
		if err != nil {
			return nil, err
		}
		feeJoints = append(feeJoints, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feeJoints, nil
}

func scanRowIntoFee(rows *sql.Rows) (*Fee, error) {
	u := new(Fee)
	err := rows.Scan(
		&u.Name,
		&u.Description,
		&u.FeeGroup,
		&u.Type,
		&u.Condition,
		&u.ConditionAmt.Amount,
		&u.ConditionAmt.Currency,
		&u.Ratio,
		&u.IsApplyInConditionAmtOnly,
		&u.IsEnable,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func scanRowIntoFeeJoint(rows *sql.Rows) (*FeeJoint, error) {
	u := new(FeeJoint)
	err := rows.Scan(
		&u.ID,
		&u.FeeName,
		&u.Description,
		&u.Amt.Amount,
		&u.Amt.Currency,
		&u.ParentTxId,
		&u.OwnerId,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
