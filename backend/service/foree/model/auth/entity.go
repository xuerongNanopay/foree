package auth

import "time"

type AuthenticatorStatus string

const (
	Unverify AuthenticatorStatus = "unverify"
	Active   AuthenticatorStatus = "active"
	Suspend  AuthenticatorStatus = "suspend"
	Disable  AuthenticatorStatus = "disable"
)

type Authenticator struct {
	id uint64

	verifiedAt time.Timer
	createAt   time.Timer
	updateAt   time.Timer
}
