package http

import (
	"strconv"
	"time"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

type AuditLogHandler struct {
	svc port.AuditLogService
}

func NewAuditLogHandler(svc port.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{svc}
}

type listAuditLogsRequest struct {
	Skip          uint64 `form:"skip"`
	Limit         uint64 `form:"limit"`
	Page          uint64 `form:"page"`
	PageSize      uint64 `form:"page_size"`
	UserID        uint64 `form:"user_id"`
	Action        string `form:"action"`
	ResourceType  string `form:"resource_type"`
	ResourceID    uint64 `form:"resource_id"`
	StartTime     string `form:"start_time"`
	EndTime       string `form:"end_time"`
}

func (ah *AuditLogHandler) ListAuditLogs(ctx *gin.Context) {
	var req listAuditLogsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	skip := req.Skip
	limit := req.Limit

	if req.Page > 0 && req.PageSize > 0 {
		skip = (req.Page - 1) * req.PageSize
		limit = req.PageSize
	}

	if limit == 0 {
		limit = 20
	}

	params := make(map[string]interface{})

	if req.UserID > 0 {
		params["user_id"] = req.UserID
	}
	if req.Action != "" {
		params["action"] = req.Action
	}
	if req.ResourceType != "" {
		params["resource_type"] = req.ResourceType
	}
	if req.ResourceID > 0 {
		params["resource_id"] = req.ResourceID
	}
	if req.StartTime != "" {
		if t, err := time.Parse("2006-01-02", req.StartTime); err == nil {
			params["start_time"] = t
		}
	}
	if req.EndTime != "" {
		if t, err := time.Parse("2006-01-02", req.EndTime); err == nil {
			params["end_time"] = t.Add(24 * time.Hour)
		}
	}

	logs, err := ah.svc.ListAuditLogs(ctx, params, skip, limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	count, err := ah.svc.CountAuditLogs(ctx, params)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := make([]auditLogResponse, 0, len(logs))
	for _, log := range logs {
		rsp = append(rsp, newAuditLogResponse(&log))
	}

	handleSuccess(ctx, toMap(newMeta(count, limit, skip), rsp, "audit_logs"))
}

type auditLogResponse struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	UserName     string `json:"user_name"`
	Action       string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID   uint64 `json:"resource_id"`
	OldData      string `json:"old_data"`
	NewData      string `json:"new_data"`
	IPAddress    string `json:"ip_address"`
	CreatedAt    string `json:"created_at"`
}

func newAuditLogResponse(log *domain.AuditLog) auditLogResponse {
	return auditLogResponse{
		ID:           log.ID,
		UserID:       log.UserID,
		UserName:     log.UserName,
		Action:       string(log.Action),
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		OldData:      string(log.OldData),
		NewData:      string(log.NewData),
		IPAddress:    log.IPAddress,
		CreatedAt:    log.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func getResourceType(path string) string {
	switch {
	case contains(path, "/users/"):
		return "users"
	case contains(path, "/payments/"):
		return "payments"
	case contains(path, "/categories/"):
		return "categories"
	case contains(path, "/products/"):
		return "products"
	case contains(path, "/orders/"):
		return "orders"
	case contains(path, "/consignment/consignors"):
		return "consignors"
	case contains(path, "/consignment/items"):
		return "consignments"
	case contains(path, "/consignment/settlements"):
		return "settlements"
	case contains(path, "/suppliers/"):
		return "suppliers"
	case contains(path, "/purchases/"):
		return "purchases"
	default:
		return "unknown"
	}
}

func getResourceID(path string) uint64 {
	parts := splitPath(path)
	for i, part := range parts {
		if i > 0 && isNumeric(part) {
			id, _ := strconv.ParseUint(part, 10, 64)
			return id
		}
	}
	return 0
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(s)] != s && s[0:len(substr)] == substr || len(s) > len(substr) && containsHelper(s, substr)
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitPath(path string) []string {
	result := []string{}
	current := ""
	for _, c := range path {
		if c == '/' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
