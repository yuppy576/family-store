package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	db *postgres.DB
}

func NewCategoryRepository(db *postgres.DB) *CategoryRepository {
	return &CategoryRepository{db}
}

func (cr *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	storeID := getStoreIDFromContext(ctx)
	query := cr.db.QueryBuilder.Insert("categories").
		Columns("name", "store_id").
		Values(category.Name, storeID).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.StoreID,
	)
	if err != nil {
		if errCode := cr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return category, nil
}

func (cr *CategoryRepository) GetCategoryByID(ctx context.Context, id uint64) (*domain.Category, error) {
	var category domain.Category
	storeID := getStoreIDFromContext(ctx)

	query := cr.db.QueryBuilder.Select("*").
		From("categories").
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"store_id": storeID}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.StoreID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (cr *CategoryRepository) ListCategories(ctx context.Context, skip, limit uint64) ([]domain.Category, error) {
	var category domain.Category
	var categories []domain.Category
	storeID := getStoreIDFromContext(ctx)

	query := cr.db.QueryBuilder.Select("*").
		From("categories").
		Where(sq.Eq{"store_id": storeID}).
		OrderBy("id").
		Limit(limit).
		Offset(skip * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.StoreID,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (cr *CategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	storeID := getStoreIDFromContext(ctx)
	query := cr.db.QueryBuilder.Update("categories").
		Set("name", category.Name).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": category.ID}).
		Where(sq.Eq{"store_id": storeID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.StoreID,
	)
	if err != nil {
		if errCode := cr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return category, nil
}

func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id uint64) error {
	storeID := getStoreIDFromContext(ctx)
	query := cr.db.QueryBuilder.Delete("categories").
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"store_id": storeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = cr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}