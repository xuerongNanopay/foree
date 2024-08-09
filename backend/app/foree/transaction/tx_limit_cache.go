package transaction

import (
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLTxLimitCacheInsert = `
		INSERT INTO tx_limit_cache
		(
			identity, used_amount, used_currency,
			max_amount, max_currency
		) VALUES (?,?,?,?,?)
	`
	sQLTxLimitCacheGetUniqueByIdentity = `
		SELECT
			t.id, t.identity, t.used_amount, t.used_currency,
			t.max_amount, t.max_currency,
			t.create_at, t.update_at
		FROM tx_limit_cache t
		WHERE identity = ?
	`
)

// Improve the performance of limit check.
// We down need to do the range query over transaction table to aggregate current usage.
// Identity format {anyUniqueId}_YYYY_MM_DD
type TxLimitCache struct {
	ID       int64
	Identity string
	UsedAmt  types.AmountData
	MaxAmt   types.AmountData
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

func NewTxLimitCacheRepo(db *sql.DB) *TxLimitCacheRepo {
	return &TxLimitCacheRepo{db: db}
}

type TxLimitCacheRepo struct {
	db *sql.DB
}

func (repo *TxLimitCacheRepo) InsertTxLimitCache(tx TxLimitCache) (int64, error) {
	result, err := repo.db.Exec(
		sQLTxLimitCacheInsert,
		tx.Identity,
		tx.UsedAmt.Amount,
		tx.UsedAmt.Curreny,
		tx.MaxAmt.Amount,
		tx.MaxAmt.Curreny,
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

func scanRowIntoTxLimitCache(rows *sql.Rows) (*TxLimitCache, error) {
	tx := new(TxLimitCache)
	err := rows.Scan(
		&tx.ID,
		&tx.Identity,
		&tx.UsedAmt.Amount,
		&tx.UsedAmt.Curreny,
		&tx.MaxAmt.Amount,
		&tx.MaxAmt.Curreny,
		&tx.CreateAt,
		&tx.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func GenerateIdentity(referenceId int64) string {
	now := time.Now()
	loc, err := time.LoadLocation("America/Toronto")
	if err == nil {
		now = now.In(loc)
	}
	return fmt.Sprintf("%v_%s", referenceId, now.Format(time.DateOnly))
}
