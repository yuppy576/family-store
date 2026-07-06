package http

import (
	"strconv"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/gin-gonic/gin"
)

// stringToUint64 is a helper function to convert a string to uint64
func stringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

// getAuthPayload is a helper function to get the auth payload from the context
func getAuthPayload(ctx *gin.Context, key string) *domain.TokenPayload {
	val, exists := ctx.Get(key)
	if !exists {
		return nil
	}
	payload, ok := val.(*domain.TokenPayload)
	if !ok {
		return nil
	}
	return payload
}

// toMap is a helper function to add meta and data to a map
func toMap(m meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}
