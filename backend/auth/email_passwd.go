package auth

import "time"

type EmailPasswdStatus string

const (
	EPStatusWaitingVerify EmailPasswdStatus = "WAITING_VERIFY"
	EPStatusPassExpire    EmailPasswdStatus = "PASSWORD_EXPIRE"
	EPStatusActive        EmailPasswdStatus = "ACTIVE"
	EPStatusSuspend       EmailPasswdStatus = "SUSPEND"
	EPStatusDisable       EmailPasswdStatus = "DISABLE"
)

type EmailPasswd struct {
	ID             uint64            `json:"id"`
	Status         EmailPasswdStatus `json:"status"`
	Email          string            `json:"email"`
	Passowrd       string            `json:"-"`
	VerifyCode     string            `json:"-"`
	CodeVerifiedAt time.Time         `json:"codeVerifiedAt"`
	CreateAt       time.Time         `json:"createAt"`
	UpdateAt       time.Time         `json:"updateAt"`
	UserId         uint64            `json:"userId"`
}
