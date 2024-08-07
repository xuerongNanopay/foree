package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	SQLFeeGetAll = `
		SELECT
			f.id, f.description, f.type, f.condition,
			f.condition_amount, f.condition_currency,
			f.fee_amount, f.fee_currency,
			f.is_apply_in_condition_amount_only
			f.is_enable, f.create_at, f.update_at
		FROM fees as f
	`
	SQLFeeGetUniqueById = `
		SELECT
			f.id, f.description, f.type, f.condition,
			f.condition_amount, f.condition_currency,
			f.fee_amount, f.fee_currency,
			f.is_apply_in_condition_amount_only
			f.is_enable, f.create_at, f.update_at
		FROM fees as f
		Where f.id=?
	`
	SQLFeeJointInsert = `
		INSERT INTO fees
		(
			fee_id, description, amount, currency,
			transaction_id, owner_id
		) VALUES(?,?,?,?,?,?)
	`
	SQLFeeJointGetByTransactionId = `
		SELECT
			f.fee_id, f.description, f.amount, f.currency,
			f.transaction_id, f.owner_id, f.create_at, f.update_at
		FROM fee_joint as f
		Where f.transaction_id=?
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
	ID            int64
	FeeId         string
	Description   string
	Amt           types.AmountData
	TransactionId int64
	OwnerId       int64
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
}

func NewFee(db *sql.DB) *FeeRepo {
	return &FeeRepo{db: db}
}

type FeeRepo struct {
	db *sql.DB
}

func (repo *FeeRepo) GetFeeById(id int64) (*Fee, error) {
	rows, err := repo.db.Query(SQLFeeGetUniqueById, id)

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

func (repo *FeeRepo) GetAllFee(id int64) ([]*Fee, error) {
	rows, err := repo.db.Query(SQLFeeGetAll)

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

func (repo *FeeRepo) InsertFeeJoint(feeJoint FeeJoint) (int64, error) {
	result, err := repo.db.Exec(
		SQLFeeJointInsert,
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

func (repo *FeeRepo) GetAllFeeJoinbyTransactionId(transactionId int64) ([]*FeeJoint, error) {
	rows, err := repo.db.Query(SQLFeeJointGetByTransactionId, transactionId)

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
