package foree_service

import (
	"context"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
)

const PromotionOnboard = "ONBOARD_PROMOTION"
const PromotionReferral = "REFERRAL_PROMOTION"

const promotionCacheExpiry time.Duration = 4 * time.Minute
const promotionCacheRefreshInterval time.Duration = 2 * time.Minute

func NewPromotionService(promotionRepo *promotion.PromotionRepo) *PromotionService {
	promotionService := &PromotionService{
		promotionRepo:               promotionRepo,
		promotionCacheInsertChan:    make(chan string, 1),
		promotionCacheUpdateChan:    make(chan string, 1),
		promotionCacheRefreshTicker: time.NewTicker(txLimitCacheRefreshInterval),
	}

	return promotionService
}

type PromotionService struct {
	promotionRepo               *promotion.PromotionRepo
	cache                       sync.Map
	promotionCacheInsertChan    chan string
	promotionCacheUpdateChan    chan string
	promotionCacheRefreshTicker *time.Ticker
	pp                          promotion.Promotion
}

func (p *PromotionService) start() {
	for {
		select {
		case promotionName := <-p.promotionCacheInsertChan:
			promo, err := p.promotionRepo.GetUniquePromotionByName(context.TODO(), promotionName)
			if err != nil {
				foree_logger.Logger.Error("Promotion_Cache_Insert_Fail", "promotionName", promotionName, "cause", err.Error())
			} else {
				p.cache.Store(promotionName, CacheItem[promotion.Promotion]{
					item:      *promo,
					expiredAt: time.Now().Add(promotionCacheExpiry),
				})
			}
		case promotionName := <-p.promotionCacheUpdateChan:
			promo, err := p.promotionRepo.GetUniquePromotionByName(context.TODO(), promotionName)
			if err != nil {
				foree_logger.Logger.Error("Promotion_Cache_Update_Fail", "promotionName", promotionName, "cause", err.Error())
			} else {
				p.cache.Swap(promotionName, CacheItem[promotion.Promotion]{
					item:      *promo,
					expiredAt: time.Now().Add(promotionCacheExpiry),
				})
			}

		case <-p.promotionCacheRefreshTicker.C:
			length := 0
			p.cache.Range(func(k, _ interface{}) bool {
				length += 1
				promotionName, _ := k.(string)
				promo, err := p.promotionRepo.GetUniquePromotionByName(context.TODO(), promotionName)
				if err != nil {
					foree_logger.Logger.Error("Promotion_Cache_Refresh_Fail", "limitGroup", promotionName, "cause", err.Error())
				} else {
					p.cache.Swap(promotionName, CacheItem[promotion.Promotion]{
						item:      *promo,
						expiredAt: time.Now().Add(promotionCacheExpiry),
					})
				}
				return true
			})
			if length > 32 {
				foree_logger.Logger.Warn("Promotion_Cache_Refresh", "message", "Size of rate cache is greater than 32, please check if promotionCacheExpiry and promotionCacheRefreshInterval are still suitable with this cache size.")
			}
		}
	}
}
