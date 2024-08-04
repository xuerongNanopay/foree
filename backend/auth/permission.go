package auth

import "time"

// Super: *::*::*
// app::service::methods
type Group struct {
	ID          string       `json:"id"`
	Description string       `json:"description"`
	IsEnable    bool         `json:"isEnable"`
	CreateAt    time.Timer   `json:"createAt"`
	UpdateAt    time.Timer   `json:"updateAt"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	IsEnable    bool       `json:"isEnable"`
	CreateAt    time.Timer `json:"createAt"`
	UpdateAt    time.Timer `json:"updateAt"`
}

type GroupPermissionJoin struct {
	GroupID      string     `json:"groupId"`
	PermissionID string     `json:"permissionId"`
	IsEnable     bool       `json:"isEnable"`
	CreateAt     time.Timer `json:"createAt"`
	UpdateAt     time.Timer `json:"updateAt"`
}
