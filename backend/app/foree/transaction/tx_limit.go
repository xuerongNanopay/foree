package transaction

import (
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLTxLimitGetUniqueById = `
		SELECT
			l.id, l.amount, l.currency,
			l.is_min_limit, l.is_enable,
			l.create_at, l.update_at
		FROM tx_limit l
		where l.id = ?
	`
)

type TxLimit struct {
	ID         string
	Amt        types.AmountData
	IsMinLimit bool
	IsEnable   bool
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
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

	if f.ID == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoTxLimit(rows *sql.Rows) (*TxLimit, error) {
	l := new(TxLimit)
	err := rows.Scan(
		&l.ID,
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
