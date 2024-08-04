package user

import "time"

type UserStatus string

const (
	PENDING UserStatus = "PENDING"
	ACTIVE  UserStatus = "ACTIVE"
	SUSPEND UserStatus = "SUSPEND"
	DISABLE UserStatus = "DISABLE"
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
