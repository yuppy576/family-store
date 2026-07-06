package http

import (
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	svc port.SubscriptionService
}

func NewSubscriptionHandler(svc port.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc}
}

type renewSubscriptionRequest struct {
	Plan   string `json:"plan" binding:"required" example:"PERSONAL"`
	Months int    `json:"months" binding:"required,min=1" example:"12"`
}

func (sh *SubscriptionHandler) RenewSubscription(ctx *gin.Context) {
	var req renewSubscriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	storeID := getStoreIDFromContext(ctx)

	plan := domain.SubscriptionPlan(req.Plan)
	if plan != domain.PlanPersonal && plan != domain.PlanProfessional && plan != domain.PlanLifetime {
		handleError(ctx, domain.ErrInvalidSubscriptionPlan)
		return
	}

	sub, err := sh.svc.RenewSubscription(ctx, storeID, plan, req.Months)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, newSubscriptionResponse(sub))
}

func (sh *SubscriptionHandler) GetSubscription(ctx *gin.Context) {
	storeID := getStoreIDFromContext(ctx)

	sub, err := sh.svc.GetSubscription(ctx, storeID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, newSubscriptionResponse(sub))
}

func (sh *SubscriptionHandler) ActivateSubscription(ctx *gin.Context) {
	storeID := getStoreIDFromContext(ctx)

	sub, err := sh.svc.ActivateSubscription(ctx, storeID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, newSubscriptionResponse(sub))
}

type subscriptionResponse struct {
	ID        uint64            `json:"id"`
	StoreID   uint64            `json:"store_id"`
	Plan      string            `json:"plan"`
	Status    string            `json:"status"`
	StartDate string            `json:"start_date"`
	EndDate   *string           `json:"end_date"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

func newSubscriptionResponse(sub *domain.Subscription) subscriptionResponse {
	var endDate *string
	if sub.EndDate != nil {
		ed := sub.EndDate.Format("2006-01-02 15:04:05")
		endDate = &ed
	}

	return subscriptionResponse{
		ID:        sub.ID,
		StoreID:   sub.StoreID,
		Plan:      string(sub.Plan),
		Status:    string(sub.Status),
		StartDate: sub.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:   endDate,
		CreatedAt: sub.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: sub.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func getStoreIDFromContext(ctx *gin.Context) uint64 {
	val, _ := ctx.Get(storeIDKey)
	if val == nil {
		return 1
	}
	switch v := val.(type) {
	case uint64:
		return v
	case int64:
		return uint64(v)
	case int:
		return uint64(v)
	default:
		return 1
	}
}