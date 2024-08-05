package auth

import (
	"database/sql"
	"time"
)

const (
	SQLPermissionByGroupId = `
		SELECT 
			p.id
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

func (repo *PermissionRepo) GetAllByGroupId(groupId string) ([]Permission, error) {
	return nil, nil
}
