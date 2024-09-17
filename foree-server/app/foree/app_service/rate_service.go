package foree_service

import (
	"context"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
)

const rateCacheExpiry time.Duration = 5 * time.Minute
const rateCacheRefreshInterval time.Duration = 3 * time.Minute

func NewRateService() *RateService {
	rateService := &RateService{
		rateCache:              make(map[string]CacheItem[transaction.Rate], 8),
		rateCacheInsertChan:    make(chan transaction.Rate, 1),
		rateCacheUpdateChan:    make(chan string, 1),
		rateCacheRefreshTicker: time.NewTicker(rateCacheRefreshInterval),
	}
	rateService.start()
	return rateService
}

type RateService struct {
	rateRepo               *transaction.RateRepo
	cache                  sync.Map
	rateCache              map[string]CacheItem[transaction.Rate]
	rateCacheInsertChan    chan transaction.Rate
	rateCacheUpdateChan    chan string
	rateCacheRefreshTicker *time.Ticker
}

func (r *RateService) start() {
	for {
		select {
		case rate := <-r.rateCacheInsertChan:
			r.cache.Store(rate.GetId(), CacheItem[transaction.Rate]{
				item:      rate,
				expiredAt: time.Now().Add(rateCacheExpiry),
			})
		case rateId := <-r.rateCacheUpdateChan:
			rate, err := r.rateRepo.GetUniqueRateById(context.TODO(), rateId)
			if err != nil {
				foree_logger.Logger.Error("Rate_Cache_Update_Fail", "userId", "rateId", rate, "cause", err.Error())
			} else {
				r.cache.Swap(rate.GetId(), CacheItem[transaction.Rate]{
					item:      *rate,
					expiredAt: time.Now().Add(rateCacheExpiry),
				})
			}

		case _ = <-r.rateCacheRefreshTicker.C:
			length := 0
			r.cache.Range(func(k, _ interface{}) bool {
				rateId, _ := k.(string)
				rate, err := r.rateRepo.GetUniqueRateById(context.TODO(), rateId)
				if err != nil {
					foree_logger.Logger.Error("Rate_Cache_Refresh_Fail", "userId", "rateId", rate, "cause", err.Error())
				} else {
					r.cache.Swap(rate.GetId(), CacheItem[transaction.Rate]{
						item:      *rate,
						expiredAt: time.Now().Add(rateCacheExpiry),
					})
				}
				return true
			})
			if length > 64 {
				foree_logger.Logger.Error("Rate_Cache_Refresh", "message", "Size of rate cache is greater than 64, please check if rateCacheExpiry and rateCacheRefreshInterval are still suitable with this cache size.")
			}
		}
	}
}

func (r *RateService) GetRate(src, dest string, validIn time.Duration) (*transaction.Rate, error) {
	rateId := transaction.GenerateRateId(src, dest)

	value, ok := r.cache.Load(rateId)

	if !ok {
		rate, err := r.rateRepo.GetUniqueRateById(context.TODO(), rateId)
		if err != nil {
			return nil, err
		}
		select {
		case r.rateCacheInsertChan <- *rate:
		default:
		}
	}

	cacheItem, _ := value.(CacheItem[transaction.Rate])

	if cacheItem.expiredAt.Before(time.Now()) {
		r.rateCacheUpdateChan <- rateId
	}
	rate := cacheItem.item
	return &rate, nil
}
