package service

import (
	"context"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
)

type SubscriptionService struct {
	subRepo  port.SubscriptionRepository
	storeRepo port.StoreRepository
}

func NewSubscriptionService(subRepo port.SubscriptionRepository, storeRepo port.StoreRepository) *SubscriptionService {
	return &SubscriptionService{subRepo, storeRepo}
}

func (ss *SubscriptionService) CreateSubscription(ctx context.Context, storeID uint64, plan domain.SubscriptionPlan) (*domain.Subscription, error) {
	startDate := time.Now()
	var endDate *time.Time

	if plan == domain.PlanTrial {
		trialEnd := startDate.Add(14 * 24 * time.Hour)
		endDate = &trialEnd
	}

	sub := &domain.Subscription{
		StoreID:   storeID,
		Plan:      plan,
		Status:    domain.SubTrial,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if plan != domain.PlanTrial {
		sub.Status = domain.SubActive
	}

	return ss.subRepo.CreateSubscription(ctx, sub)
}

func (ss *SubscriptionService) GetSubscription(ctx context.Context, storeID uint64) (*domain.Subscription, error) {
	return ss.subRepo.GetSubscriptionByStoreID(ctx, storeID)
}

func (ss *SubscriptionService) RenewSubscription(ctx context.Context, storeID uint64, plan domain.SubscriptionPlan, months int) (*domain.Subscription, error) {
	sub, err := ss.subRepo.GetSubscriptionByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var endDate time.Time

	if sub.EndDate != nil && sub.EndDate.After(now) {
		endDate = sub.EndDate.Add(time.Duration(months) * 30 * 24 * time.Hour)
	} else {
		endDate = now.Add(time.Duration(months) * 30 * 24 * time.Hour)
	}

	sub.Plan = plan
	sub.Status = domain.SubActive
	sub.EndDate = &endDate

	return ss.subRepo.UpdateSubscription(ctx, sub)
}

func (ss *SubscriptionService) ActivateSubscription(ctx context.Context, storeID uint64) (*domain.Subscription, error) {
	sub, err := ss.subRepo.GetSubscriptionByStoreID(ctx, storeID)
	if err != nil {
		return nil, err
	}

	sub.Status = domain.SubActive

	return ss.subRepo.UpdateSubscription(ctx, sub)
}

func (ss *SubscriptionService) CheckSubscriptionStatus(ctx context.Context, storeID uint64) (bool, string, error) {
	sub, err := ss.subRepo.GetSubscriptionByStoreID(ctx, storeID)
	if err != nil {
		return false, "", err
	}

	now := time.Now()

	if sub.Status == domain.SubFrozen {
		return false, "租户已冻结", nil
	}

	if sub.EndDate != nil && sub.EndDate.Before(now) {
		return false, "订阅已过期", nil
	}

	return true, "", nil
}

func (ss *SubscriptionService) ProcessExpiredSubscriptions(ctx context.Context) error {
	expiredSubs, err := ss.subRepo.GetExpiredSubscriptions(ctx)
	if err != nil {
		return err
	}

	for _, sub := range expiredSubs {
		sub.Status = domain.SubFrozen
		_, err := ss.subRepo.UpdateSubscription(ctx, &sub)
		if err != nil {
			return err
		}

		store, err := ss.storeRepo.GetStoreByID(ctx, sub.StoreID)
		if err != nil {
			return err
		}

		store.Status = domain.StoreFrozen
		_, err = ss.storeRepo.UpdateStore(ctx, store)
		if err != nil {
			return err
		}
	}

	return nil
}