package handlers

import (
	"net/http"

	"hellomix-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	transactionService *services.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// InitiateExchange handles POST /api/v1/exchange/initiate
func (th *TransactionHandler) InitiateExchange(c *gin.Context) {
	var req services.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	transaction, err := th.transactionService.CreateTransaction(c.Request.Context(), &req)
	if err != nil {
		logrus.Errorf("Failed to create transaction: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create transaction",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"transaction_id":    transaction.ID,
			"payment_address":   transaction.PaymentAddress,
			"btc_amount":        transaction.BTCAmount,
			"output_currency":   transaction.OutputCurrency,
			"output_addresses":  transaction.OutputAddresses,
			"estimated_output":  transaction.EstimatedOutput,
			"fee":              transaction.Fee,
			"status":           transaction.Status,
			"created_at":       transaction.CreatedAt,
		},
	})
}

// GetTransactionStatus handles GET /api/v1/exchange/status/:id
func (th *TransactionHandler) GetTransactionStatus(c *gin.Context) {
	idParam := c.Param("id")
	transactionID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction ID",
		})
		return
	}

	transaction, err := th.transactionService.GetTransaction(c.Request.Context(), transactionID)
	if err != nil {
		logrus.Errorf("Failed to get transaction: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"transaction_id":    transaction.ID,
			"payment_address":   transaction.PaymentAddress,
			"btc_amount":        transaction.BTCAmount,
			"output_currency":   transaction.OutputCurrency,
			"output_addresses":  transaction.OutputAddresses,
			"estimated_output":  transaction.EstimatedOutput,
			"fee":              transaction.Fee,
			"status":           transaction.Status,
			"created_at":       transaction.CreatedAt,
			"updated_at":       transaction.UpdatedAt,
		},
	})
}
