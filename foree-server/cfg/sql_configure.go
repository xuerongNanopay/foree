package cfg

import (
	"database/sql"
	"time"
)

// Get configure from mysql server.

type configuration struct {
	Name            string
	RawValue        string
	RefreshInterval int64
	CreatedAt       *time.Time `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
}

type SQLConfigure struct {
	db      *sql.DB
	configs map[string]Config[any]
}

func scanRowIntoConfiguration(rows *sql.Rows) (*configuration, error) {
	c := new(configuration)
	err := rows.Scan(
		&c.Name,
		&c.RawValue,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}
