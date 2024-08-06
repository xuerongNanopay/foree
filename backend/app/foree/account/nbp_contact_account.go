package account

type ForeeContactAccount struct {
	ID        int64
	FirstName string
	LastName  string
	OwnerId   int64
	Status    AccountStatus
}
