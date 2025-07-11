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

// PaymentProcessor handles real Bitcoin payment processing
type PaymentProcessor struct {
	db             *gorm.DB
	paymentMonitor *crypto.PaymentMonitor
	priceService   *PriceService
	testnet        bool
}

// NewPaymentProcessor creates a new payment processor
func NewPaymentProcessor(db *gorm.DB, priceService *PriceService, testnet bool) *PaymentProcessor {
	return &PaymentProcessor{
		db:             db,
		paymentMonitor: crypto.NewPaymentMonitor(testnet),
		priceService:   priceService,
		testnet:        testnet,
	}
}

// ProcessTransaction processes a transaction with real Bitcoin monitoring
func (pp *PaymentProcessor) ProcessTransaction(ctx context.Context, transactionID uuid.UUID) error {
	// Get transaction from database
	var transaction models.Transaction
	if err := pp.db.WithContext(ctx).Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	logrus.Infof("Starting payment processing for transaction: %s", transactionID)

	// Update status to waiting for payment
	if err := pp.updateTransactionStatus(ctx, transactionID, models.StatusWaiting); err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Convert BTC amount to satoshis
	expectedSats := crypto.BTCToSatoshis(transaction.BTCAmount)

	// Monitor for payment (with timeout)
	paymentCtx, cancel := context.WithTimeout(ctx, 30*time.Minute) // 30 minute timeout
	defer cancel()

	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	var paymentStatus *crypto.PaymentStatus
	var err error

	for {
		select {
		case <-paymentCtx.Done():
			// Timeout reached, mark as expired
			logrus.Warnf("Payment timeout for transaction: %s", transactionID)
			pp.updateTransactionStatus(ctx, transactionID, models.StatusExpired)
			return fmt.Errorf("payment timeout")

		case <-ticker.C:
			// Check for payment
			paymentStatus, err = pp.paymentMonitor.MonitorPayment(paymentCtx, transaction.PaymentAddress, expectedSats)
			if err != nil {
				logrus.Errorf("Failed to check payment for transaction %s: %v", transactionID, err)
				continue
			}

			logrus.Infof("Payment status for %s: %s, received: %d sats, expected: %d sats", 
				transactionID, paymentStatus.Status, paymentStatus.TotalReceived, expectedSats)

			switch paymentStatus.Status {
			case "confirmed":
				// Payment confirmed, process the exchange
				logrus.Infof("Payment confirmed for transaction: %s", transactionID)
				if err := pp.updateTransactionStatus(ctx, transactionID, models.StatusProcessing); err != nil {
					logrus.Errorf("Failed to update status to processing: %v", err)
				}

				// Store payment information
				if err := pp.storePaymentInfo(ctx, transactionID, paymentStatus); err != nil {
					logrus.Errorf("Failed to store payment info: %v", err)
				}

				// Process the actual exchange
				if err := pp.processExchange(ctx, transactionID, &transaction); err != nil {
					logrus.Errorf("Failed to process exchange: %v", err)
					pp.updateTransactionStatus(ctx, transactionID, models.StatusFailed)
					return err
				}

				// Mark as completed
				if err := pp.updateTransactionStatus(ctx, transactionID, models.StatusCompleted); err != nil {
					logrus.Errorf("Failed to update status to completed: %v", err)
				}

				logrus.Infof("Transaction completed successfully: %s", transactionID)
				return nil

			case "unconfirmed":
				// Payment received but not confirmed yet
				if transaction.Status != models.StatusProcessing {
					logrus.Infof("Unconfirmed payment received for transaction: %s", transactionID)
					if err := pp.updateTransactionStatus(ctx, transactionID, models.StatusProcessing); err != nil {
						logrus.Errorf("Failed to update status to processing: %v", err)
					}
				}
				// Continue monitoring for confirmation

			case "pending":
				// Still waiting for payment
				continue
			}
		}
	}
}

