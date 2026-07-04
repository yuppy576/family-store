package service

import (
	"context"
	"time"
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
)

type SupplierService struct {
	repo port.SupplierRepository
}

func NewSupplierService(repo port.SupplierRepository) *SupplierService {
	return &SupplierService{repo}
}

func (s *SupplierService) CreateSupplier(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error) {
	supplier.CreatedAt = time.Now()
	supplier.UpdatedAt = time.Now()
	return s.repo.CreateSupplier(ctx, supplier)
}

func (s *SupplierService) GetSupplier(ctx context.Context, id uint64) (*domain.Supplier, error) {
	return s.repo.GetSupplierByID(ctx, id)
}

func (s *SupplierService) ListSuppliers(ctx context.Context, search string, skip, limit uint64) ([]domain.Supplier, error) {
	return s.repo.ListSuppliers(ctx, search, skip, limit)
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) (*domain.Supplier, error) {
	supplier.UpdatedAt = time.Now()
	return s.repo.UpdateSupplier(ctx, supplier)
}

func (s *SupplierService) DeleteSupplier(ctx context.Context, id uint64) error {
	return s.repo.DeleteSupplier(ctx, id)
}

func (s *SupplierService) CreatePurchase(ctx context.Context, purchase *domain.Purchase, items []domain.PurchaseItem) (*domain.Purchase, error) {
	purchase.CreatedAt = time.Now()
	purchase.UpdatedAt = time.Now()
	if purchase.Status == "" {
		purchase.Status = domain.PurchaseCompleted
	}
	p, err := s.repo.CreatePurchase(ctx, purchase)
	if err != nil {
		return nil, err
	}
	// Create purchase items and update product stock
	for _, item := range items {
		item.PurchaseID = p.ID
		item.CreatedAt = time.Now()
		item.TotalPrice = float64(item.Quantity) * item.UnitPrice
		if _, err := s.repo.CreatePurchaseItem(ctx, &item); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (s *SupplierService) GetPurchase(ctx context.Context, id uint64) (*domain.Purchase, error) {
	return s.repo.GetPurchaseByID(ctx, id)
}

func (s *SupplierService) ListPurchases(ctx context.Context, skip, limit uint64) ([]domain.Purchase, error) {
	return s.repo.ListPurchases(ctx, skip, limit)
}
