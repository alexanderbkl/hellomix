package routes

import (
	"hellomix-backend/internal/api/handlers"
	"hellomix-backend/internal/api/middleware"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(
	transactionHandler *handlers.TransactionHandler,
	priceHandler *handlers.PriceHandler,
	addressHandler *handlers.AddressHandler,
	healthHandler *handlers.HealthHandler,
	redisClient *redis.Client,
	rateLimit int,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.Security())
	r.Use(middleware.RequestID())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Rate limiting middleware
	rateLimiter := middleware.NewRateLimiter(redisClient, rateLimit, time.Minute)
	r.Use(rateLimiter.Middleware())

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", healthHandler.Health)

		// Price endpoints
		v1.GET("/prices", priceHandler.GetPrices)

		// Exchange endpoints
		exchange := v1.Group("/exchange")
		{
			exchange.POST("/initiate", transactionHandler.InitiateExchange)
			exchange.GET("/status/:id", transactionHandler.GetTransactionStatus)
			exchange.GET("/payment/:id", transactionHandler.GetPaymentStatus)
		}

		// Address endpoints
		addresses := v1.Group("/addresses")
		{
			addresses.POST("/generate", addressHandler.GenerateBitcoinAddress)
			addresses.POST("/validate", addressHandler.ValidateAddress)
		}

		// Supported currencies
		v1.GET("/supported-currencies", healthHandler.GetSupportedCurrencies)
	}

	// Serve static files (for frontend)
	r.Static("/static", "./static")
	
	// Catch-all route for SPA
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error": "Route not found",
			"path": c.Request.URL.Path,
		})
	})

	return r
}
