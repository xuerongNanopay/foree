package account

type InteracAccount struct {
	ID         int64
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	OwnerId    int64
	Status     AccountStatus
}
