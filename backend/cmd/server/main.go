package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hellomix-backend/internal/api/handlers"
	"hellomix-backend/internal/api/routes"
	"hellomix-backend/internal/config"
	"hellomix-backend/internal/database"
	"hellomix-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Configure logger
	configureLogger(cfg.Server.Mode)

	logrus.Info("Starting HelloMix Backend Server...")

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logrus.Warnf("Redis connection failed: %v", err)
		logrus.Info("Continuing without Redis caching...")
		redisClient = nil
	} else {
		logrus.Info("Connected to Redis")
	}

	// Initialize services
	priceService := services.NewPriceService(db.DB, redisClient, cfg.API.CoinGeckoAPIKey)
	
	// Use testnet from configuration
	testnet := cfg.Wallet.Testnet
	transactionService := services.NewTransactionService(db.DB, priceService, testnet)

	// Initialize handlers
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	priceHandler := handlers.NewPriceHandler(priceService)
	addressHandler := handlers.NewAddressHandler()
	healthHandler := handlers.NewHealthHandler()

	// Setup routes
	router := routes.SetupRoutes(
		transactionHandler,
		priceHandler,
		addressHandler,
		healthHandler,
		redisClient,
		cfg.API.RateLimit,
	)

	// Create HTTP server
	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.Timeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	} else {
		logrus.Info("Server exited gracefully")
	}

	// Close database connection
	if err := db.Close(); err != nil {
		logrus.Errorf("Failed to close database: %v", err)
	}

	// Close Redis connection
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			logrus.Errorf("Failed to close Redis: %v", err)
		}
	}

	logrus.Info("Server shutdown complete")
}

func configureLogger(mode string) {
	// Configure logger based on environment
	if mode == "debug" || mode == "development" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		logrus.Debug("Debug logging enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	logrus.SetOutput(os.Stdout)
}
