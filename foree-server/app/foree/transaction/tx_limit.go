package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLTxLimitGetUniqueByLimitGroup = `
		SELECT
			l.name, l.limit_group, l.min_amount, l.min_currency,
			l.max_amount, l.max_currency, l.created_at, l.updated_at
		FROM tx_limit l
		where l.limit_group = ?
	`
)

type TxLimit struct {
	Name       string           `json:"name"`
	LimitGroup string           `json:"limitGroup"`
	MinAmt     types.AmountData `json:"minLimit"`
	MaxAmt     types.AmountData `json:"maxLimit"`
	CreatedAt  *time.Time       `json:"createdAt"`
	UpdatedAt  *time.Time       `json:"updatedAt"`
}

func NewTxLimitRepo(db *sql.DB) *TxLimitRepo {
	return &TxLimitRepo{db: db}
}

type TxLimitRepo struct {
	db *sql.DB
}

func (repo *TxLimitRepo) GetUniqueTxLimitByLimitGroup(limitGroup string) (*TxLimit, error) {
	rows, err := repo.db.Query(sQLTxLimitGetUniqueByLimitGroup, limitGroup)

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

	if f == nil || f.Name == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoTxLimit(rows *sql.Rows) (*TxLimit, error) {
	l := new(TxLimit)
	err := rows.Scan(
		&l.Name,
		&l.LimitGroup,
		&l.MinAmt.Amount,
		&l.MinAmt.Currency,
		&l.MaxAmt.Amount,
		&l.MaxAmt.Currency,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return l, nil
}
