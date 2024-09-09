package promotion

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLGiftGetUniqueByCode = `
		SELECT
			g.code, g.description, g.amount, g.currency, g,require_min
			g.is_enable, g.start_time, g.end_time, g.created_at, g.updated_at
		FROM gift as g
		Where g.code = ?
	`
)

// Control the life cycle of promotion.
type Gift struct {
	Code        string           `json:"code"`
	Description string           `json:"description"`
	Amt         types.AmountData `json:"Amt"`
	IsEnable    bool             `json:"isEnable"`
	StartTime   time.Time        `json:"startTime"`
	EndTime     time.Time        `json:"endTime"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

func (p *Gift) IsValid() bool {
	if !p.IsEnable {
		return false
	}

	if now := time.Now().Unix(); now > p.StartTime.Unix() || (now > p.EndTime.Unix() && !p.EndTime.IsZero()) {
		return false
	}
	return true
}

func NewGiftRepo(db *sql.DB) *GiftRepo {
	return &GiftRepo{db: db}
}

type GiftRepo struct {
	db *sql.DB
}

func (repo *GiftRepo) GetUniqueGiftByCode(ctx context.Context, code string) (*Gift, error) {
	rows, err := repo.db.Query(sQLGiftGetUniqueByCode, code)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Gift

	for rows.Next() {
		f, err = scanRowIntoGift(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.Code == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoGift(rows *sql.Rows) (*Gift, error) {
	p := new(Gift)
	err := rows.Scan(
		&p.Code,
		&p.Description,
		&p.Amt.Amount,
		&p.Amt.Currency,
		&p.IsEnable,
		&p.StartTime,
		&p.EndTime,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}
