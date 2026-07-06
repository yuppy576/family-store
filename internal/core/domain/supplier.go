package domain

import "time"

type Supplier struct {
	ID            uint64    `json:"id"`
	Name          string    `json:"name"`
	ContactPerson string    `json:"contact_person"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	Memo          string    `json:"memo"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	StoreID       uint64    `json:"store_id"`
}

type PurchaseStatus string

const (
	PurchasePending   PurchaseStatus = "PENDING"
	PurchaseCompleted PurchaseStatus = "COMPLETED"
	PurchaseCancelled PurchaseStatus = "CANCELLED"
)

type Purchase struct {
	ID           uint64          `json:"id"`
	SupplierID   uint64          `json:"supplier_id"`
	Operator     string          `json:"operator"`
	TotalAmount  float64         `json:"total_amount"`
	Status       PurchaseStatus  `json:"status"`
	Remark       string          `json:"remark"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	Supplier     *Supplier       `json:"supplier,omitempty"`
	StoreID      uint64          `json:"store_id"`
}

type PurchaseItem struct {
	ID         uint64    `json:"id"`
	PurchaseID uint64    `json:"purchase_id"`
	ProductID  uint64    `json:"product_id"`
	Quantity   int64     `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	Product    *Product  `json:"product,omitempty"`
	StoreID    uint64    `json:"store_id"`
}
