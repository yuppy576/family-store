package port

import (
	"context"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

type SupplierRepository interface {
	CreateSupplier(ctx context.Context, s *domain.Supplier) (*domain.Supplier, error)
	GetSupplierByID(ctx context.Context, id uint64) (*domain.Supplier, error)
	ListSuppliers(ctx context.Context, search string, skip, limit uint64) ([]domain.Supplier, error)
	UpdateSupplier(ctx context.Context, s *domain.Supplier) (*domain.Supplier, error)
	DeleteSupplier(ctx context.Context, id uint64) error

	CreatePurchase(ctx context.Context, p *domain.Purchase) (*domain.Purchase, error)
	GetPurchaseByID(ctx context.Context, id uint64) (*domain.Purchase, error)
	ListPurchases(ctx context.Context, skip, limit uint64) ([]domain.Purchase, error)
	UpdatePurchaseStatus(ctx context.Context, id uint64, status domain.PurchaseStatus) error

	CreatePurchaseItem(ctx context.Context, pi *domain.PurchaseItem) (*domain.PurchaseItem, error)
	ListPurchaseItems(ctx context.Context, purchaseID uint64) ([]domain.PurchaseItem, error)
}

type SupplierService interface {
	CreateSupplier(ctx context.Context, s *domain.Supplier) (*domain.Supplier, error)
	GetSupplier(ctx context.Context, id uint64) (*domain.Supplier, error)
	ListSuppliers(ctx context.Context, search string, skip, limit uint64) ([]domain.Supplier, error)
	UpdateSupplier(ctx context.Context, s *domain.Supplier) (*domain.Supplier, error)
	DeleteSupplier(ctx context.Context, id uint64) error

	CreatePurchase(ctx context.Context, p *domain.Purchase, items []domain.PurchaseItem) (*domain.Purchase, error)
	GetPurchase(ctx context.Context, id uint64) (*domain.Purchase, error)
	ListPurchases(ctx context.Context, skip, limit uint64) ([]domain.Purchase, error)
	ListPurchaseItems(ctx context.Context, purchaseID uint64) ([]domain.PurchaseItem, error)
}
