package auth

import "time"

type SessionService interface {
	HasPermission(session Session, permission string) (bool, error)
}

type Session struct {
	ID          uint64       `json:"id"`
	UserId      uint64       `json:"userId"`
	User        User         `json:"user"`
	Permissions []Permission `json:"permission"`
	UserAgent   string       `json:"userAgent"`
	Ip          string       `json:"op"`
	ExpireAt    time.Time
}
