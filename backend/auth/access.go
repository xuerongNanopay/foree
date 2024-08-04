package auth

import "time"

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
