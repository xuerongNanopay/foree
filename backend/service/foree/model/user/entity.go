package user

import "time"

type UserStatus string

const (
	Pending UserStatus = "pending"
	Active  UserStatus = "active"
	Suspend UserStatus = "suspend"
	Disable UserStatus = "disable"
)

type User struct {
	id uint64

	firstname  string
	middlename string
	lastname   string
	age        uint8
	dob        time.Time

	address1 string
	address2 string
	city     string
	province string
	country  string

	phoneNumber string
	email       string

	group string

	onboardAt time.Timer
	createAt  time.Timer
	updateAt  time.Timer
}
