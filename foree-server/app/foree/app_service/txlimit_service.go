package foree_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/auth"
)

const txLimitCacheExpiry time.Duration = 4 * time.Minute
const txLimitCacheRefreshInterval time.Duration = 2 * time.Minute

func NewTxLimitService(
	txLimitRepo *transaction.TxLimitRepo,
	dailyTxLimiteRepo *transaction.DailyTxLimitRepo,
) *TxLimitService {
	txLimitService := &TxLimitService{
		txLimitRepo:               txLimitRepo,
		dailyTxLimiteRepo:         dailyTxLimiteRepo,
		txlimitCacheInsertChan:    make(chan string, 1),
		txlimitCacheUpdateChan:    make(chan string, 1),
		txlimitCacheRefreshTicker: time.NewTicker(txLimitCacheRefreshInterval),
	}
	txLimitService.start()
	return txLimitService
}

type TxLimitService struct {
	txLimitRepo               *transaction.TxLimitRepo
	cache                     sync.Map
	dailyTxLimiteRepo         *transaction.DailyTxLimitRepo
	txlimitCacheInsertChan    chan string
	txlimitCacheUpdateChan    chan string
	txlimitCacheRefreshTicker *time.Ticker
}

func (t *TxLimitService) start() {
	for {
		select {
		case limitGroup := <-t.txlimitCacheInsertChan:
			txLimit, err := t.txLimitRepo.GetUniqueTxLimitByLimitGroup(limitGroup)
			if err != nil {
				foree_logger.Logger.Error("TxLimit_Cache_Insert_Fail", "limitGroup", limitGroup, "cause", err.Error())
			} else {
				t.cache.Store(limitGroup, CacheItem[transaction.TxLimit]{
					item:      *txLimit,
					expiredAt: time.Now().Add(rateCacheExpiry),
				})
			}
		case limitGroup := <-t.txlimitCacheUpdateChan:
			txLimit, err := t.txLimitRepo.GetUniqueTxLimitByLimitGroup(limitGroup)
			if err != nil {
				foree_logger.Logger.Error("TxLimit_Cache_Update_Fail", "limitGroup", limitGroup, "cause", err.Error())
			} else {
				t.cache.Swap(limitGroup, CacheItem[transaction.TxLimit]{
					item:      *txLimit,
					expiredAt: time.Now().Add(rateCacheExpiry),
				})
			}

		case <-t.txlimitCacheRefreshTicker.C:
			length := 0
			t.cache.Range(func(k, _ interface{}) bool {
				length += 1
				limitGroup, _ := k.(string)
				txLimit, err := t.txLimitRepo.GetUniqueTxLimitByLimitGroup(limitGroup)
				if err != nil {
					foree_logger.Logger.Error("TxLimit_Cache_Refresh_Fail", "limitGroup", limitGroup, "cause", err.Error())
				} else {
					t.cache.Swap(limitGroup, CacheItem[transaction.TxLimit]{
						item:      *txLimit,
						expiredAt: time.Now().Add(txLimitCacheExpiry),
					})
				}
				return true
			})
			if length > 32 {
				foree_logger.Logger.Error("TxLimit_Cache_Refresh", "message", "Size of rate cache is greater than 32, please check if txLimitCacheExpiry and txLimitCacheRefreshInterval are still suitable with this cache size.")
			}
		}
	}
}

func (t *TxLimitService) getTxLimit(ctx context.Context, limitGroup string) (*transaction.TxLimit, error) {
	value, ok := t.cache.Load(limitGroup)
	if !ok {
		txLimit, err := t.txLimitRepo.GetUniqueTxLimitByLimitGroup(limitGroup)
		if err != nil {
			return nil, err
		}
		select {
		case t.txlimitCacheInsertChan <- limitGroup:
		default:
		}
		return txLimit, nil
	}
	cacheItem, _ := value.(CacheItem[[]*transaction.Fee])
	if cacheItem.expiredAt.Before(time.Now()) {
		select {
		case t.txlimitCacheUpdateChan <- limitGroup:
		default:
		}
	}
	return t.txLimitRepo.GetUniqueTxLimitByLimitGroup(limitGroup)
}

func (t *TxLimitService) addDailyTxLimit(ctx context.Context, session auth.Session, amt types.AmountData) (*transaction.DailyTxLimit, error) {
	dailyLimit, err := t.getDailyTxLimit(ctx, session)
	if err != nil {
		return nil, err
	}

	dailyLimit.UsedAmt.Amount += amt.Amount

	if err := t.dailyTxLimiteRepo.UpdateDailyTxLimitById(ctx, *dailyLimit); err != nil {
		return nil, err
	}

	return dailyLimit, nil
}

// I don't case race condition here, cause create transaction will rescure it.
func (t *TxLimitService) getDailyTxLimit(ctx context.Context, session auth.Session) (*transaction.DailyTxLimit, error) {
	reference := transaction.GenerateDailyTxLimitReference(session.UserId)
	dailyLimit, err := t.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
	if err != nil {
		return nil, err
	}

	// If not create one.
	if dailyLimit == nil {
		txLimit, err := t.getTxLimit(ctx, session.UserGroup.TransactionLimitGroup)
		if err != nil {
			return nil, err
		}
		if txLimit == nil {
			return nil, fmt.Errorf("transaction limit no found for group `%v`", session.UserGroup.TransactionLimitGroup)
		}

		dailyLimit = &transaction.DailyTxLimit{
			Reference: reference,
			UsedAmt: types.AmountData{
				Amount:   0.0,
				Currency: txLimit.MaxAmt.Currency,
			},
			MaxAmt: types.AmountData{
				Amount:   txLimit.MaxAmt.Amount,
				Currency: txLimit.MaxAmt.Currency,
			},
			OwnerId: session.UserId,
		}
		_, err = t.dailyTxLimiteRepo.InsertDailyTxLimit(ctx, *dailyLimit)
		if err != nil {
			return nil, err
		}
		dl, err := t.dailyTxLimiteRepo.GetUniqueDailyTxLimitByReference(ctx, reference)
		if err != nil {
			return nil, err
		}
		dailyLimit = dl
	}
	return dailyLimit, nil
}
