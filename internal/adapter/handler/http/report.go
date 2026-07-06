package http

import (
	"time"

	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	orderSvc port.OrderService
}

func NewReportHandler(orderSvc port.OrderService) *ReportHandler {
	return &ReportHandler{orderSvc}
}

type salesStatsRequest struct {
	StartDate string `form:"start_date" example:"2024-01-01"`
	EndDate   string `form:"end_date" example:"2024-01-31"`
}

func (rh *ReportHandler) GetSalesStats(ctx *gin.Context) {
	var req salesStatsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		endDate = time.Now()
	}

	stats, err := rh.orderSvc.GetSalesStats(ctx, startDate, endDate)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, stats)
}

func (rh *ReportHandler) GetDailySales(ctx *gin.Context) {
	var req salesStatsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		startDate = time.Now().AddDate(0, 0, -30)
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		endDate = time.Now()
	}

	sales, err := rh.orderSvc.GetDailySales(ctx, startDate, endDate)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, sales)
}