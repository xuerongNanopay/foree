package auth

import "time"

type UserExtra struct {
	ID                 int64     `json:"id"`
	Pob                string    `json:"pob"`
	Cor                string    `json:"cor"`
	Nationality        string    `json:"nationality"`
	OccupationCategory string    `json:"occupationCategory"`
	OccupationName     string    `json:"occupationName"`
	OwnerId            int64     `json:"ownerId"`
	CreateAt           time.Time `json:"createAt"`
	UpdateAt           time.Time `json:"updateAt"`
}
