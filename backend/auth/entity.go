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

type UserStatus string

const (
	UserStatusInitial UserStatus = "INITIAL"
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusSuspend UserStatus = "SUSPEND"
	UserStatusDisable UserStatus = "DISABLE"
)

type User struct {
	ID          uint64     `json:"id"`
	Group       string     `json:"-"`
	Status      UserStatus `json:"status"`
	FirstName   string     `json:"firstName"`
	MiddleName  string     `json:"middleName"`
	LastName    string     `json:"lastName"`
	Age         int        `json:"age"`
	Dob         time.Time  `json:"dob"`
	Nationality string     `json:"nationality"`
	Address1    string     `json:"address1"`
	Address2    string     `json:"address2"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	Country     string     `json:"country"`
	PhoneNumber string     `json:"phoneNumber"`
	Email       string     `json:"email"`
	CreateAt    time.Timer `json:"createAt"`
	UpdateAt    time.Timer `json:"updateAt"`
	// OccupationId int64      `json:"-"`
	// Occupation   string     `json:"occupation"`
	AvatarUrl string `json:"avatarUrl"`
}

type Session struct {
	ID       uint64 `json:"id"`
	UserId   uint64 `json:"userId"`
	User     User   `json:"user"`
	ExpireAt time.Time
}

type BasicAccessStatus string

const (
	BAStatusWaitingVerify BasicAccessStatus = "WAITING_VERIFY"
	BAStatusPassExpire    BasicAccessStatus = "PASSWORD_EXPIRE"
	BAStatusActive        BasicAccessStatus = "ACTIVE"
	BAStatusSuspend       BasicAccessStatus = "SUSPEND"
	BAStatusDisable       BasicAccessStatus = "DISABLE"
)

type BasicAccess struct {
	ID             uint64            `json:"id"`
	Status         BasicAccessStatus `json:"status"`
	Email          string            `json:"email"`
	Passowrd       string            `json:"-"`
	VerifyCode     string            `json:"-"`
	CodeVerifiedAt time.Time         `json:"codeVerifiedAt"`
	CreateAt       time.Timer        `json:"createAt"`
	UpdateAt       time.Timer        `json:"updateAt"`
	UserId         uint64            `json:"userId"`
}
