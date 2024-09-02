package account

type AccountStatus string

const (
	AccountStatusInitial AccountStatus = "INITIAL"
	AccountStatusActive  AccountStatus = "ACTIVE"
	AccountStatusSuspend AccountStatus = "SUSPEND"
	AccountStatusDisable AccountStatus = "DISABLE"
	AccountStatusDelete  AccountStatus = "DELETE"
)
