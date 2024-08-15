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
			p.Name, p.quantity, p.amount, p.currency,
			p.start_time, p.end_time, p.is_enable,
			p.create_at, p.update_at
		WHERE p.name = ?
	`
	sQLPromotionUpdateByName = `
		SELECT 
			p.quantity = ?
		WHERE p.name = ?
	`
)

type Promotion struct {
	Name      string           `json:"id"`
	Quantity  int32            `json:"limit"`
	Amt       types.AmountData `json:"Amt"`
	StartTime time.Time        `json:"startTime"`
	EndTime   time.Time        `json:"endTime"`
	IsEnable  bool             `json:"isEnable"`
	CreateAt  time.Time        `json:"createAt"`
	UpdateAt  time.Time        `json:"updateAt"`
}

func (p *Promotion) CanApply() bool {
	if !p.IsEnable {
		return false
	}

	if p.Quantity <= 0 {
		return false
	}

	now := time.Now().Unix()

	if now > p.StartTime.Unix() || (now > p.EndTime.Unix() && !p.EndTime.IsZero()) {
		return false
	}

	return true
}

func NewPromotionRepo(db *sql.DB) *PromotionRepo {
	return &PromotionRepo{db: db}
}

type PromotionRepo struct {
	db *sql.DB
}

func (repo *PromotionRepo) UpdatePromotionByName(ctx context.Context, p Promotion) error {
	_, err := repo.db.Exec(
		sQLPromotionUpdateByName,
		p.Quantity,
		p.Name,
	)
	if err != nil {
		return err
	}
	return nil
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

	if f.Name == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoPromotion(rows *sql.Rows) (*Promotion, error) {
	u := new(Promotion)
	err := rows.Scan(
		&u.Name,
		&u.Quantity,
		&u.Amt.Amount,
		&u.Amt.Curreny,
		&u.StartTime,
		&u.EndTime,
		&u.IsEnable,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// func (p *Promotion) CanApply(name string) bool {
// 	if name != p.Name {
// 		return false
// 	}

// 	if atomic.LoadInt32(&p.Quantity) == 0 {
// 		return false
// 	}

// 	atomic.AddInt32(&p.Quantity, -1)

// 	now := time.Now().Unix()

// 	if now > p.StartTime.Unix() || (now > p.EndTime.Unix() && !p.EndTime.IsZero()) {
// 		return false
// 	}

// 	return true
// }
