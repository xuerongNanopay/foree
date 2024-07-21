package user

import "time"

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

	createAt time.Timer
	updateAt time.Timer
}
