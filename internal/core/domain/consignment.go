package domain

import "time"

// ConsignmentStatus represents the status of a consignment item
type ConsignmentStatus string

const (
	ConsignmentStatusOnSale    ConsignmentStatus = "ON_SALE"
	ConsignmentStatusSold      ConsignmentStatus = "SOLD"
	ConsignmentStatusExpired   ConsignmentStatus = "EXPIRED"
	ConsignmentStatusReturned  ConsignmentStatus = "RETURNED"
	ConsignmentStatusCancelled ConsignmentStatus = "CANCELLED"
)

// TransferStatus represents the vehicle transfer progress status
type TransferStatus string

const (
	TransferPendingInsp  TransferStatus = "PENDING_INSPECTION"
	TransferInspected    TransferStatus = "INSPECTED"
	TransferTransferring TransferStatus = "TRANSFERRING"
	TransferTransferred  TransferStatus = "TRANSFERRED"
	TransferSettled      TransferStatus = "SETTLED"
)

// SettlementType represents the type of settlement
type SettlementType string

const (
	SettlementSold   SettlementType = "SOLD_SETTLEMENT"
	SettlementReturn SettlementType = "RETURN_SETTLEMENT"
	SettlementRenew  SettlementType = "RENEWAL"
)

// Consignor represents a person who consigns items for sale
type Consignor struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	IDCard    string    `json:"id_card"`
	Address   string    `json:"address"`
	Memo      string    `json:"memo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Consignment represents an item placed on consignment
type Consignment struct {
	ID               uint64            `json:"id"`
	ConsignorID      uint64            `json:"consignor_id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Images           []string          `json:"images"`
	Category         string            `json:"category"`
	ExpectedPrice    float64           `json:"expected_price"`
	RecommendedPrice float64           `json:"recommended_price"`
	FinalPrice       float64           `json:"final_price"`
	CommissionRate   float64           `json:"commission_rate"`
	CommissionAmount float64           `json:"commission_amount"`
	Status           ConsignmentStatus `json:"status"`
	ContractEnd      time.Time         `json:"contract_end"`
	IsVehicle        bool              `json:"is_vehicle"`
	Memo             string            `json:"memo"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Consignor        *Consignor        `json:"consignor,omitempty"`
}

// ConsignmentVehicle stores vehicle-specific consignment details
type ConsignmentVehicle struct {
	ID               uint64    `json:"id"`
	ConsignmentID    uint64    `json:"consignment_id"`
	VIN              string    `json:"vin"`
	PlateNumber      string    `json:"plate_number"`
	Brand            string    `json:"brand"`
	Model            string    `json:"model"`
	Year             int32     `json:"year"`
	Mileage          int32     `json:"mileage"`
	Displacement     string    `json:"displacement"`
	Color            string    `json:"color"`
	InspectionExpire time.Time `json:"inspection_expire"`
	InsuranceExpire  time.Time `json:"insurance_expire"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TransferProgress tracks the vehicle transfer process
type TransferProgress struct {
	ID          uint64         `json:"id"`
	VehicleID   uint64         `json:"vehicle_id"`
	Status      TransferStatus `json:"status"`
	Remark      string         `json:"remark"`
	Attachment  string         `json:"attachment"`
	Operator    string         `json:"operator"`
	CreatedAt   time.Time      `json:"created_at"`
}

// ConsignmentSettlement records financial settlements
type ConsignmentSettlement struct {
	ID                uint64         `json:"id"`
	ConsignmentID     uint64         `json:"consignment_id"`
	Type              SettlementType `json:"type"`
	SalePrice         float64        `json:"sale_price"`
	CommissionAmount  float64        `json:"commission_amount"`
	SettlementAmount  float64        `json:"settlement_amount"`
	RenewalFee        float64        `json:"renewal_fee"`
	RenewalMonths     int32          `json:"renewal_months"`
	NewEndDate        time.Time      `json:"new_end_date"`
	Remark            string         `json:"remark"`
	CreatedAt         time.Time      `json:"created_at"`
}
