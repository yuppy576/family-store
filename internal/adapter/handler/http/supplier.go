package http

import (
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct{ svc port.SupplierService }

func NewSupplierHandler(svc port.SupplierService) *SupplierHandler { return &SupplierHandler{svc} }

type supplierReq struct {
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contact_person"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	Memo          string `json:"memo"`
}

func (h *SupplierHandler) Create(ctx *gin.Context) {
	var req supplierReq
	if err := ctx.ShouldBindJSON(&req); err != nil { validationError(ctx, err); return }
	r, err := h.svc.CreateSupplier(ctx, &domain.Supplier{Name: req.Name, ContactPerson: req.ContactPerson, Phone: req.Phone, Address: req.Address, Memo: req.Memo})
	if err != nil { handleError(ctx, err); return }
	handleSuccess(ctx, r)
}

func (h *SupplierHandler) Get(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil { validationError(ctx, err); return }
	r, err := h.svc.GetSupplier(ctx, req.ID)
	if err != nil { handleError(ctx, err); return }
	handleSuccess(ctx, r)
}

func (h *SupplierHandler) List(ctx *gin.Context) {
	search := ctx.Query("q")
	skip, _ := stringToUint64(ctx.DefaultQuery("skip", "0"))
	limit, _ := stringToUint64(ctx.DefaultQuery("limit", "20"))
	r, err := h.svc.ListSuppliers(ctx, search, skip, limit)
	if err != nil { handleError(ctx, err); return }
	m := newMeta(uint64(len(r)), limit, skip)
	handleSuccess(ctx, toMap(m, r, "suppliers"))
}

func (h *SupplierHandler) Update(ctx *gin.Context) {
	var reqID idRequest
	if err := ctx.ShouldBindUri(&reqID); err != nil { validationError(ctx, err); return }
	var req supplierReq
	if err := ctx.ShouldBindJSON(&req); err != nil { validationError(ctx, err); return }
	r, err := h.svc.UpdateSupplier(ctx, &domain.Supplier{ID: reqID.ID, Name: req.Name, ContactPerson: req.ContactPerson, Phone: req.Phone, Address: req.Address, Memo: req.Memo})
	if err != nil { handleError(ctx, err); return }
	handleSuccess(ctx, r)
}

func (h *SupplierHandler) Delete(ctx *gin.Context) {
	var req idRequest
	if err := ctx.ShouldBindUri(&req); err != nil { validationError(ctx, err); return }
	if err := h.svc.DeleteSupplier(ctx, req.ID); err != nil { handleError(ctx, err); return }
	handleSuccess(ctx, nil)
}
