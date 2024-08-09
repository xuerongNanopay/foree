package auth

import (
	"database/sql"
	"strings"
	"time"
)

const (
	sQLPermissionByGroupName = `
		SELECT 
			p.name, p.description
		FROM permissions as p
		INNER JOIN group_permission_joint as gpt ON p.name = gpt.permission_name
		WHERE gpt.is_enable = true and gpt.group_name = ?
	`
)

// ReadOnly
// Super: *::*::*
// app::service::methods
type Group struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreateAt    time.Time    `json:"createAt"`
	UpdateAt    time.Time    `json:"updateAt"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"createAt"`
	UpdateAt    time.Time `json:"updateAt"`
}

type GroupPermissionJoint struct {
	GroupName      string    `json:"groupName"`
	PermissionName string    `json:"permissionName"`
	IsEnable       bool      `json:"isEnable"`
	CreateAt       time.Time `json:"createAt"`
	UpdateAt       time.Time `json:"updateAt"`
}

func NewPermission(db *sql.DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

type PermissionRepo struct {
	db *sql.DB
}

func (repo *PermissionRepo) GetAllPermissionByGroupName(groupName string) ([]*Permission, error) {
	rows, err := repo.db.Query(sQLPermissionByGroupName)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ps := make([]*Permission, 16)
	for rows.Next() {
		p, err := scanRowIntoPermission(rows)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ps, nil
}

func scanRowIntoPermission(rows *sql.Rows) (*Permission, error) {
	p := new(Permission)
	err := rows.Scan(
		&p.Name,
		&p.Description,
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
