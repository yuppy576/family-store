package domain

import (
	"time"

	"github.com/google/uuid"
)

// Order is an entity that represents an order
type Order struct {
	ID           uint64
	UserID       uint64
	PaymentID    uint64
	CustomerName string
	TotalPrice   float64
	TotalPaid    float64
	TotalReturn  float64
	ReceiptCode  uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	StoreID      uint64
	User         *User
	Payment      *Payment
	Products     []OrderProduct
}

type SalesStats struct {
	TotalOrders  uint64  `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalPaid    float64 `json:"total_paid"`
}

type DailySales struct {
	Date         time.Time `json:"date"`
	TotalOrders  uint64    `json:"total_orders"`
	TotalRevenue float64   `json:"total_revenue"`
}
