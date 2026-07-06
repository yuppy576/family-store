package repository

import (
	"context"
	"database/sql"
)

// nullString converts a string to sql.NullString for empty string check
func nullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}

	return sql.NullString{
		String: value,
		Valid:  true,
	}
}

// nullUint64 converts an uint64 to sql.NullInt64 for empty uint64 check
func nullUint64(value uint64) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{}
	}

	valueInt64 := int64(value)

	return sql.NullInt64{
		Int64: valueInt64,
		Valid: true,
	}
}

// nullInt64 converts an int64 to sql.NullInt64 for empty int64 check
func nullInt64(value int64) sql.NullInt64 {
	if value == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: value,
		Valid: true,
	}
}

// nullFloat64 converts a float64 to sql.NullFloat64 for empty float64 check
func nullFloat64(value float64) sql.NullFloat64 {
	if value == 0 {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{
		Float64: value,
		Valid:   true,
	}
}

// getStoreIDFromContext retrieves store_id from context
// Returns 1 as default if not found for backward compatibility
func getStoreIDFromContext(ctx context.Context) uint64 {
	if ctx == nil {
		return 1
	}
	val := ctx.Value("store_id")
	if val == nil {
		return 1
	}
	switch v := val.(type) {
	case uint64:
		return v
	case int:
		return uint64(v)
	case int64:
		return uint64(v)
	default:
		return 1
	}
}
