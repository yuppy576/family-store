package service

import (
	"context"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/bagashiz/go-pos/internal/core/util"
)

type StoreService struct {
	storeRepo      port.StoreRepository
	userRepo       port.UserRepository
	tokenSvc       port.TokenService
	cache          port.CacheRepository
	subscriptionSvc port.SubscriptionService
}

func NewStoreService(storeRepo port.StoreRepository, userRepo port.UserRepository, tokenSvc port.TokenService, cache port.CacheRepository, subscriptionSvc port.SubscriptionService) *StoreService {
	return &StoreService{
		storeRepo:      storeRepo,
		userRepo:       userRepo,
		tokenSvc:       tokenSvc,
		cache:          cache,
		subscriptionSvc: subscriptionSvc,
	}
}

func (ss *StoreService) Register(ctx context.Context, name, email, password string) (*domain.Store, *domain.User, string, error) {
	domainName := generateDomain(name)

	trialEnd := time.Now().Add(14 * 24 * time.Hour)

	store := &domain.Store{
		Name:     name,
		Domain:   domainName,
		Status:   domain.StoreTrial,
		TrialEnd: &trialEnd,
	}

	store, err := ss.storeRepo.CreateStore(ctx, store)
	if err != nil {
		return nil, nil, "", err
	}

	ctxWithStoreID := context.WithValue(ctx, "store_id", store.ID)

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, nil, "", domain.ErrInternal
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     domain.Admin,
		StoreID:  store.ID,
	}

	user, err = ss.userRepo.CreateUser(ctxWithStoreID, user)
	if err != nil {
		return nil, nil, "", err
	}

	_, err = ss.subscriptionSvc.CreateSubscription(ctx, store.ID, domain.PlanTrial)
	if err != nil {
		return nil, nil, "", err
	}

	accessToken, err := ss.tokenSvc.CreateToken(user)
	if err != nil {
		return nil, nil, "", domain.ErrTokenCreation
	}

	return store, user, accessToken, nil
}

func (ss *StoreService) GetStoreByDomain(ctx context.Context, domain string) (*domain.Store, error) {
	return ss.storeRepo.GetStoreByDomain(ctx, domain)
}

func (ss *StoreService) GetStoreByID(ctx context.Context, id uint64) (*domain.Store, error) {
	return ss.storeRepo.GetStoreByID(ctx, id)
}

func generateDomain(name string) string {
	cleanName := util.Slugify(name)
	return cleanName + ".store.yuppy576.top"
}