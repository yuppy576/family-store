package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

/**
 * ConsignmentRepository implements port.ConsignmentRepository interface
 */
type ConsignmentRepository struct {
	db *postgres.DB
}

func NewConsignmentRepository(db *postgres.DB) *ConsignmentRepository {
	return &ConsignmentRepository{db}
}

// ── Consignor ────────────────────────────────────────────────

func (cr *ConsignmentRepository) CreateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error) {
	query := cr.db.QueryBuilder.Insert("consignors").
		Columns("name", "phone", "id_card", "address", "memo", "created_at", "updated_at").
		Values(consignor.Name, consignor.Phone, consignor.IDCard, consignor.Address, consignor.Memo, consignor.CreatedAt, consignor.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&consignor.ID, &consignor.Name, &consignor.Phone, &consignor.IDCard,
		&consignor.Address, &consignor.Memo, &consignor.CreatedAt, &consignor.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return consignor, nil
}

func (cr *ConsignmentRepository) GetConsignorByID(ctx context.Context, id uint64) (*domain.Consignor, error) {
	var c domain.Consignor
	query := cr.db.QueryBuilder.Select("*").From("consignors").Where(sq.Eq{"id": id}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&c.ID, &c.Name, &c.Phone, &c.IDCard, &c.Address, &c.Memo, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, handleDBError(err)
	}
	return &c, nil
}

func (cr *ConsignmentRepository) ListConsignors(ctx context.Context, search string, skip, limit uint64) ([]domain.Consignor, error) {
	q := cr.db.QueryBuilder.Select("*").From("consignors")
	if search != "" {
		q = q.Where(sq.Or{
			sq.ILike{"name": "%" + search + "%"},
			sq.ILike{"phone": "%" + search + "%"},
		})
	}
	q = q.OrderBy("id DESC").Offset(skip).Limit(limit)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consignors []domain.Consignor
	for rows.Next() {
		var c domain.Consignor
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.IDCard, &c.Address, &c.Memo, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		consignors = append(consignors, c)
	}
	return consignors, nil
}

func (cr *ConsignmentRepository) UpdateConsignor(ctx context.Context, consignor *domain.Consignor) (*domain.Consignor, error) {
	query := cr.db.QueryBuilder.Update("consignors").
		Set("name", consignor.Name).
		Set("phone", consignor.Phone).
		Set("id_card", consignor.IDCard).
		Set("address", consignor.Address).
		Set("memo", consignor.Memo).
		Set("updated_at", consignor.UpdatedAt).
		Where(sq.Eq{"id": consignor.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&consignor.ID, &consignor.Name, &consignor.Phone, &consignor.IDCard,
		&consignor.Address, &consignor.Memo, &consignor.CreatedAt, &consignor.UpdatedAt,
	)
	if err != nil {
		return nil, handleDBError(err)
	}
	return consignor, nil
}

func (cr *ConsignmentRepository) DeleteConsignor(ctx context.Context, id uint64) error {
	query := cr.db.QueryBuilder.Delete("consignors").Where(sq.Eq{"id": id})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = cr.db.Exec(ctx, sql, args...)
	return err
}

// ── Consignment ──────────────────────────────────────────────

func (cr *ConsignmentRepository) CreateConsignment(ctx context.Context, cons *domain.Consignment) (*domain.Consignment, error) {
	query := cr.db.QueryBuilder.Insert("consignments").
		Columns("consignor_id", "name", "description", "images", "category",
			"expected_price", "recommended_price", "final_price",
			"commission_rate", "commission_amount", "status",
			"contract_end", "is_vehicle", "memo", "created_at", "updated_at").
		Values(cons.ConsignorID, cons.Name, cons.Description, cons.Images, cons.Category,
			cons.ExpectedPrice, cons.RecommendedPrice, cons.FinalPrice,
			cons.CommissionRate, cons.CommissionAmount, cons.Status,
			cons.ContractEnd, cons.IsVehicle, cons.Memo, cons.CreatedAt, cons.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&cons.ID, &cons.ConsignorID, &cons.Name, &cons.Description, &cons.Images,
		&cons.Category, &cons.ExpectedPrice, &cons.RecommendedPrice, &cons.FinalPrice,
		&cons.CommissionRate, &cons.CommissionAmount, &cons.Status,
		&cons.ContractEnd, &cons.IsVehicle, &cons.Memo, &cons.CreatedAt, &cons.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return cons, nil
}

func (cr *ConsignmentRepository) GetConsignmentByID(ctx context.Context, id uint64) (*domain.Consignment, error) {
	var c domain.Consignment
	query := cr.db.QueryBuilder.Select("*").From("consignments").Where(sq.Eq{"id": id}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&c.ID, &c.ConsignorID, &c.Name, &c.Description, &c.Images,
		&c.Category, &c.ExpectedPrice, &c.RecommendedPrice, &c.FinalPrice,
		&c.CommissionRate, &c.CommissionAmount, &c.Status,
		&c.ContractEnd, &c.IsVehicle, &c.Memo, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, handleDBError(err)
	}
	return &c, nil
}

func (cr *ConsignmentRepository) ListConsignments(ctx context.Context, status string, skip, limit uint64) ([]domain.Consignment, error) {
	q := cr.db.QueryBuilder.Select("*").From("consignments")
	if status != "" {
		q = q.Where(sq.Eq{"status": status})
	}
	q = q.OrderBy("id DESC").Offset(skip).Limit(limit)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.Consignment
	for rows.Next() {
		var c domain.Consignment
		if err := rows.Scan(&c.ID, &c.ConsignorID, &c.Name, &c.Description, &c.Images,
			&c.Category, &c.ExpectedPrice, &c.RecommendedPrice, &c.FinalPrice,
			&c.CommissionRate, &c.CommissionAmount, &c.Status,
			&c.ContractEnd, &c.IsVehicle, &c.Memo, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, nil
}
func (cr *ConsignmentRepository) ListExpiringConsignments(ctx context.Context, withinDays int32) ([]domain.Consignment, error) {
	rows, err := cr.db.Query(ctx,
		"SELECT * FROM consignments WHERE status = $1 AND contract_end <= CURRENT_DATE + CAST($2 AS INTEGER) ORDER BY contract_end ASC",
		"ON_SALE", withinDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Consignment
	for rows.Next() {
		var c domain.Consignment
		if err := rows.Scan(&c.ID, &c.ConsignorID, &c.Name, &c.Description, &c.Images,
			&c.Category, &c.ExpectedPrice, &c.RecommendedPrice, &c.FinalPrice,
			&c.CommissionRate, &c.CommissionAmount, &c.Status,
			&c.ContractEnd, &c.IsVehicle, &c.Memo, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, nil
}

func (cr *ConsignmentRepository) UpdateConsignment(ctx context.Context, cons *domain.Consignment) (*domain.Consignment, error) {
	query := cr.db.QueryBuilder.Update("consignments").
		Set("name", cons.Name).Set("description", cons.Description).
		Set("category", cons.Category).Set("expected_price", cons.ExpectedPrice).
		Set("recommended_price", cons.RecommendedPrice).Set("final_price", cons.FinalPrice).
		Set("commission_rate", cons.CommissionRate).Set("commission_amount", cons.CommissionAmount).
		Set("status", cons.Status).Set("contract_end", cons.ContractEnd).
		Set("memo", cons.Memo).Set("updated_at", cons.UpdatedAt).
		Where(sq.Eq{"id": cons.ID}).Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&cons.ID, &cons.ConsignorID, &cons.Name, &cons.Description, &cons.Images,
		&cons.Category, &cons.ExpectedPrice, &cons.RecommendedPrice, &cons.FinalPrice,
		&cons.CommissionRate, &cons.CommissionAmount, &cons.Status,
		&cons.ContractEnd, &cons.IsVehicle, &cons.Memo, &cons.CreatedAt, &cons.UpdatedAt,
	)
	if err != nil {
		return nil, handleDBError(err)
	}
	return cons, nil
}

func (cr *ConsignmentRepository) DeleteConsignment(ctx context.Context, id uint64) error {
	query := cr.db.QueryBuilder.Delete("consignments").Where(sq.Eq{"id": id})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = cr.db.Exec(ctx, sql, args...)
	return err
}

// ── Vehicle ──────────────────────────────────────────────────

func (cr *ConsignmentRepository) CreateVehicle(ctx context.Context, v *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error) {
	q := cr.db.QueryBuilder.Insert("consignment_vehicles").
		Columns("consignment_id", "vin", "plate_number", "brand", "model", "year",
			"mileage", "displacement", "color", "inspection_expire", "insurance_expire",
			"created_at", "updated_at").
		Values(v.ConsignmentID, v.VIN, v.PlateNumber, v.Brand, v.Model, v.Year,
			v.Mileage, v.Displacement, v.Color, v.InspectionExpire, v.InsuranceExpire,
			v.CreatedAt, v.UpdatedAt).
		Suffix("RETURNING *")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&v.ID, &v.ConsignmentID, &v.VIN, &v.PlateNumber, &v.Brand, &v.Model,
		&v.Year, &v.Mileage, &v.Displacement, &v.Color,
		&v.InspectionExpire, &v.InsuranceExpire, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (cr *ConsignmentRepository) GetVehicleByConsignmentID(ctx context.Context, consignmentID uint64) (*domain.ConsignmentVehicle, error) {
	var v domain.ConsignmentVehicle
	q := cr.db.QueryBuilder.Select("*").From("consignment_vehicles").Where(sq.Eq{"consignment_id": consignmentID}).Limit(1)
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&v.ID, &v.ConsignmentID, &v.VIN, &v.PlateNumber, &v.Brand, &v.Model,
		&v.Year, &v.Mileage, &v.Displacement, &v.Color,
		&v.InspectionExpire, &v.InsuranceExpire, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, handleDBError(err)
	}
	return &v, nil
}

func (cr *ConsignmentRepository) UpdateVehicle(ctx context.Context, v *domain.ConsignmentVehicle) (*domain.ConsignmentVehicle, error) {
	q := cr.db.QueryBuilder.Update("consignment_vehicles").
		Set("vin", v.VIN).Set("plate_number", v.PlateNumber).Set("brand", v.Brand).
		Set("model", v.Model).Set("year", v.Year).Set("mileage", v.Mileage).
		Set("displacement", v.Displacement).Set("color", v.Color).
		Set("inspection_expire", v.InspectionExpire).Set("insurance_expire", v.InsuranceExpire).
		Set("updated_at", v.UpdatedAt).
		Where(sq.Eq{"id": v.ID}).Suffix("RETURNING *")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&v.ID, &v.ConsignmentID, &v.VIN, &v.PlateNumber, &v.Brand, &v.Model,
		&v.Year, &v.Mileage, &v.Displacement, &v.Color,
		&v.InspectionExpire, &v.InsuranceExpire, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, handleDBError(err)
	}
	return v, nil
}

// ── Transfer Progress ────────────────────────────────────────

func (cr *ConsignmentRepository) CreateTransferProgress(ctx context.Context, p *domain.TransferProgress) (*domain.TransferProgress, error) {
	q := cr.db.QueryBuilder.Insert("consignment_transfer_progress").
		Columns("vehicle_id", "status", "remark", "attachment", "operator", "created_at").
		Values(p.VehicleID, p.Status, p.Remark, p.Attachment, p.Operator, p.CreatedAt).
		Suffix("RETURNING *")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&p.ID, &p.VehicleID, &p.Status, &p.Remark, &p.Attachment, &p.Operator, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (cr *ConsignmentRepository) ListTransferProgress(ctx context.Context, vehicleID uint64) ([]domain.TransferProgress, error) {
	q := cr.db.QueryBuilder.Select("*").From("consignment_transfer_progress").
		Where(sq.Eq{"vehicle_id": vehicleID}).OrderBy("id DESC")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.TransferProgress
	for rows.Next() {
		var p domain.TransferProgress
		if err := rows.Scan(&p.ID, &p.VehicleID, &p.Status, &p.Remark, &p.Attachment, &p.Operator, &p.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, nil
}

// ── Settlement ───────────────────────────────────────────────

func (cr *ConsignmentRepository) CreateSettlement(ctx context.Context, s *domain.ConsignmentSettlement) (*domain.ConsignmentSettlement, error) {
	q := cr.db.QueryBuilder.Insert("consignment_settlements").
		Columns("consignment_id", "type", "sale_price", "commission_amount",
			"settlement_amount", "renewal_fee", "renewal_months", "new_end_date", "remark", "created_at").
		Values(s.ConsignmentID, s.Type, s.SalePrice, s.CommissionAmount,
			s.SettlementAmount, s.RenewalFee, s.RenewalMonths, s.NewEndDate, s.Remark, s.CreatedAt).
		Suffix("RETURNING *")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&s.ID, &s.ConsignmentID, &s.Type, &s.SalePrice, &s.CommissionAmount,
		&s.SettlementAmount, &s.RenewalFee, &s.RenewalMonths, &s.NewEndDate,
		&s.Remark, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (cr *ConsignmentRepository) ListSettlements(ctx context.Context, consignmentID uint64) ([]domain.ConsignmentSettlement, error) {
	q := cr.db.QueryBuilder.Select("*").From("consignment_settlements").
		Where(sq.Eq{"consignment_id": consignmentID}).OrderBy("id DESC")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.ConsignmentSettlement
	for rows.Next() {
		var s domain.ConsignmentSettlement
		if err := rows.Scan(&s.ID, &s.ConsignmentID, &s.Type, &s.SalePrice,
			&s.CommissionAmount, &s.SettlementAmount, &s.RenewalFee, &s.RenewalMonths,
			&s.NewEndDate, &s.Remark, &s.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

// handleDBError translates common DB errors to domain errors
func handleDBError(err error) error {
	if err.Error() == "no rows in result set" {
		return domain.ErrDataNotFound
	}
	return err
}
