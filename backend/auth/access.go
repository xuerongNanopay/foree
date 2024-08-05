package auth

import "time"

type AccessStatus string

const (
	AccessStatusWaitingVerify AccessStatus = "WAITING_VERIFY"
	AccessStatusPassExpire    AccessStatus = "PASSWORD_EXPIRE"
	AccessStatusActive        AccessStatus = "ACTIVE"
	AccessStatusSuspend       AccessStatus = "SUSPEND"
	AccessStatusDisable       AccessStatus = "DISABLE"
)

type EmailPasswdAccess struct {
	ID             uint64       `json:"id"`
	Status         AccessStatus `json:"status"`
	Email          string       `json:"email"`
	Passowrd       string       `json:"-"`
	VerifyCode     string       `json:"-"`
	CodeVerifiedAt time.Time    `json:"codeVerifiedAt"`
	CreateAt       time.Time    `json:"createAt"`
	UpdateAt       time.Time    `json:"updateAt"`
	UserId         uint64       `json:"userId"`
}
