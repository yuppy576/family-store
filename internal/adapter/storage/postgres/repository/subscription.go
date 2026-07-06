package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
)

type SubscriptionRepository struct {
	db *postgres.DB
}

func NewSubscriptionRepository(db *postgres.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db}
}

func (sr *SubscriptionRepository) CreateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	query := sr.db.QueryBuilder.Insert("subscriptions").
		Columns("store_id", "plan", "status", "start_date", "end_date").
		Values(sub.StoreID, sub.Plan, sub.Status, sub.StartDate, sub.EndDate).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&sub.ID,
		&sub.StoreID,
		&sub.Plan,
		&sub.Status,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (sr *SubscriptionRepository) GetSubscriptionByStoreID(ctx context.Context, storeID uint64) (*domain.Subscription, error) {
	var sub domain.Subscription

	query := sr.db.QueryBuilder.Select("*").
		From("subscriptions").
		Where(sq.Eq{"store_id": storeID}).
		OrderBy("created_at DESC").
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&sub.ID,
		&sub.StoreID,
		&sub.Plan,
		&sub.Status,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (sr *SubscriptionRepository) UpdateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	query := sr.db.QueryBuilder.Update("subscriptions").
		Set("plan", sub.Plan).
		Set("status", sub.Status).
		Set("start_date", sub.StartDate).
		Set("end_date", sub.EndDate).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": sub.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&sub.ID,
		&sub.StoreID,
		&sub.Plan,
		&sub.Status,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (sr *SubscriptionRepository) GetExpiringSubscriptions(ctx context.Context, days int) ([]domain.Subscription, error) {
	var sub domain.Subscription
	var subs []domain.Subscription

	now := time.Now()
	expiringDate := now.Add(time.Duration(days) * 24 * time.Hour)

	query := sr.db.QueryBuilder.Select("*").
		From("subscriptions").
		Where(sq.And{
			sq.Eq{"status": domain.SubTrial},
			sq.LtOrEq{"end_date": expiringDate},
			sq.Gt{"end_date": now},
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := sr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&sub.ID,
			&sub.StoreID,
			&sub.Plan,
			&sub.Status,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}

func (sr *SubscriptionRepository) GetExpiredSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	var sub domain.Subscription
	var subs []domain.Subscription

	now := time.Now()

	query := sr.db.QueryBuilder.Select("*").
		From("subscriptions").
		Where(sq.And{
			sq.Eq{"status": domain.SubTrial},
			sq.Lt{"end_date": now},
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := sr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&sub.ID,
			&sub.StoreID,
			&sub.Plan,
			&sub.Status,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	return subs, nil
}