package http

import (
	"log/slog"
	"strings"

	"github.com/bagashiz/go-pos/internal/adapter/config"
	"github.com/bagashiz/go-pos/internal/core/port"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new HTTP router
func NewRouter(
	config *config.HTTP,
	token port.TokenService,
	userHandler UserHandler,
	authHandler AuthHandler,
	paymentHandler PaymentHandler,
	categoryHandler CategoryHandler,
	productHandler ProductHandler,
	orderHandler OrderHandler,
	consignmentHandler ConsignmentHandler,
	supplierHandler SupplierHandler,
	purchaseHandler PurchaseHandler,
	auditLogHandler AuditLogHandler,
	auditLogService port.AuditLogService,
	storeHandler StoreHandler,
	subscriptionHandler SubscriptionHandler,
	reportHandler ReportHandler,
) (*Router, error) {
	// Disable debug mode in production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig), subdomainMiddleware())

	// Custom validators
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
			return nil, err
		}

		if err := v.RegisterValidation("payment_type", paymentTypeValidator); err != nil {
			return nil, err
		}

	}

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.POST("/login", authHandler.Login)

			authUser := user.Group("/").Use(authMiddleware(token))
			{
				authUser.GET("/", userHandler.ListUsers)
				authUser.GET("/me", userHandler.GetMe)
				authUser.GET("/:id", userHandler.GetUser)

				admin := authUser.Use(adminMiddleware(), auditMiddleware(auditLogService))
				{
					admin.PUT("/:id", userHandler.UpdateUser)
					admin.DELETE("/:id", userHandler.DeleteUser)
				}
			}
		}
		payment := v1.Group("/payments").Use(authMiddleware(token))
		{
			payment.GET("/", paymentHandler.ListPayments)
			payment.GET("/:id", paymentHandler.GetPayment)

			admin := payment.Use(adminMiddleware(), auditMiddleware(auditLogService))
			{
				admin.POST("/", paymentHandler.CreatePayment)
				admin.PUT("/:id", paymentHandler.UpdatePayment)
				admin.DELETE("/:id", paymentHandler.DeletePayment)
			}
		}
		category := v1.Group("/categories").Use(authMiddleware(token))
		{
			category.GET("/", categoryHandler.ListCategories)
			category.GET("/:id", categoryHandler.GetCategory)

			admin := category.Use(adminMiddleware(), auditMiddleware(auditLogService))
			{
				admin.POST("/", categoryHandler.CreateCategory)
				admin.PUT("/:id", categoryHandler.UpdateCategory)
				admin.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}
		product := v1.Group("/products").Use(authMiddleware(token))
		{
			product.GET("/", productHandler.ListProducts)
			product.GET("/low-stock", productHandler.ListLowStockProducts)
			product.GET("/:id", productHandler.GetProduct)

			admin := product.Use(adminMiddleware(), auditMiddleware(auditLogService))
			{
				admin.POST("/", productHandler.CreateProduct)
				admin.PUT("/:id", productHandler.UpdateProduct)
				admin.DELETE("/:id", productHandler.DeleteProduct)
			}
		}
		og := v1.Group("/orders").Use(authMiddleware(token), auditMiddleware(auditLogService))
		{
			og.POST("/", orderHandler.CreateOrder)
			og.POST("", orderHandler.CreateOrder)
			og.GET("/", orderHandler.ListOrders)
			og.GET("", orderHandler.ListOrders)
			og.GET("/:id", orderHandler.GetOrder)
		}
	}

	// ── Consignment Routes ──────────────────────────────────────
	// Read routes (authenticated users)
	consignment := v1.Group("/consignment").Use(authMiddleware(token))
	{
		consignment.GET("/consignors", consignmentHandler.ListConsignors)
		consignment.GET("/consignors/:id", consignmentHandler.GetConsignor)
		consignment.GET("/items", consignmentHandler.ListConsignments)
		consignment.GET("/items/:id", consignmentHandler.GetConsignment)
		consignment.POST("/items/:id/vehicle", consignmentHandler.CreateVehicle)
		consignment.GET("/items/:id/vehicle", consignmentHandler.GetVehicle)
		consignment.PUT("/items/:id/vehicle", consignmentHandler.UpdateVehicle)
		consignment.GET("/items/:id/settlements", consignmentHandler.ListSettlements)
		consignment.GET("/vehicles/:id/progress", consignmentHandler.ListTransferProgress)
		consignment.GET("/expiring", consignmentHandler.ListExpiring)
	}
	// Admin routes
	ca := v1.Group("/consignment").Use(authMiddleware(token), adminMiddleware(), auditMiddleware(auditLogService))
	{
		ca.POST("/consignors", consignmentHandler.CreateConsignor)
		ca.PUT("/consignors/:id", consignmentHandler.UpdateConsignor)
		ca.DELETE("/consignors/:id", consignmentHandler.DeleteConsignor)
		ca.POST("/items", consignmentHandler.CreateConsignment)
		ca.PUT("/items/:id", consignmentHandler.UpdateConsignment)
		ca.DELETE("/items/:id", consignmentHandler.DeleteConsignment)
		ca.POST("/items/:id/settlements", consignmentHandler.CreateSettlement)
		ca.POST("/vehicles/:id/progress", consignmentHandler.CreateTransferProgress)
	}

	// ── Store Routes ──────────────────────────────────────
	sg := v1.Group("/suppliers").Use(authMiddleware(token))
	{
		sg.GET("/", supplierHandler.List)
		sg.GET("", supplierHandler.List)  // no trailing slash
		sg.GET("/:id", supplierHandler.Get)
		sa := sg.Use(adminMiddleware(), auditMiddleware(auditLogService))
		{
			sa.POST("/", supplierHandler.Create)
			sa.POST("", supplierHandler.Create)
			sa.PUT("/:id", supplierHandler.Update)
			sa.DELETE("/:id", supplierHandler.Delete)
		}
	}
	pg := v1.Group("/purchases").Use(authMiddleware(token))
	{
		pg.GET("/", purchaseHandler.List)
		pg.GET("", purchaseHandler.List)
		pg.GET("/:id", purchaseHandler.Get)
		pg.GET("/:id/items", purchaseHandler.ListItems)
		pa := pg.Use(adminMiddleware(), auditMiddleware(auditLogService))
		{
			pa.POST("/", purchaseHandler.Create)
			pa.POST("", purchaseHandler.Create)
		}
	}

	audit := v1.Group("/audit-logs").Use(authMiddleware(token), adminMiddleware())
	{
		audit.GET("/", auditLogHandler.ListAuditLogs)
	}

	stores := v1.Group("/stores")
	{
		stores.POST("/register", storeHandler.Register)
		stores.GET("/", storeHandler.GetStoreByDomain)
	}

	subscription := v1.Group("/subscription").Use(authMiddleware(token), adminMiddleware())
	{
		subscription.GET("/", subscriptionHandler.GetSubscription)
		subscription.POST("/renew", subscriptionHandler.RenewSubscription)
		subscription.POST("/activate", subscriptionHandler.ActivateSubscription)
	}

	report := v1.Group("/reports").Use(authMiddleware(token))
	{
		report.GET("/sales/stats", reportHandler.GetSalesStats)
		report.GET("/sales/daily", reportHandler.GetDailySales)
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
