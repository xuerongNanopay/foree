package foree_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
	"xue.io/go-pay/app/foree/referral"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
)

const PromotionOnboard = "ONBOARD_PROMOTION"
const PromotionReferral = "REFERRAL_PROMOTION"

const promotionCacheExpiry time.Duration = 4 * time.Minute
const promotionCacheRefreshInterval time.Duration = 2 * time.Minute

func NewPromotionService(
	promotionRepo *promotion.PromotionRepo,
	rewardRepo *transaction.RewardRepo,
	referralRepo *referral.ReferralRepo,
) *PromotionService {
	promotionService := &PromotionService{
		rewardRepo:                  rewardRepo,
		promotionRepo:               promotionRepo,
		referralRepo:                referralRepo,
		promotionCacheInsertChan:    make(chan string, 1),
		promotionCacheUpdateChan:    make(chan string, 1),
		promotionCacheRefreshTicker: time.NewTicker(promotionCacheRefreshInterval),
	}
	go promotionService.start()
	return promotionService
}

type PromotionService struct {
	*referral.ReferralRepo
	promotionRepo               *promotion.PromotionRepo
	rewardRepo                  *transaction.RewardRepo
	referralRepo                *referral.ReferralRepo
	cache                       sync.Map
	promotionCacheInsertChan    chan string
	promotionCacheUpdateChan    chan string
	promotionCacheRefreshTicker *time.Ticker
}

func (p *PromotionService) start() {
	for {
		select {
		case promotionName := <-p.promotionCacheInsertChan:
			promo, err := p.promotionRepo.GetUniquePromotionByName(context.TODO(), promotionName)
			if err != nil {
				foree_logger.Logger.Error("Promotion_Cache_Insert_FAIL", "promotionName", promotionName, "cause", err.Error())
			} else if promo == nil {
				foree_logger.Logger.Error("Promotion_Cache_Insert_FAIL", "promotionName", promotionName, "cause", "promotion no found")
			} else {
				p.cache.Store(promotionName, CacheItem[promotion.Promotion]{
					item:      *promo,
					expiredAt: time.Now().Add(promotionCacheExpiry),
				})
			}
		case promotionName := <-p.promotionCacheUpdateChan:
			promo, err := p.promotionRepo.GetUniquePromotionByName(context.TODO(), promotionName)
			if err != nil {
				foree_logger.Logger.Error("Promotion_Cache_Update_FAIL", "promotionName", promotionName, "cause", err.Error())
			} else if promo == nil {
				foree_logger.Logger.Error("Promotion_Cache_Update_FAIL", "promotionName", promotionName, "cause", "promotion no found")
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
					foree_logger.Logger.Error("Promotion_Cache_Refresh_FAIL", "promotionName", promotionName, "cause", err.Error())
				} else if promo == nil {
					foree_logger.Logger.Error("Promotion_Cache_Refresh_FAIL", "promotionName", promotionName, "cause", "promotion no found")
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

func (p *PromotionService) getPromotion(promotionName string) (*promotion.Promotion, error) {

	value, ok := p.cache.Load(promotionName)

	if !ok {
		promo, err := p.promotionRepo.GetUniquePromotionByName(context.TODO(), promotionName)
		if err != nil {
			return nil, err
		}
		if promo != nil {
			select {
			case p.promotionCacheInsertChan <- promotionName:
			default:
			}
		}
		return promo, nil
	}

	cacheItem, _ := value.(CacheItem[promotion.Promotion])

	if cacheItem.expiredAt.Before(time.Now()) {
		select {
		case p.promotionCacheUpdateChan <- promotionName:
		default:
		}
	}
	//Return old data.
	promo := cacheItem.item
	return &promo, nil
}

func (p *PromotionService) rewardOnboard(registerUser auth.User) {

	promotion, err := p.getPromotion(PromotionOnboard)

	if err != nil {
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	if promotion == nil {
		foree_logger.Logger.Debug("Reward_Onboard_FAIL", "userId", registerUser.ID, "promotionName", PromotionOnboard, "cause", "promotion no found")
		return
	}

	if !promotion.IsValid() {
		foree_logger.Logger.Debug("Reward_Onboard_FAIL", "userId", registerUser.ID, "promotionName", PromotionOnboard, "cause", "promotion is invalid")
		return
	}
	expiry := time.Now().Add(time.Hour * 24 * 180)
	reward := transaction.Reward{
		Type:        transaction.RewardTypeOnboard,
		Status:      transaction.RewardStatusActive,
		Description: "Onboard Reward",
		Amt:         promotion.Amt,
		OwnerId:     registerUser.ID,
		ExpireAt:    &expiry,
	}

	rewardId, err := p.rewardRepo.InsertReward(context.TODO(), reward)
	if err != nil {
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
	}
	foree_logger.Logger.Info("Reward_Onboard_SUCCESS", "userId", registerUser.ID, "rewardId", rewardId)

}

func (p *PromotionService) initialReferralReward(registerUser auth.User) {
	referral, err := p.referralRepo.GetUniqueReferralByRefereeId(registerUser.ID)
	if err != nil {
		foree_logger.Logger.Error("Initial_Referral_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}
	if referral == nil {
		foree_logger.Logger.Debug("Initial_Referral_Reward_FAIL", "userId", registerUser.ID, "cause", "do not have referrer")
		return
	}

	promotion, err := p.getPromotion(PromotionReferral)

	if err != nil {
		foree_logger.Logger.Error("Initial_Referral_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	if promotion == nil {
		foree_logger.Logger.Warn("Initial_Referral_Reward_FAIL", "userId", registerUser.ID, "promotionName", PromotionReferral, "cause", "promotion no found")
		return
	}

	if !promotion.IsValid() {
		foree_logger.Logger.Debug("Initial_Referral_Reward_FAIL", "userId", registerUser.ID, "promotionName", PromotionReferral, "cause", "promotion is invalid")
		return
	}

	generateReferralRewardDescription := func(refereeId int64) string {
		return fmt.Sprintf("Referral accepted by %v", refereeId)
	}

	expiry := time.Now().Add(time.Hour * 24 * 180)
	reward := transaction.Reward{
		Type:        transaction.RewardTypeReferal,
		Status:      transaction.RewardStatusInitial,
		Description: generateReferralRewardDescription(registerUser.ID),
		Amt:         promotion.Amt,
		OwnerId:     referral.ReferrerId,
		ExpireAt:    &expiry,
	}

	rewardId, err := p.rewardRepo.InsertReward(context.TODO(), reward)
	if err != nil {
		foree_logger.Logger.Error("Initial_Referral_Reward_FAIL", "userId", registerUser.ID, "cause", "referrerId", referral.ID, err.Error())
	}
	foree_logger.Logger.Info("Initial_Referral_Reward_SUCCESS", "userId", registerUser.ID, "referrerId", referral.ID, "rewardId", rewardId)
}
