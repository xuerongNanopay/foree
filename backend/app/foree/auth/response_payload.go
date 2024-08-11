package auth

import "xue.io/go-pay/auth"

type UserDTO struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
	Status     string `json:"status"`
	AvatarUrl  string `json:"avatarUrl"`
}

func NewUserDTO(user *auth.User) *UserDTO {
	return &UserDTO{
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Status:     string(user.Status),
		AvatarUrl:  user.AvatarUrl,
	}
}
