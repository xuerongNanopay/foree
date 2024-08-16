package transaction

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLPromoCodeGetUniqueByCode = `
		SELECT
			P.code, p.description, p.min_amount, p.min_currency, p.limit_per_acc,
			p.start_time, p.end_time, p.create_at, p.update_at
		FROM promo_code as p
		Where p.code = ?
	`
)

type PromoCode struct {
	Code        string           `json:"code"`
	Description string           `json:"description"`
	MinAmt      types.AmountData `json:"minAmt"`
	LimitPerAcc int              `json:"limit_per_acc"`
	StartTime   time.Time        `json:"startTime"`
	EndTime     time.Time        `json:"endTime"`
	CreateAt    time.Time        `json:"createAt"`
	UpdateAt    time.Time        `json:"updateAt"`
}

func NewPromoCodeRepo(db *sql.DB) *PromoCodeRepo {
	return &PromoCodeRepo{db: db}
}

type PromoCodeRepo struct {
	db *sql.DB
}

func (repo *PromoCodeRepo) GetUniquePromoCodeByCode(ctx context.Context, code string) (*PromoCode, error) {
	rows, err := repo.db.Query(sQLPromoCodeGetUniqueByCode, code)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *PromoCode

	for rows.Next() {
		f, err = scanRowIntoPromoCode(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.Code == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoPromoCode(rows *sql.Rows) (*PromoCode, error) {
	p := new(PromoCode)
	err := rows.Scan(
		&p.Code,
		&p.Description,
		&p.MinAmt.Amount,
		&p.MinAmt.Curreny,
		&p.StartTime,
		&p.EndTime,
		&p.CreateAt,
		&p.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}
