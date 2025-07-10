package services

import (
	"context"
	"fmt"
	"time"

	"hellomix-backend/internal/models"
	"hellomix-backend/pkg/crypto"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TransactionService handles cryptocurrency exchange transactions
type TransactionService struct {
	db             *gorm.DB
	priceService   *PriceService
	bitcoinService *crypto.BitcoinService
	validator      *crypto.AddressValidator
}

// NewTransactionService creates a new transaction service
func NewTransactionService(db *gorm.DB, priceService *PriceService) *TransactionService {
	return &TransactionService{
		db:             db,
		priceService:   priceService,
		bitcoinService: crypto.NewBitcoinService(false), // Use mainnet
		validator:      crypto.NewAddressValidator(),
	}
}

// CreateTransactionRequest represents a request to create a new transaction
type CreateTransactionRequest struct {
	BTCAmount       float64                 `json:"btc_amount" binding:"required,gt=0"`
	OutputCurrency  string                  `json:"output_currency" binding:"required"`
	OutputAddresses []models.OutputAddress  `json:"output_addresses" binding:"required,min=1,max=7"`
}

// CreateTransaction creates a new exchange transaction
func (ts *TransactionService) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*models.Transaction, error) {
	// Validate output currency
	if !ts.isSupportedCurrency(req.OutputCurrency) {
		return nil, fmt.Errorf("unsupported output currency: %s", req.OutputCurrency)
	}

	// Validate output addresses
	if err := ts.validateOutputAddresses(req.OutputAddresses, req.OutputCurrency); err != nil {
		return nil, fmt.Errorf("invalid output addresses: %w", err)
	}

	// Validate percentage allocation
	if err := ts.validatePercentageAllocation(req.OutputAddresses); err != nil {
		return nil, fmt.Errorf("invalid percentage allocation: %w", err)
	}

	// Generate payment address
	paymentAddress, err := ts.bitcoinService.GenerateAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment address: %w", err)
	}

	// Calculate estimated output
	estimatedOutput, err := ts.calculateEstimatedOutput(ctx, req.BTCAmount, req.OutputCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate estimated output: %w", err)
	}

	// Calculate fee
	fee := ts.calculateFee(req.BTCAmount, req.OutputCurrency)

	// Create transaction
	transaction := &models.Transaction{
		ID:              uuid.New(),
		BTCAmount:       req.BTCAmount,
		OutputCurrency:  req.OutputCurrency,
		OutputAddresses: models.OutputAddresses(req.OutputAddresses),
		PaymentAddress:  paymentAddress,
		Status:          models.StatusPending,
		Fee:             fee,
		EstimatedOutput: estimatedOutput,
	}

	if err := ts.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	logrus.Infof("Created new transaction: %s", transaction.ID)
	
	// Start background processing
	go ts.processTransactionAsync(transaction.ID)

	return transaction, nil
}

// GetTransaction retrieves a transaction by ID
func (ts *TransactionService) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := ts.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// UpdateTransactionStatus updates the status of a transaction
func (ts *TransactionService) UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status string) error {
	result := ts.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update transaction status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	logrus.Infof("Updated transaction %s status to %s", id, status)
	return nil
}

// isSupportedCurrency checks if the currency is supported
func (ts *TransactionService) isSupportedCurrency(currency string) bool {
	supportedCurrencies := []string{"BTC", "ETH", "USDT", "USDC", "ADA", "SOL", "MATIC"}
	for _, supported := range supportedCurrencies {
		if currency == supported {
			return true
		}
	}
	return false
}

// validateOutputAddresses validates the output addresses
func (ts *TransactionService) validateOutputAddresses(addresses []models.OutputAddress, currency string) error {
	if len(addresses) == 0 {
		return fmt.Errorf("at least one output address is required")
	}

	if len(addresses) > 7 {
		return fmt.Errorf("maximum 7 output addresses allowed")
	}

	for i, addr := range addresses {
		if addr.Address == "" {
			return fmt.Errorf("address %d is empty", i+1)
		}

		if !ts.validator.ValidateAddress(addr.Address, currency) {
			return fmt.Errorf("invalid address %d for currency %s", i+1, currency)
		}

		if addr.Percentage <= 0 || addr.Percentage > 100 {
			return fmt.Errorf("invalid percentage for address %d: %.2f", i+1, addr.Percentage)
		}
	}

	return nil
}

// validatePercentageAllocation validates that percentages add up to 100%
func (ts *TransactionService) validatePercentageAllocation(addresses []models.OutputAddress) error {
	total := 0.0
	for _, addr := range addresses {
		total += addr.Percentage
	}

	if total < 99.9 || total > 100.1 { // Allow for small floating point errors
		return fmt.Errorf("percentage allocation must equal 100%%, got %.2f%%", total)
	}

	return nil
}

// calculateEstimatedOutput calculates the estimated output amount
func (ts *TransactionService) calculateEstimatedOutput(ctx context.Context, btcAmount float64, outputCurrency string) (float64, error) {
	if outputCurrency == "BTC" {
		return btcAmount, nil
	}

	// Get exchange rate
	outputAmount, err := ts.priceService.CalculateExchangeRate(ctx, "BTC", outputCurrency, btcAmount)
	if err != nil {
		return 0, err
	}

	// Subtract fee
	fee := ts.calculateFee(btcAmount, outputCurrency)
	
	// Convert fee to output currency
	feeInOutputCurrency, err := ts.priceService.CalculateExchangeRate(ctx, "BTC", outputCurrency, fee)
	if err != nil {
		return 0, err
	}

	return outputAmount - feeInOutputCurrency, nil
}

// calculateFee calculates the transaction fee
func (ts *TransactionService) calculateFee(btcAmount float64, currency string) float64 {
	// Fee structure: 0.5% for most currencies, 0.2% for BTC
	feeRate := 0.005 // 0.5%
	if currency == "BTC" {
		feeRate = 0.002 // 0.2%
	}

	return btcAmount * feeRate
}

// processTransactionAsync processes the transaction in the background
func (ts *TransactionService) processTransactionAsync(transactionID uuid.UUID) {
	ctx := context.Background()
	
	// Simulate transaction processing
	stages := []struct {
		status   string
		duration time.Duration
	}{
		{models.StatusWaiting, 30 * time.Second},
		{models.StatusProcessing, 2 * time.Minute},
		{models.StatusCompleted, 0},
	}

	for _, stage := range stages {
		time.Sleep(stage.duration)
		
		if err := ts.UpdateTransactionStatus(ctx, transactionID, stage.status); err != nil {
			logrus.Errorf("Failed to update transaction status: %v", err)
			// Mark as failed
			ts.UpdateTransactionStatus(ctx, transactionID, models.StatusFailed)
			return
		}
		
		logrus.Infof("Transaction %s moved to status: %s", transactionID, stage.status)
	}
}

// GetTransactionHistory gets transaction history (for admin purposes)
func (ts *TransactionService) GetTransactionHistory(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	
	query := ts.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset)
	
	if err := query.Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	return transactions, nil
}
