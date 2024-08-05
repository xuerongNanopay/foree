package auth

import "time"

type AccessStatus string

const (
	ASStatusWaitingVerify AccessStatus = "WAITING_VERIFY"
	ASStatusPassExpire    AccessStatus = "PASSWORD_EXPIRE"
	ASStatusActive        AccessStatus = "ACTIVE"
	ASStatusSuspend       AccessStatus = "SUSPEND"
	ASStatusDisable       AccessStatus = "DISABLE"
)

type EmailPasswdAccess struct {
	ID             uint64       `json:"id"`
	Status         AccessStatus `json:"status"`
	Email          string       `json:"email"`
	Passowrd       string       `json:"-"`
	VerifyCode     string       `json:"-"`
	CodeVerifiedAt time.Time    `json:"codeVerifiedAt"`
	CreateAt       time.Timer   `json:"createAt"`
	UpdateAt       time.Timer   `json:"updateAt"`
	UserId         uint64       `json:"userId"`
}
