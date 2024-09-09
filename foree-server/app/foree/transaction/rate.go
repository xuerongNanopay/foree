package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sqlRateInsert = `
		INSERT INTO rate
		(
			id, src_amount, src_currency,
			dest_amount, dest_currency
		) VALUES(?,?,?,?,?)
	`
	sQLRateUpdateById = `
		UPDATE rate SET
			src_amount = ?, dest_amount = ?
		WHERE id = ?
	`
	sQLRateGetUniqueById = `
	    SELECT 
            r.id, r.src_amount, r.src_currency,
			r.dest_amount, r.dest_currency, 
			r.created_at, r.updated_at
        FROM rate r
        where t.id = ?
	`
)

// ID format: {src_currency}-{dest_currency}
type Rate struct {
	ID        string           `json:"id"`
	SrcAmt    types.AmountData `json:"srcAmt"`
	DestAmt   types.AmountData `json:"destAmt"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

func NewRateRepo(db *sql.DB) *RateRepo {
	return &RateRepo{db: db}
}

type RateRepo struct {
	db *sql.DB
}

func (repo *RateRepo) InsertRate(ctx context.Context, r Rate) (string, error) {
	result, err := repo.db.Exec(
		sqlRateInsert,
		r.GetId(),
		r.SrcAmt.Amount,
		strings.ToUpper(r.SrcAmt.Currency),
		r.GetForwardRate(),
		strings.ToUpper(r.DestAmt.Currency),
	)
	if err != nil {
		return "", err
	}
	_, qerr := result.LastInsertId()
	if qerr != nil {
		return "", qerr
	}
	return r.GetId(), nil
}

func (repo *RateRepo) UpdateRateById(ctx context.Context, r Rate) error {
	_, err := repo.db.Exec(sQLRateUpdateById, r.SrcAmt.Amount, r.GetForwardRate(), r.GetId())
	if err != nil {
		return err
	}
	return nil
}

func (repo *RateRepo) GetUniqueRateById(ctx context.Context, id string) (*Rate, error) {
	rows, err := repo.db.Query(sQLRateGetUniqueById, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Rate

	for rows.Next() {
		f, err = scanRowIntoRate(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.ID == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoRate(rows *sql.Rows) (*Rate, error) {
	u := new(Rate)
	err := rows.Scan(
		&u.ID,
		&u.SrcAmt.Amount,
		&u.SrcAmt.Currency,
		&u.DestAmt.Amount,
		&u.DestAmt.Currency,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Rate) GetId() string {
	return GenerateRateId(r.SrcAmt.Currency, r.DestAmt.Currency)
}

func (r *Rate) ToSummary() string {
	return fmt.Sprintf("$%.2f %s : %.2f %s", r.SrcAmt.Amount, r.SrcAmt.Currency, r.DestAmt.Amount, r.DestAmt.Currency)
}

func (r *Rate) GetForwardRate() float64 {
	return math.Round((float64(r.DestAmt.Amount)/float64(r.SrcAmt.Amount))*100) / 100
}

func (r *Rate) CalculateForwardAmount(amount float64) float64 {
	return math.Round((float64(r.DestAmt.Amount)*amount/float64(r.SrcAmt.Amount))*100) / 100
}

func (r *Rate) GetBackwardRate() float64 {
	return math.Round((float64(r.SrcAmt.Amount)/float64(r.DestAmt.Amount))*100) / 100
}

func (r *Rate) CalculateBackwardAmount(amount float64) float64 {
	return math.Round((float64(r.SrcAmt.Amount)*amount/float64(r.DestAmt.Amount))*100) / 100
}

func GenerateRateId(srcCurrency, destCurrency string) string {
	return fmt.Sprintf("%s-%s", strings.ToUpper(srcCurrency), strings.ToUpper(destCurrency))
}
