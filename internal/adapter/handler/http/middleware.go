package http

import (
	"bytes"
	"context"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/bagashiz/go-pos/internal/core/domain"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-gonic/gin"
)

var storeServiceForSubdomain port.StoreService
var subscriptionServiceForMiddleware port.SubscriptionService

func SetStoreServiceForSubdomain(svc port.StoreService) {
	storeServiceForSubdomain = svc
}

func SetSubscriptionServiceForMiddleware(svc port.SubscriptionService) {
	subscriptionServiceForMiddleware = svc
}

const (
	authorizationHeaderKey     = "authorization"
	authorizationType          = "bearer"
	authorizationPayloadKey    = "authorization_payload"
	auditLogServiceKey         = "audit_log_service"
	storeIDKey                 = "store_id"
)

type auditLogMiddlewareService interface {
	CreateAuditLog(ctx *gin.Context, log *domain.AuditLog) (*domain.AuditLog, error)
}

func auditMiddleware(svc port.AuditLogService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		if method != "POST" && method != "PUT" && method != "DELETE" {
			ctx.Next()
			return
		}

		var bodyBytes []byte
		if ctx.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		ctx.Next()

		if ctx.Writer.Status() < 400 {
			payload := getAuthPayload(ctx, authorizationPayloadKey)
			if payload == nil {
				return
			}

			action := domain.AuditAction(method)
			if method == "POST" {
				action = domain.AuditActionCreate
			} else if method == "PUT" {
				action = domain.AuditActionUpdate
			} else if method == "DELETE" {
				action = domain.AuditActionDelete
			}

			resourceType := getResourceTypeFromPath(ctx.Request.URL.Path)
			resourceID := getResourceIDFromPath(ctx.Request.URL.Path)

			clientIP := getClientIP(ctx)

			log := &domain.AuditLog{
				UserID:       payload.UserID,
				UserName:     "",
				Action:       action,
				ResourceType: resourceType,
				ResourceID:   resourceID,
				NewData:      bodyBytes,
				IPAddress:    clientIP,
			}

			logCopy := *log
			ctxCopy := context.WithValue(context.Background(), "store_id", payload.StoreID)
			go svc.CreateAuditLog(ctxCopy, &logCopy)
		}
	}
}

func getResourceTypeFromPath(path string) string {
	switch {
	case strings.Contains(path, "/users/"):
		return "users"
	case strings.Contains(path, "/payments/"):
		return "payments"
	case strings.Contains(path, "/categories/"):
		return "categories"
	case strings.Contains(path, "/products/"):
		return "products"
	case strings.Contains(path, "/orders/"):
		return "orders"
	case strings.Contains(path, "/consignment/consignors"):
		return "consignors"
	case strings.Contains(path, "/consignment/items"):
		return "consignments"
	case strings.Contains(path, "/consignment/settlements"):
		return "settlements"
	case strings.Contains(path, "/suppliers/"):
		return "suppliers"
	case strings.Contains(path, "/purchases/"):
		return "purchases"
	default:
		return "unknown"
	}
}

func getResourceIDFromPath(path string) uint64 {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if isNumericStr(parts[i]) {
			id, err := strconv.ParseUint(parts[i], 10, 64)
			if err == nil {
				return id
			}
		}
	}
	return 0
}

func isNumericStr(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func getClientIP(ctx *gin.Context) string {
	ip := ctx.GetHeader("X-Real-IP")
	if ip == "" {
		ip = ctx.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = ctx.ClientIP()
	}
	if strings.Contains(ip, ",") {
		ip = strings.Split(ip, ",")[0]
	}
	host, _, _ := net.SplitHostPort(ip)
	if host != "" {
		ip = host
	}
	return ip
}

// authMiddleware is a middleware to check if the user is authenticated
func authMiddleware(token port.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			err := domain.ErrEmptyAuthorizationHeader
			handleAbort(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			err := domain.ErrInvalidAuthorizationHeader
			handleAbort(ctx, err)
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != authorizationType {
			err := domain.ErrInvalidAuthorizationType
			handleAbort(ctx, err)
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			handleAbort(ctx, err)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Set(storeIDKey, payload.StoreID)
		ctx.Next()
	}
}

// adminMiddleware is a middleware to check if the user is an admin
func adminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := getAuthPayload(ctx, authorizationPayloadKey)

		isAdmin := payload.Role == domain.Admin
		if !isAdmin {
			err := domain.ErrForbidden
			handleAbort(ctx, err)
			return
		}

		ctx.Next()
	}
}

func subdomainMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		host := ctx.Request.Host
		host = strings.Split(host, ":")[0]

		if strings.HasSuffix(host, ".store.yuppy576.top") {
			if storeServiceForSubdomain != nil {
				store, err := storeServiceForSubdomain.GetStoreByDomain(ctx, host)
				if err == nil && store != nil {
					ctx.Set(storeIDKey, store.ID)
				}
			}
		}

		ctx.Next()
	}
}

func subscriptionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		storeIDVal, _ := ctx.Get(storeIDKey)
		if storeIDVal == nil {
			ctx.Next()
			return
		}

		var storeID uint64
		switch v := storeIDVal.(type) {
		case uint64:
			storeID = v
		case int64:
			storeID = uint64(v)
		case int:
			storeID = uint64(v)
		default:
			ctx.Next()
			return
		}

		if subscriptionServiceForMiddleware != nil {
			valid, _, err := subscriptionServiceForMiddleware.CheckSubscriptionStatus(ctx, storeID)
			if err == nil && !valid {
				handleError(ctx, domain.ErrSubscriptionExpired)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
