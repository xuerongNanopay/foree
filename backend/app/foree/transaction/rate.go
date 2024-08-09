package transaction

import (
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
			r.create_at, r.update_at
        FROM rate r
        where t.id = ?
	`
)

// ID format: {src_currency}-{dest_currency}
type Rate struct {
	ID       string           `json:"id"`
	SrcAmt   types.AmountData `json:"srcAmt"`
	DestAmt  types.AmountData `json:"destAmt"`
	CreateAt time.Time        `json:"createAt"`
	UpdateAt time.Time        `json:"updateAt"`
}

func NewRateRepo(db *sql.DB) *RateRepo {
	return &RateRepo{db: db}
}

type RateRepo struct {
	db *sql.DB
}

func (repo *ForeeTxRepo) InsertRate(r Rate) (string, error) {
	result, err := repo.db.Exec(
		sqlRateInsert,
		r.GetId(),
		r.SrcAmt.Amount,
		strings.ToUpper(r.SrcAmt.Curreny),
		r.GetForwardRate(),
		strings.ToUpper(r.DestAmt.Curreny),
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

func (repo *ForeeTxRepo) UpdateRateById(r Rate) error {
	_, err := repo.db.Exec(sQLRateUpdateById, r.SrcAmt.Amount, r.GetForwardRate(), r.GetId())
	if err != nil {
		return err
	}
	return nil
}

func (repo *ForeeTxRepo) GetUniqueRateById(id string) (*Rate, error) {
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

	if f.ID == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoRate(rows *sql.Rows) (*Rate, error) {
	u := new(Rate)
	err := rows.Scan(
		&u.ID,
		&u.SrcAmt.Amount,
		&u.SrcAmt.Curreny,
		&u.DestAmt.Amount,
		&u.DestAmt.Curreny,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Rate) GetId() string {
	return GenerateRateId(r.SrcAmt.Curreny, r.DestAmt.Curreny)
}

func (r *Rate) ToSummary() string {
	return fmt.Sprintf("$%.2f %s : %.2f %s", r.SrcAmt.Amount, r.SrcAmt.Curreny, r.DestAmt.Amount, r.DestAmt.Curreny)
}

func (r *Rate) GetForwardRate() float64 {
	return math.Round((float64(r.DestAmt.Amount)/float64(r.SrcAmt.Amount))*100) / 100
}

func (r *Rate) GetBackwardRate() float64 {
	return math.Round((float64(r.SrcAmt.Amount)/float64(r.DestAmt.Amount))*100) / 100
}

func GenerateRateId(srcCurrency, destCurrency string) string {
	return fmt.Sprintf("%s-%s", strings.ToUpper(srcCurrency), strings.ToUpper(destCurrency))
}
