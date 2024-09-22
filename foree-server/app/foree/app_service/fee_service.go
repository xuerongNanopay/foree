package foree_service

import (
	"context"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/types"
)

const feeCacheExpiry time.Duration = 4 * time.Minute
const feeCacheRefreshInterval time.Duration = 2 * time.Minute

func NewFeeService(feeRepo *transaction.FeeRepo) *FeeService {
	feeService := &FeeService{
		feeRepo:               feeRepo,
		feeCacheInsertChan:    make(chan string, 1),
		feeCacheUpdateChan:    make(chan string, 1),
		feeCacheRefreshTicker: time.NewTicker(feeCacheRefreshInterval),
	}
	go feeService.start()
	return feeService
}

type FeeService struct {
	feeRepo               *transaction.FeeRepo
	cache                 sync.Map
	feeCacheInsertChan    chan string
	feeCacheUpdateChan    chan string
	feeCacheRefreshTicker *time.Ticker
}

func (r *FeeService) start() {
	for {
		select {
		case feeGroup := <-r.feeCacheInsertChan:
			fees, err := r.feeRepo.GetAllEnableFeeByGroupName(context.TODO(), feeGroup)
			if err != nil {
				foree_logger.Logger.Error("Fee_Cache_Insert_FAIL", "feeGroup", feeGroup, "cause", err.Error())
			} else {
				r.cache.Store(feeGroup, CacheItem[[]*transaction.Fee]{
					item:      fees,
					expiredAt: time.Now().Add(rateCacheExpiry),
				})
			}
		case feeGroup := <-r.feeCacheUpdateChan:
			fees, err := r.feeRepo.GetAllEnableFeeByGroupName(context.TODO(), feeGroup)
			if err != nil {
				foree_logger.Logger.Error("Fee_Cache_Update_FAIL", "feeGroup", feeGroup, "cause", err.Error())
			} else {
				r.cache.Swap(feeGroup, CacheItem[[]*transaction.Fee]{
					item:      fees,
					expiredAt: time.Now().Add(rateCacheExpiry),
				})
			}

		case <-r.feeCacheRefreshTicker.C:
			length := 0
			r.cache.Range(func(k, _ interface{}) bool {
				length += 1
				feeGroup, _ := k.(string)
				fees, err := r.feeRepo.GetAllEnableFeeByGroupName(context.TODO(), feeGroup)
				if err != nil {
					foree_logger.Logger.Error("Fee_Cache_Refresh_FAIL", "feeGroup", feeGroup, "cause", err.Error())
				} else {
					r.cache.Swap(feeGroup, CacheItem[[]*transaction.Fee]{
						item:      fees,
						expiredAt: time.Now().Add(feeCacheExpiry),
					})
				}
				return true
			})
			if length > 32 {
				foree_logger.Logger.Warn("Fee_Cache_Refresh", "message", "Size of rate cache is greater than 64, please check if feeCacheExpiry and feeCacheRefreshInterval are still suitable with this cache size.")
			}
		}
	}
}

func (r *FeeService) getFee(feeGroup string) ([]*transaction.Fee, error) {

	value, ok := r.cache.Load(feeGroup)

	if !ok {
		fees, err := r.feeRepo.GetAllEnableFeeByGroupName(context.TODO(), feeGroup)
		if err != nil {
			return nil, err
		}
		if fees != nil && len(fees) > 0 {
			select {
			case r.feeCacheInsertChan <- feeGroup:
			default:
			}
		}
		return fees, nil
	}

	cacheItem, _ := value.(CacheItem[[]*transaction.Fee])

	if cacheItem.expiredAt.Before(time.Now()) {
		select {
		case r.feeCacheUpdateChan <- feeGroup:
		default:
		}
	}
	fee := cacheItem.item
	return fee, nil
}

func (t *FeeService) applyFee(feeGroup string, amt types.AmountData) ([]*transaction.FeeJoint, error) {
	fees, err := t.getFee(feeGroup)
	if err != nil {
		return nil, err
	}
	feeJoints := make([]*transaction.FeeJoint, 0)
	for _, f := range fees {
		fj, err := f.MaybeApplyFee(amt)
		if err != nil {
			return nil, err
		}
		if fj != nil {
			feeJoints = append(feeJoints, fj)
		}
	}
	return feeJoints, nil
}
