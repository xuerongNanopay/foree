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
	LatestActiveAt time.Time    `json:"latestActiveAt"`
	ExpireAt       time.Time    `json:"expireAt"`
	CreateAt       time.Time    `json:"createAt"`
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{mem: make(map[string]*Session, 1024)}
}

// TODO: Thread Safe.
// TODO: Improve
// TODO: use bucket to distribution in to different map.
type SessionRepo struct {
	// db *sql.DB

	mem    map[string]*Session
	rwLock sync.RWMutex
}

func (repo *SessionRepo) Insert(session *Session) (string, error) {
	session.ID = generateSessionId(0)

	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	repo.mem[session.ID] = session
	return session.ID, nil
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
	s.LatestActiveAt = time.Now()
	return s
}

func generateSessionId(bucketId int) string {
	return fmt.Sprintf("%06d-%s", bucketId, uuid.New().String())
}
