package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/constant"
	time_util "xue.io/go-pay/util/time"
)

const (
	sQLDailyTxLimitInsert = `
		INSERT INTO daily_tx_limit
		(
			reference, used_amount, used_currency,
			max_amount, max_currency
		) VALUES (?,?,?,?,?)
	`
	sQLDailyTxLimitUpdateByReference = `
        UPDATE daily_tx_limit SET 
            used_amount = ?
        WHERE reference = ?
	`
	sQLDailyTxLimitGetUniqueByReference = `
		SELECT
			t.id, t.reference, t.used_amount, t.used_currency,
			t.max_amount, t.max_currency,
			t.create_at, t.update_at
		FROM daily_tx_limit t
		WHERE reference = ?
	`
)

// Improve the performance of limit check.
// We down need to do the range query over transaction table to aggregate current usage.
// Reference format {anyUniqueId}_YYYY_MM_DD
type DailyTxLimit struct {
	ID        int64            `json:"id"`
	Reference string           `json:"reference"`
	UsedAmt   types.AmountData `json:"usedAmt"`
	MaxAmt    types.AmountData `json:"maxAmt"`
	CreateAt  time.Time        `json:"createAt"`
	UpdateAt  time.Time        `json:"updateAt"`
}

func NewDailyTxLimitRepo(db *sql.DB) *DailyTxLimitRepo {
	return &DailyTxLimitRepo{db: db}
}

type DailyTxLimitRepo struct {
	db *sql.DB
}

func (repo *DailyTxLimitRepo) InsertDailyTxLimit(ctx context.Context, tx DailyTxLimit) (int64, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error
	var result sql.Result

	if ok {
		result, err = dTx.Exec(
			sQLDailyTxLimitInsert,
			tx.Reference,
			tx.UsedAmt.Amount,
			tx.UsedAmt.Currency,
			tx.MaxAmt.Amount,
			tx.MaxAmt.Currency,
		)
	} else {
		result, err = repo.db.Exec(
			sQLDailyTxLimitInsert,
			tx.Reference,
			tx.UsedAmt.Amount,
			tx.UsedAmt.Currency,
			tx.MaxAmt.Amount,
			tx.MaxAmt.Currency,
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

func (repo *DailyTxLimitRepo) UpdateDailyTxLimitById(ctx context.Context, tx DailyTxLimit) error {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var err error

	if ok {
		_, err = dTx.Exec(sQLDailyTxLimitUpdateByReference, tx.UsedAmt.Amount, tx.ID)
	} else {
		_, err = repo.db.Exec(sQLDailyTxLimitUpdateByReference, tx.UsedAmt.Amount, tx.ID)
	}
	if err != nil {
		return err
	}
	return nil
}

func (repo *DailyTxLimitRepo) GetUniqueDailyTxLimitByReference(ctx context.Context, reference string) (*DailyTxLimit, error) {
	rows, err := repo.db.Query(sQLDailyTxLimitGetUniqueByReference, reference)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *DailyTxLimit

	for rows.Next() {
		f, err = scanRowIntoDailyTxLimit(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.ID == 0 {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoDailyTxLimit(rows *sql.Rows) (*DailyTxLimit, error) {
	tx := new(DailyTxLimit)
	err := rows.Scan(
		&tx.ID,
		&tx.Reference,
		&tx.UsedAmt.Amount,
		&tx.UsedAmt.Currency,
		&tx.MaxAmt.Amount,
		&tx.MaxAmt.Currency,
		&tx.CreateAt,
		&tx.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func GenerateDailyTxLimitReference(userId int64) string {
	now := time_util.NowInToronto()
	return fmt.Sprintf("%v_%s", userId, now.Format(time.DateOnly))
}
