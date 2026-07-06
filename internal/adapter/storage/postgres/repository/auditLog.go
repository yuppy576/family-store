package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type AuditLogRepository struct {
	db *postgres.DB
}

func NewAuditLogRepository(db *postgres.DB) *AuditLogRepository {
	return &AuditLogRepository{db}
}

func (ar *AuditLogRepository) CreateAuditLog(ctx context.Context, log *domain.AuditLog) (*domain.AuditLog, error) {
	storeID := getStoreIDFromContext(ctx)
	query := ar.db.QueryBuilder.Insert("audit_logs").
		Columns("user_id", "user_name", "action", "resource_type", "resource_id", "old_data", "new_data", "ip_address", "store_id").
		Values(log.UserID, log.UserName, log.Action, log.ResourceType, log.ResourceID, log.OldData, log.NewData, log.IPAddress, storeID).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ar.db.QueryRow(ctx, sql, args...).Scan(
		&log.ID,
		&log.UserID,
		&log.UserName,
		&log.Action,
		&log.ResourceType,
		&log.ResourceID,
		&log.OldData,
		&log.NewData,
		&log.IPAddress,
		&log.CreatedAt,
		&log.StoreID,
	)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (ar *AuditLogRepository) ListAuditLogs(ctx context.Context, params map[string]interface{}, skip, limit uint64) ([]domain.AuditLog, error) {
	var log domain.AuditLog
	var logs []domain.AuditLog
	storeID := getStoreIDFromContext(ctx)

	query := ar.db.QueryBuilder.Select("*").
		From("audit_logs").
		Where(sq.Eq{"store_id": storeID}).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(skip)

	if userID, ok := params["user_id"].(uint64); ok && userID > 0 {
		query = query.Where(sq.Eq{"user_id": userID})
	}
	if action, ok := params["action"].(string); ok && action != "" {
		query = query.Where(sq.Eq{"action": action})
	}
	if resourceType, ok := params["resource_type"].(string); ok && resourceType != "" {
		query = query.Where(sq.Eq{"resource_type": resourceType})
	}
	if resourceID, ok := params["resource_id"].(uint64); ok && resourceID > 0 {
		query = query.Where(sq.Eq{"resource_id": resourceID})
	}
	if startTime, ok := params["start_time"].(time.Time); ok {
		query = query.Where(sq.GtOrEq{"created_at": startTime})
	}
	if endTime, ok := params["end_time"].(time.Time); ok {
		query = query.Where(sq.LtOrEq{"created_at": endTime})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ar.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.UserName,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.OldData,
			&log.NewData,
			&log.IPAddress,
			&log.CreatedAt,
			&log.StoreID,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (ar *AuditLogRepository) CountAuditLogs(ctx context.Context, params map[string]interface{}) (uint64, error) {
	var count uint64
	storeID := getStoreIDFromContext(ctx)

	query := ar.db.QueryBuilder.Select("COUNT(*)").From("audit_logs").Where(sq.Eq{"store_id": storeID})

	if userID, ok := params["user_id"].(uint64); ok && userID > 0 {
		query = query.Where(sq.Eq{"user_id": userID})
	}
	if action, ok := params["action"].(string); ok && action != "" {
		query = query.Where(sq.Eq{"action": action})
	}
	if resourceType, ok := params["resource_type"].(string); ok && resourceType != "" {
		query = query.Where(sq.Eq{"resource_type": resourceType})
	}
	if resourceID, ok := params["resource_id"].(uint64); ok && resourceID > 0 {
		query = query.Where(sq.Eq{"resource_id": resourceID})
	}
	if startTime, ok := params["start_time"].(time.Time); ok {
		query = query.Where(sq.GtOrEq{"created_at": startTime})
	}
	if endTime, ok := params["end_time"].(time.Time); ok {
		query = query.Where(sq.LtOrEq{"created_at": endTime})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	err = ar.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return count, nil
}