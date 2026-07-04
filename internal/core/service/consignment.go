package service

import (
	"context"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
)

/**
 * ConsignmentService implements port.ConsignmentService interface
 */
type ConsignmentService struct {
	repo port.ConsignmentRepository
}

// NewConsignmentService creates a new consignment service instance
func NewConsignmentService(repo port.ConsignmentRepository) *ConsignmentService {
	return &ConsignmentService{repo}
}

// ── Consignor ────────────────────────────────────────────────

func (cs *ConsignmentService) CreateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error) {
	consignor.CreatedAt = time.Now()
	consignor.UpdatedAt = time.Now()
	return cs.repo.CreateConsignor(ctx, consignor)
}

func (cs *ConsignmentService) GetConsignor(ctx context.Context, id uint64) (*domain.Consignor, error) {
	return cs.repo.GetConsignorByID(ctx, id)
}

func (cs *ConsignmentService) ListConsignors(ctx context.Context, search string, skip, limit uint64) ([]domain.Consignor, error) {
	return cs.repo.ListConsignors(ctx, search, skip, limit)
}

func (cs *ConsignmentService) UpdateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error) {
	consignor.UpdatedAt = time.Now()
	return cs.repo.UpdateConsignor(ctx, consignor)
}

func (cs *ConsignmentService) DeleteConsignor(ctx context.Context, id uint64) error {
	return cs.repo.DeleteConsignor(ctx, id)
}

// ── Consignment ──────────────────────────────────────────────

func (cs *ConsignmentService) CreateConsignment(ctx context.Context, consignment *domain.Consignment) (*domain.Consignment, error) {
	consignment.CreatedAt = time.Now()
	consignment.UpdatedAt = time.Now()
	if consignment.Status == "" {
		consignment.Status = domain.ConsignmentStatusOnSale
	}
	return cs.repo.CreateConsignment(ctx, consignment)
}

func (cs *ConsignmentService) GetConsignment(ctx context.Context, id uint64) (*domain.Consignment, error) {
	return cs.repo.GetConsignmentByID(ctx, id)
}

func (cs *ConsignmentService) ListConsignments(ctx context.Context, status string, skip, limit uint64) ([]domain.Consignment, error) {
	return cs.repo.ListConsignments(ctx, status, skip, limit)
}

func (cs *ConsignmentService) ListExpiringConsignments(ctx context.Context, withinDays int32) ([]domain.Consignment, error) {
	return cs.repo.ListExpiringConsignments(ctx, withinDays)
}
func (cs *ConsignmentService) UpdateConsignment(ctx context.Context, consignment *domain.Consignment) (*domain.Consignment, error) {
	consignment.UpdatedAt = time.Now()
	return cs.repo.UpdateConsignment(ctx, consignment)
}

func (cs *ConsignmentService) DeleteConsignment(ctx context.Context, id uint64) error {
	return cs.repo.DeleteConsignment(ctx, id)
}

// ── Vehicle ──────────────────────────────────────────────────

func (cs *ConsignmentService) CreateVehicle(ctx context.Context, vehicle *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error) {
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()
	return cs.repo.CreateVehicle(ctx, vehicle)
}

func (cs *ConsignmentService) GetVehicle(ctx context.Context, consignmentID uint64) (*domain.ConsignmentVehicle, error) {
	return cs.repo.GetVehicleByConsignmentID(ctx, consignmentID)
}

func (cs *ConsignmentService) UpdateVehicle(ctx context.Context, vehicle *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error) {
	vehicle.UpdatedAt = time.Now()
	return cs.repo.UpdateVehicle(ctx, vehicle)
}

// ── Transfer Progress ────────────────────────────────────────

func (cs *ConsignmentService) CreateTransferProgress(ctx context.Context, progress *domain.TransferProgress) (*domain.TransferProgress, error) {
	progress.CreatedAt = time.Now()
	return cs.repo.CreateTransferProgress(ctx, progress)
}

func (cs *ConsignmentService) ListTransferProgress(ctx context.Context, vehicleID uint64) ([]domain.TransferProgress, error) {
	return cs.repo.ListTransferProgress(ctx, vehicleID)
}

// ── Settlement ───────────────────────────────────────────────

func (cs *ConsignmentService) CreateSettlement(ctx context.Context, settlement *domain.ConsignmentSettlement) (*domain.ConsignmentSettlement, error) {
	settlement.CreatedAt = time.Now()
	return cs.repo.CreateSettlement(ctx, settlement)
}

func (cs *ConsignmentService) ListSettlements(ctx context.Context, consignmentID uint64) ([]domain.ConsignmentSettlement, error) {
	return cs.repo.ListSettlements(ctx, consignmentID)
}
