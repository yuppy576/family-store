package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	GetSubscriptionByStoreID(ctx context.Context, storeID uint64) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	GetExpiringSubscriptions(ctx context.Context, days int) ([]domain.Subscription, error)
	GetExpiredSubscriptions(ctx context.Context) ([]domain.Subscription, error)
}

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, storeID uint64, plan domain.SubscriptionPlan) (*domain.Subscription, error)
	GetSubscription(ctx context.Context, storeID uint64) (*domain.Subscription, error)
	RenewSubscription(ctx context.Context, storeID uint64, plan domain.SubscriptionPlan, months int) (*domain.Subscription, error)
	ActivateSubscription(ctx context.Context, storeID uint64) (*domain.Subscription, error)
	CheckSubscriptionStatus(ctx context.Context, storeID uint64) (bool, string, error)
	ProcessExpiredSubscriptions(ctx context.Context) error
}