package promotion

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/types"
)

const (
	sQLPromotionGetUniqueByName = `
		SELECT
			g.name, g.description, g.amount, g.currency,
			g.is_enable, g.start_time, g.end_time, g.created_at, g.updated_at
		FROM promotion as g
		Where g.name = ?
	`
)

// Control the life cycle of promotion.
type Promotion struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Amt         types.AmountData `json:"Amt"`
	IsEnable    bool             `json:"isEnable"`
	StartTime   *time.Time       `json:"startTime"`
	EndTime     *time.Time       `json:"endTime"`
	CreatedAt   *time.Time       `json:"createdAt"`
	UpdatedAt   *time.Time       `json:"updatedAt"`
}

func (p *Promotion) IsValid() bool {
	if !p.IsEnable {
		return false
	}

	if p.StartTime == nil && p.EndTime == nil {
		return true
	}

	now := time.Now()

	if p.StartTime == nil && now.Before(*p.EndTime) {
		return true
	}

	if p.EndTime == nil && now.After(*p.StartTime) {
		return true
	}

	if now.After(*p.StartTime) && now.Before(*p.EndTime) {
		return true
	}

	return false
}

func NewPromotionRepo(db *sql.DB) *PromotionRepo {
	return &PromotionRepo{db: db}
}

type PromotionRepo struct {
	db *sql.DB
}

func (repo *PromotionRepo) GetUniquePromotionByName(ctx context.Context, name string) (*Promotion, error) {
	rows, err := repo.db.Query(sQLPromotionGetUniqueByName, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Promotion

	for rows.Next() {
		f, err = scanRowIntoPromotion(rows)
		if err != nil {
			return nil, err
		}
	}

	if f == nil || f.Name == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoPromotion(rows *sql.Rows) (*Promotion, error) {
	p := new(Promotion)
	err := rows.Scan(
		&p.Name,
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
