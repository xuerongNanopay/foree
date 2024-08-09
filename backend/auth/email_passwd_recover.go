package auth

import "time"

type EmailPasswdRecover struct {
	ID            string    `json:"id"`
	EmailPasswdId int16     `json:"emailPasswdId"`
	ExpiryAt      time.Time `json:"expiryAt"`
	CreateAt      time.Time `json:"createAt"`
	UpdateAt      time.Time `json:"updateAt"`
}
