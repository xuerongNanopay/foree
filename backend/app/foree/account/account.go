package account

type AccountStatus string

const (
	AccountStatusInitial = "INITIAL"
	AccountStatusActive  = "ACTIVE"
	AccountStatusSuspend = "SUSPEND"
	AccountStatusDisable = "DISABLE"
	AccountStatusDelete  = "DELETE"
)
