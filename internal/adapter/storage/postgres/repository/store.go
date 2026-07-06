package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type StoreRepository struct {
	db *postgres.DB
}

func NewStoreRepository(db *postgres.DB) *StoreRepository {
	return &StoreRepository{db}
}

func (sr *StoreRepository) CreateStore(ctx context.Context, store *domain.Store) (*domain.Store, error) {
	query := sr.db.QueryBuilder.Insert("stores").
		Columns("name", "domain", "status", "trial_end").
		Values(store.Name, store.Domain, store.Status, store.TrialEnd).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&store.ID,
		&store.Name,
		&store.Domain,
		&store.Status,
		&store.TrialEnd,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if errCode := sr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return store, nil
}

func (sr *StoreRepository) GetStoreByDomain(ctx context.Context, domainName string) (*domain.Store, error) {
	var store domain.Store

	query := sr.db.QueryBuilder.Select("*").
		From("stores").
		Where(sq.Eq{"domain": domainName}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&store.ID,
		&store.Name,
		&store.Domain,
		&store.Status,
		&store.TrialEnd,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &store, nil
}

func (sr *StoreRepository) GetStoreByID(ctx context.Context, id uint64) (*domain.Store, error) {
	var store domain.Store

	query := sr.db.QueryBuilder.Select("*").
		From("stores").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&store.ID,
		&store.Name,
		&store.Domain,
		&store.Status,
		&store.TrialEnd,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &store, nil
}

func (sr *StoreRepository) GetAllStores(ctx context.Context) ([]domain.Store, error) {
	var stores []domain.Store

	query := sr.db.QueryBuilder.Select("*").
		From("stores").
		OrderBy("id DESC")

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
		var store domain.Store
		err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Domain,
			&store.Status,
			&store.TrialEnd,
			&store.CreatedAt,
			&store.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	return stores, nil
}

func (sr *StoreRepository) UpdateStore(ctx context.Context, store *domain.Store) (*domain.Store, error) {
	query := sr.db.QueryBuilder.Update("stores").
		Set("name", store.Name).
		Set("status", store.Status).
		Set("trial_end", store.TrialEnd).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": store.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = sr.db.QueryRow(ctx, sql, args...).Scan(
		&store.ID,
		&store.Name,
		&store.Domain,
		&store.Status,
		&store.TrialEnd,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if errCode := sr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return store, nil
}