package auth

import "time"

type AuthenticatorStatus uint

const (
	Unverify AuthenticatorStatus = iota + 1
	Active
	Suspend
	Disable
)

type Authenticator struct {
	id uint64

	verifiedAt time.Timer
	createAt   time.Timer
	updateAt   time.Timer
}
