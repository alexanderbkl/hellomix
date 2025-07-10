package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"hellomix-backend/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// PriceService handles cryptocurrency price operations
type PriceService struct {
	db           *gorm.DB
	redis        *redis.Client
	httpClient   *http.Client
	apiKey       string
	cacheExpiry  time.Duration
}

// NewPriceService creates a new price service
func NewPriceService(db *gorm.DB, redisClient *redis.Client, apiKey string) *PriceService {
	return &PriceService{
		db:          db,
		redis:       redisClient,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		apiKey:      apiKey,
		cacheExpiry: 5 * time.Minute, // Cache prices for 5 minutes
	}
}

// CoinGeckoResponse represents the response from CoinGecko API
type CoinGeckoResponse map[string]map[string]float64

// GetPrices fetches current prices for supported cryptocurrencies
func (ps *PriceService) GetPrices(ctx context.Context) (map[string]float64, error) {
	// First, try to get prices from Redis cache
	cachedPrices, err := ps.getPricesFromCache(ctx)
	if err == nil && len(cachedPrices) > 0 {
		logrus.Debug("Returning prices from cache")
		return cachedPrices, nil
	}

	// If cache miss, fetch from API
	logrus.Info("Fetching prices from CoinGecko API")
	prices, err := ps.fetchPricesFromAPI(ctx)
	if err != nil {
		logrus.Errorf("Failed to fetch prices from API: %v", err)
		// Try to get from database as fallback
		return ps.getPricesFromDB(ctx)
	}

	// Cache the prices
	if err := ps.cachePrices(ctx, prices); err != nil {
		logrus.Warnf("Failed to cache prices: %v", err)
	}

	// Store in database
	if err := ps.storePricesInDB(ctx, prices); err != nil {
		logrus.Warnf("Failed to store prices in database: %v", err)
	}

	return prices, nil
}

// fetchPricesFromAPI fetches prices from CoinGecko API
func (ps *PriceService) fetchPricesFromAPI(ctx context.Context) (map[string]float64, error) {
	currencies := []string{"bitcoin", "ethereum", "tether", "usd-coin", "cardano", "solana", "polygon"}
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", 
		strings.Join(currencies, ","))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if ps.apiKey != "" {
		req.Header.Set("X-CG-Demo-API-Key", ps.apiKey)
	}

	resp, err := ps.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response CoinGeckoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert to our format
	prices := make(map[string]float64)
	mapping := map[string]string{
		"bitcoin":  "BTC",
		"ethereum": "ETH",
		"tether":   "USDT",
		"usd-coin": "USDC",
		"cardano":  "ADA",
		"solana":   "SOL",
		"polygon":  "MATIC",
	}

	for apiName, symbol := range mapping {
		if priceData, exists := response[apiName]; exists {
			if usdPrice, exists := priceData["usd"]; exists {
				prices[symbol] = usdPrice
			}
		}
	}

	return prices, nil
}

// getPricesFromCache retrieves prices from Redis cache
func (ps *PriceService) getPricesFromCache(ctx context.Context) (map[string]float64, error) {
	prices := make(map[string]float64)
	currencies := []string{"BTC", "ETH", "USDT", "USDC", "ADA", "SOL", "MATIC"}

	for _, currency := range currencies {
		key := fmt.Sprintf("price:%s", currency)
		priceStr, err := ps.redis.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		
		var price float64
		if err := json.Unmarshal([]byte(priceStr), &price); err != nil {
			continue
		}
		
		prices[currency] = price
	}

	return prices, nil
}

// cachePrices stores prices in Redis cache
func (ps *PriceService) cachePrices(ctx context.Context, prices map[string]float64) error {
	pipe := ps.redis.Pipeline()
	
	for currency, price := range prices {
		key := fmt.Sprintf("price:%s", currency)
		priceBytes, _ := json.Marshal(price)
		pipe.Set(ctx, key, priceBytes, ps.cacheExpiry)
	}
	
	_, err := pipe.Exec(ctx)
	return err
}

// storePricesInDB stores prices in the database
func (ps *PriceService) storePricesInDB(ctx context.Context, prices map[string]float64) error {
	for currency, price := range prices {
		priceCache := models.PriceCache{
			Currency:    currency,
			PriceUSD:    price,
			LastUpdated: time.Now(),
		}
		
		// Use UPSERT to update existing or create new
		if err := ps.db.WithContext(ctx).Save(&priceCache).Error; err != nil {
			logrus.Errorf("Failed to save price for %s: %v", currency, err)
		}
	}
	
	return nil
}

// getPricesFromDB retrieves prices from database (fallback)
func (ps *PriceService) getPricesFromDB(ctx context.Context) (map[string]float64, error) {
	var priceCaches []models.PriceCache
	if err := ps.db.WithContext(ctx).Find(&priceCaches).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch prices from database: %w", err)
	}

	prices := make(map[string]float64)
	for _, pc := range priceCaches {
		prices[pc.Currency] = pc.PriceUSD
	}

	return prices, nil
}

// GetPrice gets the price for a specific currency
func (ps *PriceService) GetPrice(ctx context.Context, currency string) (float64, error) {
	prices, err := ps.GetPrices(ctx)
	if err != nil {
		return 0, err
	}
	
	price, exists := prices[currency]
	if !exists {
		return 0, fmt.Errorf("price not found for currency: %s", currency)
	}
	
	return price, nil
}

// CalculateExchangeRate calculates the exchange rate between two currencies
func (ps *PriceService) CalculateExchangeRate(ctx context.Context, fromCurrency, toCurrency string, amount float64) (float64, error) {
	prices, err := ps.GetPrices(ctx)
	if err != nil {
		return 0, err
	}
	
	fromPrice, exists := prices[fromCurrency]
	if !exists {
		return 0, fmt.Errorf("price not found for currency: %s", fromCurrency)
	}
	
	toPrice, exists := prices[toCurrency]
	if !exists {
		return 0, fmt.Errorf("price not found for currency: %s", toCurrency)
	}
	
	// Convert amount from fromCurrency to USD, then to toCurrency
	usdValue := amount * fromPrice
	result := usdValue / toPrice
	
	return result, nil
}
