package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

type StoreRepository interface {
	CreateStore(ctx context.Context, store *domain.Store) (*domain.Store, error)
	GetStoreByDomain(ctx context.Context, domain string) (*domain.Store, error)
	GetStoreByID(ctx context.Context, id uint64) (*domain.Store, error)
	GetAllStores(ctx context.Context) ([]domain.Store, error)
	UpdateStore(ctx context.Context, store *domain.Store) (*domain.Store, error)
}

type StoreService interface {
	Register(ctx context.Context, name, email, password string) (*domain.Store, *domain.User, string, error)
	GetStoreByDomain(ctx context.Context, domain string) (*domain.Store, error)
	GetStoreByID(ctx context.Context, id uint64) (*domain.Store, error)
}