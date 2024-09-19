package foree_service

import (
	"context"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/transaction"
)

const rateCacheExpiry time.Duration = 4 * time.Minute
const rateCacheRefreshInterval time.Duration = 2 * time.Minute

func NewRateService(rateRepo *transaction.RateRepo) *RateService {
	rateService := &RateService{
		rateRepo:               rateRepo,
		rateCacheInsertChan:    make(chan transaction.Rate, 1),
		rateCacheUpdateChan:    make(chan string, 1),
		rateCacheRefreshTicker: time.NewTicker(rateCacheRefreshInterval),
	}
	go rateService.start()
	return rateService
}

type RateService struct {
	rateRepo               *transaction.RateRepo
	cache                  sync.Map
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
				foree_logger.Logger.Error("Rate_Cache_Update_Fail", "rateId", rateId, "cause", err.Error())
			} else {
				r.cache.Swap(rate.GetId(), CacheItem[transaction.Rate]{
					item:      *rate,
					expiredAt: time.Now().Add(rateCacheExpiry),
				})
			}

		case <-r.rateCacheRefreshTicker.C:
			length := 0
			r.cache.Range(func(k, _ interface{}) bool {
				length += 1
				rateId, _ := k.(string)
				rate, err := r.rateRepo.GetUniqueRateById(context.TODO(), rateId)
				if err != nil {
					foree_logger.Logger.Error("Rate_Cache_Refresh_Fail", "rateId", rateId, "cause", err.Error())
				} else {
					r.cache.Swap(rate.GetId(), CacheItem[transaction.Rate]{
						item:      *rate,
						expiredAt: time.Now().Add(rateCacheExpiry),
					})
				}
				return true
			})
			if length > 64 {
				foree_logger.Logger.Warn("Rate_Cache_Refresh", "message", "Size of rate cache is greater than 64, please check if rateCacheExpiry and rateCacheRefreshInterval are still suitable with this cache size.")
			}
		}
	}
}

func (r *RateService) GetRate(src, dest string) (*transaction.Rate, error) {
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
		return rate, nil
	}

	cacheItem, _ := value.(CacheItem[transaction.Rate])

	if cacheItem.expiredAt.Before(time.Now()) {
		select {
		case r.rateCacheUpdateChan <- rateId:
		default:
		}
	}
	//Return old data.
	rate := cacheItem.item
	return &rate, nil
}
