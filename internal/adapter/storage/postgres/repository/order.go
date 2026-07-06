package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

/**
 * OrderRepository implements port.OrderRepository interface
 * and provides an access to the postgres database
 */
type OrderRepository struct {
	db *postgres.DB
}

// NewOrderRepository creates a new order repository instance
func NewOrderRepository(db *postgres.DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}

// CreateOrder creates a new order in the database
func (or *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	var product domain.Product
	var products []domain.OrderProduct
	storeID := getStoreIDFromContext(ctx)

	orderQuery := or.db.QueryBuilder.Insert("orders").
		Columns("user_id", "payment_id", "customer_name", "total_price", "total_paid", "total_return", "store_id").
		Values(order.UserID, order.PaymentID, order.CustomerName, order.TotalPrice, order.TotalPaid, order.TotalReturn, storeID).
		Suffix("RETURNING *")

	err := pgx.BeginFunc(ctx, or.db, func(tx pgx.Tx) error {
		sql, args, err := orderQuery.ToSql()
		if err != nil {
			return err
		}

		err = tx.QueryRow(ctx, sql, args...).Scan(
			&order.ID,
			&order.UserID,
			&order.PaymentID,
			&order.CustomerName,
			&order.TotalPrice,
			&order.TotalPaid,
			&order.TotalReturn,
			&order.ReceiptCode,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.StoreID,
		)
		if err != nil {
			return err
		}

		for _, orderProduct := range order.Products {
			orderProductQuery := or.db.QueryBuilder.Insert("order_products").
				Columns("order_id", "product_id", "quantity", "total_price", "store_id").
				Values(order.ID, orderProduct.ProductID, orderProduct.Quantity, orderProduct.TotalPrice, storeID).
				Suffix("RETURNING *")

			sql, args, err := orderProductQuery.ToSql()
			if err != nil {
				return err
			}

			err = tx.QueryRow(ctx, sql, args...).Scan(
				&orderProduct.ID,
				&orderProduct.OrderID,
				&orderProduct.ProductID,
				&orderProduct.Quantity,
				&orderProduct.TotalPrice,
				&orderProduct.CreatedAt,
				&orderProduct.UpdatedAt,
				&orderProduct.StoreID,
			)
			if err != nil {
				return err
			}

			products = append(products, orderProduct)

			productQuery := or.db.QueryBuilder.Update("products").
				Set("stock", sq.Expr("stock - ?", orderProduct.Quantity)).
				Set("updated_at", time.Now()).
				Where(sq.Eq{"id": orderProduct.ProductID}).
				Suffix("RETURNING stock")

			sql, args, err = productQuery.ToSql()
			if err != nil {
				return err
			}

			err = tx.QueryRow(ctx, sql, args...).Scan(
				&product.Stock,
			)
			if err != nil {
				return err
			}

			if product.Stock < 0 {
				return tx.Rollback(ctx)
			}
		}

		order.Products = products

		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, err
}

// GetOrderByID gets an order by ID from the database
func (or *OrderRepository) GetOrderByID(ctx context.Context, id uint64) (*domain.Order, error) {
	var order domain.Order
	var orderProduct domain.OrderProduct
	storeID := getStoreIDFromContext(ctx)

	orderQuery := or.db.QueryBuilder.Select("*").
		From("orders").
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"store_id": storeID}).
		Limit(1)

	orderProductQuery := or.db.QueryBuilder.Select("*").
		From("order_products").
		Where(sq.Eq{"order_id": id}).
		Where(sq.Eq{"store_id": storeID})

	err := pgx.BeginFunc(ctx, or.db, func(tx pgx.Tx) error {

		sql, args, err := orderQuery.ToSql()
		if err != nil {
			return err
		}

		err = tx.QueryRow(ctx, sql, args...).Scan(
			&order.ID,
			&order.UserID,
			&order.PaymentID,
			&order.CustomerName,
			&order.TotalPrice,
			&order.TotalPaid,
			&order.TotalReturn,
			&order.ReceiptCode,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.StoreID,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return domain.ErrDataNotFound
			}
			return err
		}

		sql, args, err = orderProductQuery.ToSql()
		if err != nil {
			return err
		}

		rows, err := tx.Query(ctx, sql, args...)
		if err != nil {
			return err
		}

		for rows.Next() {
			err = rows.Scan(
				&orderProduct.ID,
				&orderProduct.OrderID,
				&orderProduct.ProductID,
				&orderProduct.Quantity,
				&orderProduct.TotalPrice,
				&orderProduct.CreatedAt,
				&orderProduct.UpdatedAt,
				&orderProduct.StoreID,
			)
			if err != nil {
				return err
			}

			order.Products = append(order.Products, orderProduct)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// ListOrders lists all orders from the database
func (or *OrderRepository) ListOrders(ctx context.Context, skip, limit uint64) ([]domain.Order, error) {
	var order domain.Order
	var orderProduct domain.OrderProduct
	var orders []domain.Order
	storeID := getStoreIDFromContext(ctx)

	ordersQuery := or.db.QueryBuilder.Select("*").
		From("orders").
		Where(sq.Eq{"store_id": storeID}).
		OrderBy("id").
		Limit(limit).
		Offset(skip * limit)

	err := pgx.BeginFunc(ctx, or.db, func(tx pgx.Tx) error {
		sql, args, err := ordersQuery.ToSql()
		if err != nil {
			return err
		}

		rows, err := tx.Query(ctx, sql, args...)
		if err != nil {
			return err
		}

		for rows.Next() {
			err := rows.Scan(
				&order.ID,
				&order.UserID,
				&order.PaymentID,
				&order.CustomerName,
				&order.TotalPrice,
				&order.TotalPaid,
				&order.TotalReturn,
				&order.ReceiptCode,
				&order.CreatedAt,
				&order.UpdatedAt,
				&order.StoreID,
			)
			if err != nil {
				return err
			}

			orders = append(orders, order)
		}

		for i, order := range orders {
			orderProductQuery := or.db.QueryBuilder.Select("*").
				From("order_products").
				Where(sq.Eq{"order_id": order.ID}).
				Where(sq.Eq{"store_id": storeID})

			sql, args, err := orderProductQuery.ToSql()
			if err != nil {
				return err
			}

			rows, err := tx.Query(ctx, sql, args...)
			if err != nil {
				return err
			}

			for rows.Next() {
				err := rows.Scan(
					&orderProduct.ID,
					&orderProduct.OrderID,
					&orderProduct.ProductID,
					&orderProduct.Quantity,
					&orderProduct.TotalPrice,
					&orderProduct.CreatedAt,
					&orderProduct.UpdatedAt,
					&orderProduct.StoreID,
				)
				if err != nil {
					return err
				}

				orders[i].Products = append(orders[i].Products, orderProduct)
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (or *OrderRepository) GetSalesStats(ctx context.Context, startDate, endDate time.Time) (*domain.SalesStats, error) {
	var stats domain.SalesStats
	storeID := getStoreIDFromContext(ctx)

	query := or.db.QueryBuilder.Select(
		"COUNT(*) as total_orders",
		"SUM(total_price) as total_revenue",
		"SUM(total_paid) as total_paid",
	).From("orders").
		Where(sq.Eq{"store_id": storeID}).
		Where(sq.GtOrEq{"created_at": startDate}).
		Where(sq.LtOrEq{"created_at": endDate})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = or.db.QueryRow(ctx, sql, args...).Scan(
		&stats.TotalOrders,
		&stats.TotalRevenue,
		&stats.TotalPaid,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &domain.SalesStats{}, nil
		}
		return nil, err
	}

	return &stats, nil
}

func (or *OrderRepository) GetDailySales(ctx context.Context, startDate, endDate time.Time) ([]domain.DailySales, error) {
	var sales []domain.DailySales
	storeID := getStoreIDFromContext(ctx)

	query := or.db.QueryBuilder.Select(
		"DATE(created_at) as date",
		"COUNT(*) as total_orders",
		"SUM(total_price) as total_revenue",
	).From("orders").
		Where(sq.Eq{"store_id": storeID}).
		Where(sq.GtOrEq{"created_at": startDate}).
		Where(sq.LtOrEq{"created_at": endDate}).
		GroupBy("DATE(created_at)").
		OrderBy("DATE(created_at)")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := or.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.DailySales
		err := rows.Scan(&s.Date, &s.TotalOrders, &s.TotalRevenue)
		if err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}

	return sales, nil
}
