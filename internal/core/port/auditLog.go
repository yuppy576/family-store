package port

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
)

type AuditLogRepository interface {
	CreateAuditLog(ctx context.Context, log *domain.AuditLog) (*domain.AuditLog, error)
	ListAuditLogs(ctx context.Context, params map[string]interface{}, skip, limit uint64) ([]domain.AuditLog, error)
	CountAuditLogs(ctx context.Context, params map[string]interface{}) (uint64, error)
}

type AuditLogService interface {
	CreateAuditLog(ctx context.Context, log *domain.AuditLog) (*domain.AuditLog, error)
	ListAuditLogs(ctx context.Context, params map[string]interface{}, skip, limit uint64) ([]domain.AuditLog, error)
	CountAuditLogs(ctx context.Context, params map[string]interface{}) (uint64, error)
}
