package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bagashiz/go-pos/internal/adapter/storage/postgres"
	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

/**
 * UserRepository implements port.UserRepository interface
 * and provides an access to the postgres database
 */
type UserRepository struct {
	db *postgres.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db,
	}
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	storeID := getStoreIDFromContext(ctx)
	query := ur.db.QueryBuilder.Insert("users").
		Columns("name", "email", "password", "store_id").
		Values(user.Name, user.Email, user.Password, storeID).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.StoreID,
	)
	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID gets a user by ID from the database
func (ur *UserRepository) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User
	storeID := getStoreIDFromContext(ctx)

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"store_id": storeID}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.StoreID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail gets a user by email from the database
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	storeID := getStoreIDFromContext(ctx)

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"email": email})

	if storeID > 0 {
		query = query.Where(sq.Eq{"store_id": storeID})
	}

	query = query.Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.StoreID,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &user, nil
}

// ListUsers lists all users from the database
func (ur *UserRepository) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var user domain.User
	var users []domain.User
	storeID := getStoreIDFromContext(ctx)

	query := ur.db.QueryBuilder.Select("*").
		From("users").
		Where(sq.Eq{"store_id": storeID}).
		OrderBy("id").
		Limit(limit).
		Offset(skip * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := ur.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.StoreID,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates a user by ID in the database
func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	name := nullString(user.Name)
	email := nullString(user.Email)
	password := nullString(user.Password)
	role := nullString(string(user.Role))
	storeID := getStoreIDFromContext(ctx)

	query := ur.db.QueryBuilder.Update("users").
		Set("name", sq.Expr("COALESCE(?, name)", name)).
		Set("email", sq.Expr("COALESCE(?, email)", email)).
		Set("password", sq.Expr("COALESCE(?, password)", password)).
		Set("role", sq.Expr("COALESCE(?, role)", role)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": user.ID}).
		Where(sq.Eq{"store_id": storeID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = ur.db.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.StoreID,
	)
	if err != nil {
		if errCode := ur.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID from the database
func (ur *UserRepository) DeleteUser(ctx context.Context, id uint64) error {
	storeID := getStoreIDFromContext(ctx)
	query := ur.db.QueryBuilder.Delete("users").
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"store_id": storeID})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = ur.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