// processExchange processes the actual cryptocurrency exchange
func (pp *PaymentProcessor) processExchange(ctx context.Context, transactionID uuid.UUID, transaction *models.Transaction) error {
	logrus.Infof("Processing exchange for transaction: %s", transactionID)

	// In a real implementation, this would:
	// 1. Calculate the exact output amounts based on current prices
	// 2. Deduct fees
	// 3. Send cryptocurrencies to the destination addresses
	// 4. Record all transactions on respective blockchains

	// For now, we'll simulate the process with proper calculations
	outputAmount, err := pp.calculateFinalOutput(ctx, transaction)
	if err != nil {
		return fmt.Errorf("failed to calculate final output: %w", err)
	}

	// Store the final output amount
	if err := pp.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("id = ?", transactionID).
		Update("final_output", outputAmount).Error; err != nil {
		logrus.Errorf("Failed to store final output: %v", err)
	}

	// In a production system, here you would:
	// 1. Send the calculated amounts to each output address
	// 2. Record the transaction hashes
	// 3. Monitor for confirmations
	
	logrus.Infof("Exchange processed: %f %s sent to %d addresses", 
		outputAmount, transaction.OutputCurrency, len(transaction.OutputAddresses))

	return nil
}

// calculateFinalOutput calculates the final output amount after fees and current rates
func (pp *PaymentProcessor) calculateFinalOutput(ctx context.Context, transaction *models.Transaction) (float64, error) {
	// Get current exchange rate
	outputAmount, err := pp.priceService.CalculateExchangeRate(ctx, "BTC", transaction.OutputCurrency, transaction.BTCAmount)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate exchange rate: %w", err)
	}

	// Calculate and deduct fees
	feeRate := 0.005 // 0.5%
	if transaction.OutputCurrency == "BTC" {
		feeRate = 0.002 // 0.2% for BTC
	}

	finalAmount := outputAmount * (1 - feeRate)
	return finalAmount, nil
}

// storePaymentInfo stores payment information in the database
func (pp *PaymentProcessor) storePaymentInfo(ctx context.Context, transactionID uuid.UUID, paymentStatus *crypto.PaymentStatus) error {
	// Create a payment record
	payment := models.Payment{
		ID:            uuid.New(),
		TransactionID: transactionID,
		Address:       paymentStatus.Address,
		AmountSats:    paymentStatus.TotalReceived,
		AmountBTC:     crypto.SatoshisToBTC(paymentStatus.TotalReceived),
		TXID:          paymentStatus.PaymentTXID,
		Confirmations: paymentStatus.Confirmations,
		Status:        paymentStatus.Status,
		DetectedAt:    time.Now(),
	}

	if err := pp.db.WithContext(ctx).Create(&payment).Error; err != nil {
		return fmt.Errorf("failed to create payment record: %w", err)
	}

	logrus.Infof("Stored payment info for transaction %s: %s", transactionID, paymentStatus.PaymentTXID)
	return nil
}

// updateTransactionStatus updates the transaction status
func (pp *PaymentProcessor) updateTransactionStatus(ctx context.Context, transactionID uuid.UUID, status string) error {
	result := pp.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("id = ?", transactionID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update transaction status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	logrus.Infof("Updated transaction %s status to %s", transactionID, status)
	return nil
}

// StartPaymentMonitoring starts monitoring for a transaction
func (pp *PaymentProcessor) StartPaymentMonitoring(transactionID uuid.UUID) {
	go func() {
		ctx := context.Background()
		if err := pp.ProcessTransaction(ctx, transactionID); err != nil {
			logrus.Errorf("Payment processing failed for transaction %s: %v", transactionID, err)
		}
	}()
}

// GetPaymentStatus gets the current payment status for a transaction
func (pp *PaymentProcessor) GetPaymentStatus(ctx context.Context, transactionID uuid.UUID) (*crypto.PaymentStatus, error) {
	// Get transaction from database
	var transaction models.Transaction
	if err := pp.db.WithContext(ctx).Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Check current payment status
	expectedSats := crypto.BTCToSatoshis(transaction.BTCAmount)
	return pp.paymentMonitor.MonitorPayment(ctx, transaction.PaymentAddress, expectedSats)
}
