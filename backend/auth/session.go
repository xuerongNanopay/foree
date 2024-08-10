package auth

import (
	"fmt"
	"strconv"
	"strings"
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
	UserId         int64        `json:"userId"`
	EmailPasswd    *EmailPasswd `json:"emailPasswd"`
	User           *User        `json:"user"`
	Permissions    []Permission `json:"permission"`
	UserAgent      string       `json:"userAgent"`
	Ip             string       `json:"ip"`
	LatestActiveAt time.Time    `json:"latestActiveAt"`
	ExpireAt       time.Time    `json:"expireAt"`
	CreateAt       time.Time    `json:"createAt"`
}

// 13 buckets, 1024 sesson of each bucket, and 12 hours session expiry
// It should be able to support at least 10000 sessions in 12 hours
func NewDefaultSessionRepo() *SessionRepo {
	// If you don't have active in 3 hours, the session will expire.
	return NewSessionRepo(12, 3, 13, 1024)
}

func NewSessionRepo(expireInHour, activeInHour, numberOfBucket, maxBucketSize int) *SessionRepo {

	mems := []map[string]*Session{}
	for i := 0; i < maxBucketSize; i++ {
		mems = append(mems, make(map[string]*Session, maxBucketSize/4))
	}
	return &SessionRepo{
		mems:           mems,
		cur:            0,
		activeInHour:   activeInHour,
		expireInHour:   expireInHour,
		numberOfBucket: numberOfBucket,
		maxBucketSize:  maxBucketSize,
	}
}

// Still have performance issue. TOOD: use atomic instead of lock
type SessionRepo struct {
	// db *sql.DB
	cur            int
	maxBucketSize  int
	expireInHour   int
	activeInHour   int
	numberOfBucket int
	mems           []map[string]*Session
	lock           sync.Mutex
}

func (repo *SessionRepo) InsertSession(session Session) (string, error) {
	session.CreateAt = time.Now()
	session.LatestActiveAt = time.Now()
	session.ExpireAt = time.Now().Add(time.Duration(time.Hour * time.Duration(repo.expireInHour)))

	repo.lock.Lock()
	defer repo.lock.Unlock()

	session.ID = generateSessionId(repo.cur)
	if len(repo.mems[repo.cur%repo.numberOfBucket]) > repo.maxBucketSize {
		go repo.purge(repo.cur)
		repo.cur += 1
		if len(repo.mems[repo.cur%repo.numberOfBucket]) != 0 {
			return "", fmt.Errorf("sesson pool is full")
		}
	}
	repo.mems[repo.cur%repo.numberOfBucket][session.ID] = &session
	return session.ID, nil
}

func (repo *SessionRepo) UpdateSession(session Session) (*Session, error) {
	session.LatestActiveAt = time.Now()
	idx, err := parseBucketId(session.ID)
	if err != nil {
		return nil, err
	}

	repo.lock.Lock()
	defer repo.lock.Unlock()

	repo.mems[idx%repo.numberOfBucket][session.ID] = &session

	return repo.mems[idx%repo.numberOfBucket][session.ID], nil
}

func (repo *SessionRepo) purge(bucketIdx int) {
	//Sleep 2 * Expiry, make sure all quote in the bucket are expiry.
	time.Sleep(time.Hour * time.Duration(repo.expireInHour/2))
	//TODO: Log
	//Clear all quote by just replace with new map
	repo.lock.Lock()
	defer repo.lock.Unlock()
	repo.mems[bucketIdx%repo.numberOfBucket] = make(map[string]*Session, repo.maxBucketSize/4)
}

func (repo *SessionRepo) Delete(id string) {
	idx, err := parseBucketId(id)
	if err != nil {
		return
	}
	repo.lock.Lock()
	defer repo.lock.Unlock()
	delete(repo.mems[idx%repo.numberOfBucket], id)
}

func (repo *SessionRepo) GetSessionUniqueById(id string) *Session {
	idx, err := parseBucketId(id)
	if err != nil {
		return nil
	}
	mem := repo.mems[idx%repo.numberOfBucket]
	if mem == nil {
		return nil
	}

	s, ok := mem[id]
	if !ok {
		return nil
	}
	now := time.Now()

	if now.Unix()-s.LatestActiveAt.Unix() > int64(repo.activeInHour*3600) {
		go repo.Delete(id)
		return nil
	}
	if now.Unix() > s.ExpireAt.Unix() {
		go repo.Delete(id)
		return nil
	}
	s.LatestActiveAt = now
	return s
}

func generateSessionId(bucketId int) string {
	return fmt.Sprintf("%06d-%s", bucketId, uuid.New().String())
}

func parseBucketId(sessionId string) (int, error) {
	s := strings.Split(sessionId, "-")
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, err
	}
	return i, nil
}
