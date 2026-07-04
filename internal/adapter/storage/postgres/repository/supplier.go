package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

type SupplierRepository struct{ db *postgres.DB }

func NewSupplierRepository(db *postgres.DB) *SupplierRepository { return &SupplierRepository{db} }

func (r *SupplierRepository) CreateSupplier(ctx context.Context, s *domain.Supplier) (*domain.Supplier, error) {
	q := r.db.QueryBuilder.Insert("suppliers").Columns("name", "contact_person", "phone", "address", "memo", "created_at", "updated_at").
		Values(s.Name, s.ContactPerson, s.Phone, s.Address, s.Memo, s.CreatedAt, s.UpdatedAt).Suffix("RETURNING *")
	sql, args, _ := q.ToSql()
	err := r.db.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Address, &s.Memo, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *SupplierRepository) GetSupplierByID(ctx context.Context, id uint64) (*domain.Supplier, error) {
	var s domain.Supplier
	sql, args, _ := r.db.QueryBuilder.Select("*").From("suppliers").Where(sq.Eq{"id": id}).Limit(1).ToSql()
	err := r.db.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Address, &s.Memo, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, handleDBError(err)
	}
	return &s, nil
}

func (r *SupplierRepository) ListSuppliers(ctx context.Context, search string, skip, limit uint64) ([]domain.Supplier, error) {
	q := r.db.QueryBuilder.Select("*").From("suppliers")
	if search != "" {
		q = q.Where(sq.Or{sq.ILike{"name": "%" + search + "%"}, sq.ILike{"phone": "%" + search + "%"}})
	}
	q = q.OrderBy("id DESC").Offset(skip).Limit(limit)
	sql, args, _ := q.ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Supplier
	for rows.Next() {
		var s domain.Supplier
		if err := rows.Scan(&s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Address, &s.Memo, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, nil
}

func (r *SupplierRepository) UpdateSupplier(ctx context.Context, s *domain.Supplier) (*domain.Supplier, error) {
	q := r.db.QueryBuilder.Update("suppliers").Set("name", s.Name).Set("contact_person", s.ContactPerson).
		Set("phone", s.Phone).Set("address", s.Address).Set("memo", s.Memo).Set("updated_at", s.UpdatedAt).
		Where(sq.Eq{"id": s.ID}).Suffix("RETURNING *")
	sql, args, _ := q.ToSql()
	err := r.db.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.Name, &s.ContactPerson, &s.Phone, &s.Address, &s.Memo, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *SupplierRepository) DeleteSupplier(ctx context.Context, id uint64) error {
	sql, args, _ := r.db.QueryBuilder.Delete("suppliers").Where(sq.Eq{"id": id}).ToSql()
	_, err := r.db.Exec(ctx, sql, args...)
	return err
}

func (r *SupplierRepository) CreatePurchase(ctx context.Context, p *domain.Purchase) (*domain.Purchase, error) {
	q := r.db.QueryBuilder.Insert("purchases").Columns("supplier_id", "operator", "total_amount", "status", "remark", "created_at", "updated_at").
		Values(p.SupplierID, p.Operator, p.TotalAmount, p.Status, p.Remark, p.CreatedAt, p.UpdatedAt).Suffix("RETURNING *")
	sql, args, _ := q.ToSql()
	err := r.db.QueryRow(ctx, sql, args...).Scan(&p.ID, &p.SupplierID, &p.Operator, &p.TotalAmount, &p.Status, &p.Remark, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *SupplierRepository) GetPurchaseByID(ctx context.Context, id uint64) (*domain.Purchase, error) {
	var p domain.Purchase
	sql, args, _ := r.db.QueryBuilder.Select("*").From("purchases").Where(sq.Eq{"id": id}).Limit(1).ToSql()
	err := r.db.QueryRow(ctx, sql, args...).Scan(&p.ID, &p.SupplierID, &p.Operator, &p.TotalAmount, &p.Status, &p.Remark, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, handleDBError(err)
	}
	return &p, nil
}

func (r *SupplierRepository) ListPurchases(ctx context.Context, skip, limit uint64) ([]domain.Purchase, error) {
	q := r.db.QueryBuilder.Select("*").From("purchases").OrderBy("id DESC").Offset(skip).Limit(limit)
	sql, args, _ := q.ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.Purchase
	for rows.Next() {
		var p domain.Purchase
		if err := rows.Scan(&p.ID, &p.SupplierID, &p.Operator, &p.TotalAmount, &p.Status, &p.Remark, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, p)
	}
	return items, nil
}

func (r *SupplierRepository) UpdatePurchaseStatus(ctx context.Context, id uint64, status domain.PurchaseStatus) error {
	sql, args, _ := r.db.QueryBuilder.Update("purchases").Set("status", status).Where(sq.Eq{"id": id}).ToSql()
	_, err := r.db.Exec(ctx, sql, args...)
	return err
}

func (r *SupplierRepository) CreatePurchaseItem(ctx context.Context, pi *domain.PurchaseItem) (*domain.PurchaseItem, error) {
	q := r.db.QueryBuilder.Insert("purchase_items").Columns("purchase_id", "product_id", "quantity", "unit_price", "total_price", "created_at").
		Values(pi.PurchaseID, pi.ProductID, pi.Quantity, pi.UnitPrice, pi.TotalPrice, pi.CreatedAt).Suffix("RETURNING *")
	sql, args, _ := q.ToSql()
	err := r.db.QueryRow(ctx, sql, args...).Scan(&pi.ID, &pi.PurchaseID, &pi.ProductID, &pi.Quantity, &pi.UnitPrice, &pi.TotalPrice, &pi.CreatedAt)
	return pi, err
}

func (r *SupplierRepository) ListPurchaseItems(ctx context.Context, purchaseID uint64) ([]domain.PurchaseItem, error) {
	q := r.db.QueryBuilder.Select("*").From("purchase_items").Where(sq.Eq{"purchase_id": purchaseID}).OrderBy("id")
	sql, args, _ := q.ToSql()
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []domain.PurchaseItem
	for rows.Next() {
		var pi domain.PurchaseItem
		if err := rows.Scan(&pi.ID, &pi.PurchaseID, &pi.ProductID, &pi.Quantity, &pi.UnitPrice, &pi.TotalPrice, &pi.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, pi)
	}
	return items, nil
}
