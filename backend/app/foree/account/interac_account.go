package foree_account

type InteracAccount struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	OwnerId   int64
	Status    AccountStatus
}
