package transaction

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

//Store exactly same payload as ForeeTx.
//We contract a initial ForeeTx, and the constomer can review the transaction detail.
// Current we store it in memory. In the futurn, can be redis....

type TxQuote struct {
	ID       string    `json:"id"`
	Tx       ForeeTx   `json:"tx"`
	OwerId   int64     `json:"owerId"`
	ExpireAt time.Time `json:"expireAt"`
	CreateAt time.Time `json:"createAt"`
}

func NewTxQuoteRepo() *TxQuoteRepo {
	return &TxQuoteRepo{}
}

type TxQuoteRepo struct {
	mem map[string]*TxQuote
	// mem2    map[string]*TxQuote

	rwLock *sync.RWMutex
}

func (repo *TxQuoteRepo) InsertTxQuote(tx *TxQuote) (string, error) {
	tx.ID = generateTxQuoteId(tx.OwerId)
	tx.CreateAt = time.Now()
	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	repo.mem[tx.ID] = tx
	return tx.ID, nil
}

func (repo *TxQuoteRepo) Delete(id string) {

	repo.rwLock.Lock()
	defer repo.rwLock.Unlock()
	delete(repo.mem, id)
}

func (repo *TxQuoteRepo) GetUniqueById(id string) *TxQuote {
	repo.rwLock.RLock()
	defer repo.rwLock.RUnlock()
	s, ok := repo.mem[id]
	if !ok {
		return nil
	}
	return s
}

func generateTxQuoteId(id int64) string {
	return fmt.Sprintf("%09d-%s", id, uuid.New().String())
}
