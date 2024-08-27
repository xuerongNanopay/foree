package auth

import (
	"database/sql"
	"strings"
	"time"
)

// ReadOnly
// Super: *::*::*
// app::service::methods
// type Role struct {
// 	Name        string       `json:"name"`
// 	Description string       `json:"description"`
// 	CreatedAt   time.Time    `json:"createdAt"`
// 	UpdatedAt   time.Time    `json:"updatedAt"`
// 	Permissions []Permission `json:"permissions"`
// }

// type Permission struct {
// 	ID          string    `json:"name"`
// 	Description string    `json:"description"`
// 	CreatedAt   time.Time `json:"createdAt"`
// 	UpdatedAt   time.Time `json:"updatedAt"`
// }

const (
	sQLRolePermissionGetAllEnabledByRoleName = `
		SELECT 
			r.role_name, r.permission, r.is_enable, r.created_at, r.update_at
		FROM role_permission as r
		WHERE r.is_enable = true and r.role_name = ?
	`
)

type RolePermission struct {
	RoleName   string    `json:"roleName"`
	Permission string    `json:"permission"`
	IsEnable   bool      `json:"isEnable"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewRolePermission(db *sql.DB) *RolePermissionRepo {
	return &RolePermissionRepo{db: db}
}

type RolePermissionRepo struct {
	db *sql.DB
}

func (repo *RolePermissionRepo) GetAllEnabledRolePermissionByRoleName(roleName string) ([]*RolePermission, error) {
	rows, err := repo.db.Query(sQLRolePermissionGetAllEnabledByRoleName, roleName)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ps := make([]*RolePermission, 16)
	for rows.Next() {
		p, err := scanRowIntoRolePermission(rows)
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

func scanRowIntoRolePermission(rows *sql.Rows) (*RolePermission, error) {
	p := new(RolePermission)
	err := rows.Scan(
		&p.RoleName,
		&p.Permission,
		&p.IsEnable,
		&p.CreatedAt,
		&p.UpdatedAt,
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
