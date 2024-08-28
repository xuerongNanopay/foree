package transaction

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

//Store exactly same payload as ForeeTx.
//We contract a initial ForeeTx, and the constomer can review the transaction detail.
// Current we store it in memory. In the futurn, can be redis....

type TxQuote struct {
	ID        string    `json:"id"`
	Tx        *ForeeTx  `json:"tx"`
	OwerId    int64     `json:"owerId"`
	ExpireAt  time.Time `json:"expireAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// set to 512
func NewTxQuoteRepo(expire, maxBucketSize int) *TxQuoteRepo {
	return &TxQuoteRepo{
		mems: [2]map[string]TxQuote{
			make(map[string]TxQuote, maxBucketSize),
			make(map[string]TxQuote, maxBucketSize),
		},
		cur:            0,
		expireInMinute: expire,
		maxBucketSize:  maxBucketSize,
	}
}

// Still have performance issue
type TxQuoteRepo struct {
	mems           [2]map[string]TxQuote
	cur            int
	maxBucketSize  int
	expireInMinute int
	rwLock         sync.RWMutex
}

func (repo *TxQuoteRepo) InsertTxQuote(ctx context.Context, tx TxQuote) (string, error) {
	tx.CreatedAt = time.Now()
	tx.ExpireAt = time.Now().Add(time.Duration(time.Minute * time.Duration(repo.expireInMinute)))
	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	tx.ID = generateTxQuoteId(repo.cur)
	if len(repo.mems[repo.cur%2]) > repo.maxBucketSize {
		go repo.purge(repo.cur)
		repo.cur += 1
	}
	repo.mems[repo.cur%2][tx.ID] = tx
	return tx.ID, nil
}

func (repo *TxQuoteRepo) purge(bucketIdx int) {
	//Sleep 2 * Expiry, make sure all quote in the bucket are expiry.
	time.Sleep(time.Minute*time.Duration(repo.expireInMinute) + time.Minute)
	//TODO: Log
	//Clear all quote by just replace with new map
	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	repo.mems[bucketIdx%2] = make(map[string]TxQuote, repo.maxBucketSize)
}

func (repo *TxQuoteRepo) Delete(id string) {
	idx, err := parseBucketId(id)
	if err != nil {
		return
	}
	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	delete(repo.mems[idx%2], id)
}

func (repo *TxQuoteRepo) GetUniqueById(ctx context.Context, id string) *TxQuote {
	idx, err := parseBucketId(id)
	if err != nil {
		return nil
	}
	repo.rwLock.RLock()
	defer repo.rwLock.RUnlock()
	s, ok := repo.mems[idx%2][id]
	if !ok {
		return nil
	}
	return &s
}

func generateTxQuoteId(bucketId int) string {
	return fmt.Sprintf("%08d-%s", bucketId, uuid.New().String())
}

func parseBucketId(quoteId string) (int, error) {
	s := strings.Split(quoteId, "-")
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, err
	}
	return i, nil
}
