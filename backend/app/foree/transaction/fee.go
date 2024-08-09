package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLFeeGetAll = `
		SELECT
			f.id, f.description, f.type, f.condition,
			f.condition_amount, f.condition_currency,
			f.fee_amount, f.fee_currency,
			f.is_apply_in_condition_amount_only
			f.is_enable, f.create_at, f.update_at
		FROM fees as f
	`
	sQLFeeGetUniqueById = `
		SELECT
			f.id, f.description, f.type, f.condition,
			f.condition_amount, f.condition_currency,
			f.fee_amount, f.fee_currency,
			f.is_apply_in_condition_amount_only
			f.is_enable, f.create_at, f.update_at
		FROM fees as f
		Where f.id = ?
	`
	sQLFeeJointInsert = `
		INSERT INTO fees
		(
			fee_id, description, amount, currency,
			transaction_id, owner_id
		) VALUES(?,?,?,?,?,?)
	`
	sQLFeeJointGetByTransactionId = `
		SELECT
			f.fee_id, f.description, f.amount, f.currency,
			f.transaction_id, f.owner_id, f.create_at, f.update_at
		FROM fee_joint as f
		Where f.transaction_id = ?
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
	ID                        string           `json:"id"`
	Description               string           `json:"description"`
	Type                      FeeType          `json:"type"`
	Condition                 FeeOperator      `json:"condition"`
	ConditionAmt              types.AmountData `json:"conditionAmt"`
	FeeAmt                    types.AmountData `json:"feeAmt"`
	IsApplyInConditionAmtOnly bool             `json:"isApplyInConditionAmtOnly"`
	IsEnable                  bool             `json:"isEnable"`
	CreateAt                  time.Time        `json:"createAt"`
	UpdateAt                  time.Time        `json:"updateAt"`
}

type FeeJoint struct {
	ID            int64            `json:"id"`
	FeeId         string           `json:"feeId"`
	Description   string           `json:"description"`
	Amt           types.AmountData `json:"amt"`
	TransactionId int64            `json:"transactionId"`
	OwnerId       int64            `json:"ownerId"`
	CreateAt      time.Time        `json:"createAt"`
	UpdateAt      time.Time        `json:"updateAt"`
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

func (repo *FeeRepo) GetUniqueFeeById(id int64) (*Fee, error) {
	rows, err := repo.db.Query(sQLFeeGetUniqueById, id)

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

	if f.ID == "" {
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

	fees := make([]*Fee, 16)
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

func (repo *FeeJointRepo) InsertFeeJoint(feeJoint FeeJoint) (int64, error) {
	result, err := repo.db.Exec(
		sQLFeeJointInsert,
		feeJoint.FeeId,
		feeJoint.Description,
		feeJoint.Amt.Amount,
		feeJoint.Amt.Curreny,
		feeJoint.TransactionId,
		feeJoint.OwnerId,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *FeeJointRepo) GetAllFeeJointbyTransactionId(transactionId int64) ([]*FeeJoint, error) {
	rows, err := repo.db.Query(sQLFeeJointGetByTransactionId, transactionId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeJoints := make([]*FeeJoint, 16)
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
		&u.ID,
		&u.Description,
		&u.Type,
		&u.Condition,
		&u.ConditionAmt.Amount,
		&u.ConditionAmt.Curreny,
		&u.FeeAmt.Amount,
		&u.FeeAmt.Curreny,
		&u.IsApplyInConditionAmtOnly,
		&u.IsEnable,
		&u.CreateAt,
		&u.UpdateAt,
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
		&u.FeeId,
		&u.Description,
		&u.Amt.Amount,
		&u.Amt.Curreny,
		&u.TransactionId,
		&u.OwnerId,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
