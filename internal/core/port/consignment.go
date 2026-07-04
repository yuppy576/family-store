package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

//go:generate mockgen -source=consignment.go -destination=mock/consignment.go -package=mock

// ConsignmentRepository is an interface for consignment data operations
type ConsignmentRepository interface {
	CreateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error)
	GetConsignorByID(ctx context.Context, id uint64) (*domain.Consignor, error)
	ListConsignors(ctx context.Context, search string, skip, limit uint64) ([]domain.Consignor, error)
	UpdateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error)
	DeleteConsignor(ctx context.Context, id uint64) error

	CreateConsignment(ctx context.Context, consignment *domain.Consignment) (*domain.Consignment, error)
	GetConsignmentByID(ctx context.Context, id uint64) (*domain.Consignment, error)
	ListConsignments(ctx context.Context, status string, skip, limit uint64) ([]domain.Consignment, error)
	ListExpiringConsignments(ctx context.Context, withinDays int32) ([]domain.Consignment, error)
	UpdateConsignment(ctx context.Context, consignment *domain.Consignment) (*domain.Consignment, error)
	DeleteConsignment(ctx context.Context, id uint64) error

	CreateVehicle(ctx context.Context, vehicle *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error)
	GetVehicleByConsignmentID(ctx context.Context, consignmentID uint64) (*domain.ConsignmentVehicle, error)
	UpdateVehicle(ctx context.Context, vehicle *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error)

	CreateTransferProgress(ctx context.Context, progress *domain.TransferProgress) (*domain.TransferProgress, error)
	ListTransferProgress(ctx context.Context, vehicleID uint64) ([]domain.TransferProgress, error)

	CreateSettlement(ctx context.Context, settlement *domain.ConsignmentSettlement) (*domain.ConsignmentSettlement, error)
	ListSettlements(ctx context.Context, consignmentID uint64) ([]domain.ConsignmentSettlement, error)
}

// ConsignmentService is an interface for consignment business logic
type ConsignmentService interface {
	CreateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error)
	GetConsignor(ctx context.Context, id uint64) (*domain.Consignor, error)
	ListConsignors(ctx context.Context, search string, skip, limit uint64) ([]domain.Consignor, error)
	UpdateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error)
	DeleteConsignor(ctx context.Context, id uint64) error

	CreateConsignment(ctx context.Context, consignment *domain.Consignment) (*domain.Consignment, error)
	GetConsignment(ctx context.Context, id uint64) (*domain.Consignment, error)
	ListConsignments(ctx context.Context, status string, skip, limit uint64) ([]domain.Consignment, error)
	ListExpiringConsignments(ctx context.Context, withinDays int32) ([]domain.Consignment, error)
	UpdateConsignment(ctx context.Context, consignment *domain.Consignment) (*domain.Consignment, error)
	DeleteConsignment(ctx context.Context, id uint64) error

	CreateVehicle(ctx context.Context, vehicle *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error)
	GetVehicle(ctx context.Context, consignmentID uint64) (*domain.ConsignmentVehicle, error)
	UpdateVehicle(ctx context.Context, vehicle *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error)

	CreateTransferProgress(ctx context.Context, progress *domain.TransferProgress) (*domain.TransferProgress, error)
	ListTransferProgress(ctx context.Context, vehicleID uint64) ([]domain.TransferProgress, error)

	CreateSettlement(ctx context.Context, settlement *domain.ConsignmentSettlement) (*domain.ConsignmentSettlement, error)
	ListSettlements(ctx context.Context, consignmentID uint64) ([]domain.ConsignmentSettlement, error)
}
