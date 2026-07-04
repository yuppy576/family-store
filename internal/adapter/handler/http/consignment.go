package http

import (
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

// ConsignmentHandler handles consignment-related HTTP requests
type ConsignmentHandler struct {
	svc port.ConsignmentService
}

func NewConsignmentHandler(svc port.ConsignmentService) *ConsignmentHandler {
	return &ConsignmentHandler{svc}
}

// ── Request / Response types ─────────────────────────────────

type consignorRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	IDCard string `json:"id_card"`
	Address string `json:"address"`
	Memo   string `json:"memo"`
}

type consignmentRequest struct {
	ConsignorID      uint64  `json:"consignor_id" binding:"required"`
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	Images           []string `json:"images"`
	Category         string  `json:"category"`
	ExpectedPrice    float64 `json:"expected_price"`
	RecommendedPrice float64 `json:"recommended_price"`
	FinalPrice       float64 `json:"final_price"`
	CommissionRate   float64 `json:"commission_rate"`
	ContractEnd      string  `json:"contract_end"`  // date string
	IsVehicle        bool    `json:"is_vehicle"`
	Memo             string  `json:"memo"`
}

type vehicleRequest struct {
	VIN              string `json:"vin"`
	PlateNumber      string `json:"plate_number"`
	Brand            string `json:"brand"`
	Model            string `json:"model"`
	Year             int32  `json:"year"`
	Mileage          int32  `json:"mileage"`
	Displacement     string `json:"displacement"`
	Color            string `json:"color"`
	InspectionExpire string `json:"inspection_expire"`
	InsuranceExpire  string `json:"insurance_expire"`
}

type transferProgressRequest struct {
	Status     string `json:"status" binding:"required"`
	Remark     string `json:"remark"`
	Attachment string `json:"attachment"`
	Operator   string `json:"operator"`
}

type settlementRequest struct {
	Type             string  `json:"type" binding:"required"`
	SalePrice        float64 `json:"sale_price"`
	CommissionAmount float64 `json:"commission_amount"`
	SettlementAmount float64 `json:"settlement_amount"`
	RenewalFee       float64 `json:"renewal_fee"`
	RenewalMonths    int32   `json:"renewal_months"`
	NewEndDate       string  `json:"new_end_date"`
	Remark           string  `json:"remark"`
}

type idRequest struct {
	ID uint64 `uri:"id" binding:"required"`
}

// ── Consignor Handlers ───────────────────────────────────────

