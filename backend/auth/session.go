package auth

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// How do we store the session. // Redis?
// type SessionService interface {
// 	HasPermission(session Session, permission string) (bool, error)
// }

type Session struct {
	ID             string       `json:"id"`
	UserId         uint64       `json:"userId"`
	User           User         `json:"user"`
	Permissions    []Permission `json:"permission"`
	UserAgent      string       `json:"userAgent"`
	Ip             string       `json:"ip"`
	LatestActiveAt time.Time    `json:"latest_active_at"`
	ExpireAt       time.Time    `json:"expire_at"`
	CreateAt       time.Time    `json:"createAt"`
	UpdateAt       time.Time    `json:"updateAt"`
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{mem: make(map[string]*Session, 1024)}
}

// TODO: Thread Safe.
// TODO: Improve
type SessionRepo struct {
	// db *sql.DB
	mem    map[string]*Session
	rwLock *sync.RWMutex
}

func (repo *SessionRepo) Insert(session *Session) (*Session, error) {
	sessionId := fmt.Sprintf("%v::%v", session.UserId, uuid.New().String())
	session.ID = sessionId
	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	repo.mem[sessionId] = session
	return nil, nil
}

func (repo *SessionRepo) Delete(id string) {

	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	delete(repo.mem, id)
}

func (repo *SessionRepo) GetUniqueById(id string) *Session {
	repo.rwLock.RLock()
	defer repo.rwLock.RUnlock()
	s, ok := repo.mem[id]
	if !ok {
		return nil
	}
	return s
}
