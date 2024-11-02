package cfg

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"xue.io/go-pay/constant"
)

// Get configure from mysql server.
const (
	sQLConfigurationGetUniqueByName = `
		SELECT
			c.name, c.raw_value, c.refresh_interval, 
			r.expire_at, r.created_at, r.updated_at
		FROM configuration as c
		Where c.name = ?
	`
	sQLConfigurationGetAllByNames = `
		SELECT
			c.name, c.raw_value, c.refresh_interval, 
			r.expire_at, r.created_at, r.updated_at
		FROM configuration as c
		Where c.name in (%v)
	`
)

type configuration struct {
	Name            string
	RawValue        string
	RefreshInterval int64
	CreatedAt       *time.Time `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
}

type configurationRepo struct {
	db *sql.DB
}

func (repo *configurationRepo) getAllConfigurationByNames(ctx context.Context, names ...string) ([]*configuration, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	args := make([]interface{}, len(names))
	p := make([]string, len(names))
	for i, n := range names {
		args[i] = n
		p[i] = "?"
	}

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(fmt.Sprintf(sQLConfigurationGetAllByNames, strings.Join(p, ",")), args...)

	} else {
		rows, err = repo.db.Query(fmt.Sprintf(sQLConfigurationGetAllByNames, strings.Join(p, ",")), args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	configures := make([]*configuration, 0)
	for rows.Next() {
		p, err := scanRowIntoConfiguration(rows)
		if err != nil {
			return nil, err
		}
		configures = append(configures, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configures, nil
}

func (repo *configurationRepo) getUniqueConfigurationByName(ctx context.Context, name string) (*configuration, error) {
	dTx, ok := ctx.Value(constant.CKdatabaseTransaction).(*sql.Tx)

	var rows *sql.Rows
	var err error

	if ok {
		rows, err = dTx.Query(sQLConfigurationGetUniqueByName, name)
	} else {
		rows, err = repo.db.Query(sQLConfigurationGetUniqueByName, name)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var c *configuration

	for rows.Next() {
		c, err = scanRowIntoConfiguration(rows)
		if err != nil {
			return nil, err
		}
	}

	if c == nil || c.Name == "" {
		return nil, nil
	}

	return c, nil
}

func scanRowIntoConfiguration(rows *sql.Rows) (*configuration, error) {
	c := new(configuration)
	err := rows.Scan(
		&c.Name,
		&c.RawValue,
		&c.RefreshInterval,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}
