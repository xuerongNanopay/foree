package basic_access

import "time"

type BasicAccessStatus string

const (
	StatusWaitingVerify BasicAccessStatus = "WAITING_VERIFY"
	StatusPassExpire    BasicAccessStatus = "PASSWORD_EXPIRE"
	StatusActive        BasicAccessStatus = "ACTIVE"
	StatusSuspend       BasicAccessStatus = "SUSPEND"
	StatusDisable       BasicAccessStatus = "DISABLE"
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
