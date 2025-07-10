package handlers

import (
	"net/http"

	"hellomix-backend/internal/services"
	"hellomix-backend/pkg/crypto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// PriceHandler handles price-related HTTP requests
type PriceHandler struct {
	priceService *services.PriceService
}

// NewPriceHandler creates a new price handler
func NewPriceHandler(priceService *services.PriceService) *PriceHandler {
	return &PriceHandler{
		priceService: priceService,
	}
}

// GetPrices handles GET /api/v1/prices
func (ph *PriceHandler) GetPrices(c *gin.Context) {
	prices, err := ph.priceService.GetPrices(c.Request.Context())
	if err != nil {
		logrus.Errorf("Failed to get prices: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cryptocurrency prices",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": prices,
		"timestamp": "now",
	})
}

// AddressHandler handles address-related HTTP requests
type AddressHandler struct {
	bitcoinService *crypto.BitcoinService
	validator      *crypto.AddressValidator
}

// NewAddressHandler creates a new address handler
func NewAddressHandler() *AddressHandler {
	return &AddressHandler{
		bitcoinService: crypto.NewBitcoinService(false), // mainnet
		validator:      crypto.NewAddressValidator(),
	}
}

// GenerateBitcoinAddress handles POST /api/v1/addresses/generate
func (ah *AddressHandler) GenerateBitcoinAddress(c *gin.Context) {
	address, err := ah.bitcoinService.GenerateAddress()
	if err != nil {
		logrus.Errorf("Failed to generate address: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate Bitcoin address",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"address": address,
		},
	})
}

// ValidateAddressRequest represents a request to validate an address
type ValidateAddressRequest struct {
	Address  string `json:"address" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

// ValidateAddress handles POST /api/v1/addresses/validate
func (ah *AddressHandler) ValidateAddress(c *gin.Context) {
	var req ValidateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	isValid := ah.validator.ValidateAddress(req.Address, req.Currency)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"valid": isValid,
			"address": req.Address,
			"currency": req.Currency,
		},
	})
}

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health handles GET /api/v1/health
func (hh *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "hellomix-backend",
		"timestamp": "now",
	})
}

// GetSupportedCurrencies handles GET /api/v1/supported-currencies
func (hh *HealthHandler) GetSupportedCurrencies(c *gin.Context) {
	currencies := []gin.H{
		{"symbol": "BTC", "name": "Bitcoin", "min_amount": 0.001, "max_amount": 10, "fee": 0.002},
		{"symbol": "ETH", "name": "Ethereum", "min_amount": 0.01, "max_amount": 100, "fee": 0.005},
		{"symbol": "USDT", "name": "Tether", "min_amount": 10, "max_amount": 50000, "fee": 0.005},
		{"symbol": "USDC", "name": "USD Coin", "min_amount": 10, "max_amount": 50000, "fee": 0.005},
		{"symbol": "ADA", "name": "Cardano", "min_amount": 100, "max_amount": 500000, "fee": 0.005},
		{"symbol": "SOL", "name": "Solana", "min_amount": 1, "max_amount": 10000, "fee": 0.005},
		{"symbol": "MATIC", "name": "Polygon", "min_amount": 100, "max_amount": 1000000, "fee": 0.005},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": currencies,
	})
}
