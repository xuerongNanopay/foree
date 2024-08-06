package account

import "time"

type InteracAccount struct {
	ID         int64
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	OwnerId    int64
	Status     AccountStatus
	CreateAt   time.Time `json:"createAt"`
	UpdateAt   time.Time `json:"updateAt"`
}