func (h *ConsignmentHandler) CreateConsignor(ctx *gin.Context) {
	var req consignorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	consignor := &domain.Consignor{
		Name: req.Name, Phone: req.Phone, IDCard: req.IDCard,
		Address: req.Address, Memo: req.Memo,
	}
	result, err := h.svc.CreateConsignor(ctx, consignor)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) GetConsignor(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	result, err := h.svc.GetConsignor(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) ListConsignors(ctx *gin.Context) {
	search := ctx.Query("q")
	skip, _ := stringToUint64(ctx.DefaultQuery("skip", "0"))
	limit, _ := stringToUint64(ctx.DefaultQuery("limit", "20"))

	consignors, err := h.svc.ListConsignors(ctx, search, skip, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}
	meta := newMeta(uint64(len(consignors)), limit, skip)
	rsp := toMap(meta, consignors, "consignors")
	handleSuccess(ctx, rsp)
}

func (h *ConsignmentHandler) UpdateConsignor(ctx *gin.Context) {
	var reqID idRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		validationError(ctx, err)
		return
	}
	var req consignorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	result, err := h.svc.UpdateConsignor(ctx, &domain.Consignor{
		ID: reqID.ID, Name: req.Name, Phone: req.Phone,
		IDCard: req.IDCard, Address: req.Address, Memo: req.Memo,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) DeleteConsignor(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := h.svc.DeleteConsignor(ctx, req.ID); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, nil)
}

// ── Consignment Handlers ─────────────────────────────────────

func (h *ConsignmentHandler) CreateConsignment(ctx *gin.Context) {
	var req consignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	cons := &domain.Consignment{
		ConsignorID: req.ConsignorID, Name: req.Name,
		Description: req.Description, Images: req.Images, Category: req.Category,
		ExpectedPrice: req.ExpectedPrice, RecommendedPrice: req.RecommendedPrice,
		FinalPrice: req.FinalPrice, CommissionRate: req.CommissionRate,
		IsVehicle: req.IsVehicle, Memo: req.Memo,
	}

	result, err := h.svc.CreateConsignment(ctx, cons)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) GetConsignment(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	result, err := h.svc.GetConsignment(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) ListConsignments(ctx *gin.Context) {
	status := ctx.Query("status")
	skip, _ := stringToUint64(ctx.DefaultQuery("skip", "0"))
	limit, _ := stringToUint64(ctx.DefaultQuery("limit", "20"))

	items, err := h.svc.ListConsignments(ctx, status, skip, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}
	meta := newMeta(uint64(len(items)), limit, skip)
	rsp := toMap(meta, items, "consignments")
	handleSuccess(ctx, rsp)
}

func (h *ConsignmentHandler) UpdateConsignment(ctx *gin.Context) {
	var reqID idRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil {
		validationError(ctx, err)
		return
	}
	var req consignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	result, err := h.svc.UpdateConsignment(ctx, &domain.Consignment{
		ID: reqID.ID, Name: req.Name, Description: req.Description,
		Category: req.Category, ExpectedPrice: req.ExpectedPrice,
		RecommendedPrice: req.RecommendedPrice, FinalPrice: req.FinalPrice,
		CommissionRate: req.CommissionRate, IsVehicle: req.IsVehicle, Memo: req.Memo,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) DeleteConsignment(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := h.svc.DeleteConsignment(ctx, req.ID); err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, nil)
}

// ── Vehicle Handlers ─────────────────────────────────────────

func (h *ConsignmentHandler) CreateVehicle(ctx *gin.Context) {
	var req vehicleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	consignmentID, _ := stringToUint64(ctx.Param("id"))

	result, err := h.svc.CreateVehicle(ctx, &domain.ConsignmentVehicle{
		ConsignmentID: consignmentID, VIN: req.VIN, PlateNumber: req.PlateNumber,
		Brand: req.Brand, Model: req.Model, Year: req.Year, Mileage: req.Mileage,
		Displacement: req.Displacement, Color: req.Color,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) GetVehicle(ctx *gin.Context) {
	consignmentID, _ := stringToUint64(ctx.Param("id"))
	result, err := h.svc.GetVehicle(ctx, consignmentID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) UpdateVehicle(ctx *gin.Context) {
	var req vehicleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	consignmentID, _ := stringToUint64(ctx.Param("id"))

	result, err := h.svc.UpdateVehicle(ctx, &domain.ConsignmentVehicle{
		ConsignmentID: consignmentID, VIN: req.VIN, PlateNumber: req.PlateNumber,
		Brand: req.Brand, Model: req.Model, Year: req.Year, Mileage: req.Mileage,
		Displacement: req.Displacement, Color: req.Color,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

// ── Transfer Progress Handlers ───────────────────────────────

func (h *ConsignmentHandler) CreateTransferProgress(ctx *gin.Context) {
	var req transferProgressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	vehicleID, _ := stringToUint64(ctx.Param("id"))

	result, err := h.svc.CreateTransferProgress(ctx, &domain.TransferProgress{
		VehicleID: vehicleID, Status: domain.TransferStatus(req.Status),
		Remark: req.Remark, Attachment: req.Attachment, Operator: req.Operator,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) ListTransferProgress(ctx *gin.Context) {
	vehicleID, _ := stringToUint64(ctx.Param("id"))
	result, err := h.svc.ListTransferProgress(ctx, vehicleID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

// ── Settlement Handlers ──────────────────────────────────────

func (h *ConsignmentHandler) CreateSettlement(ctx *gin.Context) {
	var req settlementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	consignmentID, _ := stringToUint64(ctx.Param("id"))

	result, err := h.svc.CreateSettlement(ctx, &domain.ConsignmentSettlement{
		ConsignmentID: consignmentID, Type: domain.SettlementType(req.Type),
		SalePrice: req.SalePrice, CommissionAmount: req.CommissionAmount,
		SettlementAmount: req.SettlementAmount,
	})
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}

func (h *ConsignmentHandler) ListSettlements(ctx *gin.Context) {
	consignmentID, _ := stringToUint64(ctx.Param("id"))
	result, err := h.svc.ListSettlements(ctx, consignmentID)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, result)
}
