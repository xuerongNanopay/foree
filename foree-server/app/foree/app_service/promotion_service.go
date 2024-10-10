package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
	"xue.io/go-pay/app/foree/referral"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
)

const (
	PromotionOnboard  = "ONBOARD_PROMOTION"
	PromotionReferrer = "REFERRER_PROMOTION"
	PromotionReferee  = "REFEREE_PROMOTION"
)

const promotionCacheExpiry time.Duration = 2 * time.Minute
const promotionCacheRefreshInterval time.Duration = 1 * time.Minute

func NewPromotionService(
	db *sql.DB,
	promotionRepo *promotion.PromotionRepo,
	rewardRepo *promotion.RewardRepo,
	referralRepo *referral.ReferralRepo,
	promotionRewardJointRepo *promotion.PromotionRewardJointRepo,
) *PromotionService {
	promotionService := &PromotionService{
		db:                          db,
		rewardRepo:                  rewardRepo,
		promotionRepo:               promotionRepo,
		referralRepo:                referralRepo,
		promotionRewardJointRepo:    promotionRewardJointRepo,
		promotionCacheInsertChan:    make(chan string, 1),
		promotionCacheUpdateChan:    make(chan string, 1),
		promotionCacheRefreshTicker: time.NewTicker(promotionCacheRefreshInterval),
	}
	go promotionService.startPromotionCacher()
	return promotionService
}

type PromotionService struct {
	db                          *sql.DB
	promotionRepo               *promotion.PromotionRepo
	rewardRepo                  *promotion.RewardRepo
	referralRepo                *referral.ReferralRepo
	promotionRewardJointRepo    *promotion.PromotionRewardJointRepo
	cache                       sync.Map
	promotionCacheInsertChan    chan string
	promotionCacheUpdateChan    chan string
	promotionCacheRefreshTicker *time.Ticker
}

func (p *PromotionService) startPromotionCacher() {
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
	promo, err := p.getPromotion(PromotionOnboard)

	if err != nil {
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	if promo == nil {
		foree_logger.Logger.Debug("Reward_Onboard_FAIL", "userId", registerUser.ID, "promotionName", PromotionOnboard, "cause", "promotion no found")
		return
	}

	if !promo.IsValid() {
		foree_logger.Logger.Debug("Reward_Onboard_FAIL", "userId", registerUser.ID, "promotionName", PromotionOnboard, "cause", "promotion is invalid")
		return
	}

	prj, err := p.promotionRewardJointRepo.GetUniquePromotionRewardJointByPromotionIdAndOwnerId(promo.ID, registerUser.ID)
	if err != nil {
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}
	if prj != nil {
		foree_logger.Logger.Warn("Reward_Onboard_FAIL", "userId", registerUser.ID, "promotionName", PromotionOnboard, "cause", "user already got the promotion")
		return
	}

	expiry := time.Now().Add(time.Hour * 24 * 180)
	reward := promotion.Reward{
		Type:        promotion.RewardTypeOnboard,
		Status:      promotion.RewardStatusActive,
		Description: "Onboard Reward",
		Amt:         promo.Amt,
		OwnerId:     registerUser.ID,
		ExpireAt:    &expiry,
	}

	ctx := context.TODO()
	dTx, err := p.db.Begin()
	if err != nil {
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
		dTx.Rollback()
		return
	}
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	rewardId, err := p.rewardRepo.InsertReward(ctx, reward)
	if err != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	// Update join.
	_, err = p.promotionRewardJointRepo.InsertPromotionRewardJoint(ctx, promotion.PromotionRewardJoint{
		PromotionId:      promo.ID,
		PromotionVersion: promo.Version,
		RewardId:         rewardId,
		OwnerId:          registerUser.ID,
	})
	if err != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	if err = dTx.Commit(); err != nil {
		foree_logger.Logger.Error("Reward_Onboard_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	foree_logger.Logger.Info("Reward_Onboard_SUCCESS", "userId", registerUser.ID, "rewardId", rewardId)
}

func (p *PromotionService) rewardReferrer(registerUser auth.User) {
	referral, err := p.referralRepo.GetUniqueReferralByRefereeId(registerUser.ID)
	if err != nil {
		foree_logger.Logger.Error("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}
	if referral == nil {
		foree_logger.Logger.Debug("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", "do not have referrer")
		return
	}

	promo, err := p.getPromotion(PromotionReferrer)

	if err != nil {
		foree_logger.Logger.Error("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	if promo == nil {
		foree_logger.Logger.Warn("Referrer_Reward_FAIL", "userId", registerUser.ID, "promotionName", PromotionReferrer, "cause", "promotion no found")
		return
	}

	if !promo.IsValid() {
		foree_logger.Logger.Debug("Referrer_Reward_FAIL", "userId", registerUser.ID, "promotionName", PromotionReferrer, "cause", "promotion is invalid")
		return
	}

	prj, err := p.promotionRewardJointRepo.GetUniquePromotionRewardJointByPromotionIdAndReferrerIdAndRefereeId(promo.ID, referral.ID, registerUser.ID)

	if err != nil {
		foree_logger.Logger.Error("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	if prj != nil {
		foree_logger.Logger.Warn("Referrer_Reward_FAIL", "userId", registerUser.ID, "promotionName", PromotionOnboard, "cause", "user already got the promotion")
		return
	}

	expiry := time.Now().Add(time.Hour * 24 * 180)
	reward := promotion.Reward{
		Type:        promotion.RewardTypeReferrer,
		Status:      promotion.RewardStatusActive,
		Description: fmt.Sprintf("Referrer reward from %v %v", registerUser.FirstName, registerUser.LastName),
		Amt:         promo.Amt,
		OwnerId:     referral.ReferrerId,
		ExpireAt:    &expiry,
	}

	ctx := context.TODO()
	dTx, err := p.db.Begin()
	if err != nil {
		foree_logger.Logger.Error("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		dTx.Rollback()
		return
	}
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	rewardId, err := p.rewardRepo.InsertReward(ctx, reward)
	if err != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	// Update join.
	_, err = p.promotionRewardJointRepo.InsertPromotionRewardJoint(ctx, promotion.PromotionRewardJoint{
		PromotionId:      promo.ID,
		PromotionVersion: promo.Version,
		RewardId:         rewardId,
		ReferrerId:       referral.ReferrerId,
		RefereeId:        referral.RefereeId,
		OwnerId:          referral.ReferrerId,
	})
	if err != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	if err = dTx.Commit(); err != nil {
		foree_logger.Logger.Error("Referrer_Reward_FAIL", "userId", registerUser.ID, "cause", err.Error())
		return
	}

	foree_logger.Logger.Info("Referrer_Reward_SUCCESS", "userId", registerUser.ID, "rewardId", rewardId)

}
