package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

//TODO: Redesign? Daily limit should be account base.

const (
	sQLTxLimitGetUniqueByName = `
		SELECT
			l.name, l.min_amt_amount, l.min_amt_currency,
			l.max_amt_amount, l.max_amt_currency,, l.is_enable,
			l.created_at, l.updated_at
		FROM tx_limit l
		where l.name = ?
	`
)

type TxLimit struct {
	Name      string           `json:"name"`
	MinAmt    types.AmountData `json:"minLimit"`
	MaxAmt    types.AmountData `json:"maxLimit"`
	IsEnable  bool             `json:"isEnable"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdateAt  time.Time        `json:"updatedAt"`
}

func (l *TxLimit) IsViolateLimit(amt types.AmountData) bool {
	if !l.IsEnable {
		return false
	}

	if amt.Amount < l.MinAmt.Amount {
		return false
	}

	if amt.Amount > l.MaxAmt.Amount {
		return false
	}

	return true
}

func NewTxLimitRepo(db *sql.DB) *TxLimitRepo {
	return &TxLimitRepo{db: db}
}

type TxLimitRepo struct {
	db *sql.DB
}

func (repo *TxLimitRepo) GetUniqueTxLimitByName(name string) (*TxLimit, error) {
	rows, err := repo.db.Query(sQLTxLimitGetUniqueByName, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *TxLimit

	for rows.Next() {
		f, err = scanRowIntoTxLimit(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.Name == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoTxLimit(rows *sql.Rows) (*TxLimit, error) {
	l := new(TxLimit)
	err := rows.Scan(
		&l.Name,
		&l.MinAmt.Amount,
		&l.MinAmt.Currency,
		&l.MaxAmt.Amount,
		&l.MaxAmt.Currency,
		&l.IsEnable,
		&l.CreatedAt,
		&l.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return l, nil
}
