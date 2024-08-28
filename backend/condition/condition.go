package condition

import (
	"context"
	"database/sql"
	"time"
)

const (
	sQLConditionGetUniqueByName = `
		SELECT 
			p.name, p.limit, p.start_time, p.end_time,
			p.is_enable, p.created_at, p.updated_at
		FROM conditions as p
		WHERE p.name = ?
	`
)

type Condition struct {
	Name      string    `json:"id"`
	Limit     int32     `json:"limit"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	IsEnable  bool      `json:"isEnable"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Condition) Fulfill() bool {
	if !c.IsEnable {
		return false
	}
	if c.Limit <= 0 {
		return false
	}

	now := time.Now().Unix()

	if now > c.StartTime.Unix() || (now > c.EndTime.Unix() && !c.EndTime.IsZero()) {
		return false
	}

	return true
}

func NewConditionRepo(db *sql.DB) *ConditionRepo {
	return &ConditionRepo{db: db}
}

type ConditionRepo struct {
	db *sql.DB
}

func (repo *ConditionRepo) GetUniqueConditionByName(ctx context.Context, name string) (*Condition, error) {
	rows, err := repo.db.Query(sQLConditionGetUniqueByName, name)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *Condition

	for rows.Next() {
		f, err = scanRowIntoCondition(rows)
		if err != nil {
			return nil, err
		}
	}

	if f.Name == "" {
		return nil, nil
	}

	return f, nil
}

func scanRowIntoCondition(rows *sql.Rows) (*Condition, error) {
	u := new(Condition)
	err := rows.Scan(
		&u.Name,
		&u.Limit,
		&u.StartTime,
		&u.EndTime,
		&u.IsEnable,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
