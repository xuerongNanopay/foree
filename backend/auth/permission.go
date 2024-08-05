package auth

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	SQLPermissionByGroupId = `
		SELECT 
			p.id, p.description, p.is_enable
		FROM group_permission as gp
		INNER JOIN permission as p ON gp.permission_id=p.id
		WHERE p.is_enable = true and pg.is_enable = true and pg.group_id = ?
	`
)

// ReadOnly
// Super: *::*::*
// app::service::methods
type Group struct {
	ID          string       `json:"id"`
	Description string       `json:"description"`
	IsEnable    bool         `json:"isEnable"`
	CreateAt    time.Time    `json:"createAt"`
	UpdateAt    time.Time    `json:"updateAt"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	IsEnable    bool      `json:"isEnable"`
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`
}

type GroupPermissionJoin struct {
	GroupID      string    `json:"groupId"`
	PermissionID string    `json:"permissionId"`
	IsEnable     bool      `json:"isEnable"`
	CreateAt     time.Time `json:"createAt"`
	UpdateAt     time.Time `json:"updateAt"`
}

func NewPermission(db *sql.DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

type PermissionRepo struct {
	db *sql.DB
}

func (repo *PermissionRepo) GetAllByGroupId(groupId string) ([]*Permission, error) {
	rows, err := repo.db.Query(SQLPermissionByGroupId)

	if err != nil {
		return nil, fmt.Errorf("GetAllByGroupId: %v", err)
	}
	defer rows.Close()

	ps := make([]*Permission, 16)
	for rows.Next() {
		p, err := scanRowIntoPermission(rows)
		if err != nil {
			return nil, fmt.Errorf("GetAllByGroupId: %v", err)
		}
		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAllByGroupId: %v", err)
	}

	return ps, nil
}

func scanRowIntoPermission(rows *sql.Rows) (*Permission, error) {
	p := new(Permission)
	err := rows.Scan(
		&p.ID,
		&p.Description,
		&p.IsEnable,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func IsPermissionGrand(requiredPermission string, ownedPermission string) bool {
	if requiredPermission == "" || ownedPermission == "" {
		return false
	}
	if requiredPermission == ownedPermission {
		return true
	}

	r := strings.Split(requiredPermission, "::")
	o := strings.Split(requiredPermission, "::")

	if len(r) != len(o) {
		return false
	}

	for i := 0; i < len(r); i++ {
		if r[i] != o[i] && o[i] != "*" {
			return false
		}
	}

	return true
}
