package http

import (
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct{ svc port.SupplierService }

func NewPurchaseHandler(svc port.SupplierService) *PurchaseHandler { return &PurchaseHandler{svc} }

type purchaseItemReq struct {
	ProductID uint64  `json:"product_id" binding:"required"`
	Quantity  int64   `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required,min=0"`
}

type purchaseReq struct {
	SupplierID uint64            `json:"supplier_id" binding:"required"`
	Operator   string            `json:"operator"`
	Remark     string            `json:"remark"`
	Items      []purchaseItemReq `json:"items" binding:"required,min=1"`
}

func (h *PurchaseHandler) Create(ctx *gin.Context) {
	var req purchaseReq
	if err := ctx.ShouldBindJSON(&req); err != nil { validationError(ctx, err); return }

	var items []domain.PurchaseItem
	var total float64
	for _, it := range req.Items {
		t := float64(it.Quantity) * it.UnitPrice
		items = append(items, domain.PurchaseItem{ProductID: it.ProductID, Quantity: it.Quantity, UnitPrice: it.UnitPrice, TotalPrice: t})
		total += t
	}

	operator := req.Operator
	if operator == "" {
		operator = "admin"
	}

	purchase := &domain.Purchase{SupplierID: req.SupplierID, Operator: operator, TotalAmount: total, Remark: req.Remark}
	r, err := h.svc.CreatePurchase(ctx, purchase, items)
	if err != nil { handleError(ctx, err); return }
	handleSuccess(ctx, r)
}

func (h *PurchaseHandler) Get(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil { validationError(ctx, err); return }
	r, err := h.svc.GetPurchase(ctx, req.ID)
	if err != nil { handleError(ctx, err); return }
	handleSuccess(ctx, r)
}

func (h *PurchaseHandler) List(ctx *gin.Context) {
	skip, _ := stringToUint64(ctx.DefaultQuery("skip", "0"))
	limit, _ := stringToUint64(ctx.DefaultQuery("limit", "20"))
	r, err := h.svc.ListPurchases(ctx, skip, limit)
	if err != nil { handleError(ctx, err); return }
	m := newMeta(uint64(len(r)), limit, skip)
	handleSuccess(ctx, toMap(m, r, "purchases"))
}
