package service

import (
	"context"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
)

type AuditLogService struct {
	repo  port.AuditLogRepository
	cache port.CacheRepository
}

func NewAuditLogService(repo port.AuditLogRepository, cache port.CacheRepository) *AuditLogService {
	return &AuditLogService{repo, cache}
}

func (as *AuditLogService) CreateAuditLog(ctx context.Context, log *domain.AuditLog) (*domain.AuditLog, error) {
	return as.repo.CreateAuditLog(ctx, log)
}

func (as *AuditLogService) ListAuditLogs(ctx context.Context, params map[string]interface{}, skip, limit uint64) ([]domain.AuditLog, error) {
	return as.repo.ListAuditLogs(ctx, params, skip, limit)
}

func (as *AuditLogService) CountAuditLogs(ctx context.Context, params map[string]interface{}) (uint64, error) {
	return as.repo.CountAuditLogs(ctx, params)
}
