package auth

import "time"

type UserIdentification struct {
	ID       int64     `json:"id"`
	Type     string    `json:"type"`
	Value    string    `json:"value"`
	OwnerId  int64     `json:"ownerId"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}
