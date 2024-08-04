package auth

import "time"

type UserStatus string

const (
	UserStatusInitial UserStatus = "INITIAL"
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusSuspend UserStatus = "SUSPEND"
	UserStatusDisable UserStatus = "DISABLE"
)

type User struct {
	ID          uint64     `json:"id"`
	Group       string     `json:"-"`
	Status      UserStatus `json:"status"`
	FirstName   string     `json:"firstName"`
	MiddleName  string     `json:"middleName"`
	LastName    string     `json:"lastName"`
	Age         int        `json:"age"`
	Dob         time.Time  `json:"dob"`
	Nationality string     `json:"nationality"`
	Address1    string     `json:"address1"`
	Address2    string     `json:"address2"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	Country     string     `json:"country"`
	PhoneNumber string     `json:"phoneNumber"`
	Email       string     `json:"email"`
	CreateAt    time.Timer `json:"createAt"`
	UpdateAt    time.Timer `json:"updateAt"`
	// OccupationId int64      `json:"-"`
	// Occupation   string     `json:"occupation"`
	AvatarUrl string `json:"avatarUrl"`
}
