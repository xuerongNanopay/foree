package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLTxLimitGetUniqueById = `
		SELECT
			l.name, l.amount, l.currency,
			l.is_min_limit, l.is_enable,
			l.create_at, l.update_at
		FROM tx_limit l
		where l.name = ?
	`
)

type TxLimit struct {
	Name       string           `json:"name"`
	Amt        types.AmountData `json:"amt"`
	IsMinLimit bool             `json:"isMinLimit"`
	IsEnable   bool             `json:"isEnable"`
	CreateAt   time.Time        `json:"createAt"`
	UpdateAt   time.Time        `json:"updateAt"`
}

func (l *TxLimit) IsViolateLimit(amt types.AmountData) bool {
	return false
}

func NewTxLimitRepo(db *sql.DB) *TxLimitRepo {
	return &TxLimitRepo{db: db}
}

type TxLimitRepo struct {
	db *sql.DB
}

func (repo *InteracCITxRepo) GetUniqueTxLimitById(id int64) (*TxLimit, error) {
	rows, err := repo.db.Query(sQLTxLimitGetUniqueById, id)

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
		&l.Amt.Amount,
		&l.Amt.Curreny,
		&l.IsMinLimit,
		&l.IsEnable,
		&l.CreateAt,
		&l.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return l, nil
}
