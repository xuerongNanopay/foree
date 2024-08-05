package auth

import (
	"database/sql"
	"time"
)

// How do we store the session. // Redis?
// type SessionService interface {
// 	HasPermission(session Session, permission string) (bool, error)
// }

type Session struct {
	ID          uint64       `json:"id"`
	UserId      uint64       `json:"userId"`
	User        User         `json:"user"`
	Permissions []Permission `json:"permission"`
	UserAgent   string       `json:"userAgent"`
	Ip          string       `json:"ip"`
	ExpireAt    time.Time    `json:"expire_at"`
	CreateAt    time.Time    `json:"createAt"`
	UpdateAt    time.Time    `json:"updateAt"`
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{mem: make(map[string]*Session, 1024)}
}

// TODO: Thread Safe.
type SessionRepo struct {
	// db *sql.DB
	mem map[string]*Session
}

func (repo *SessionRepo) insert(session *Session) (*Session, error) {
	// repo.mem
	return nil, nil
}

func (repo *SessionRepo) delete(id string) error {
	// repo.mem
	return nil
}

func (repo *SessionRepo) GetById(id string) (*Session, error) {
	return nil, nil
}
