package services

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"hellomix-backend/internal/models"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// WalletService handles secure wallet operations
type WalletService struct {
	db         *gorm.DB
	encryptKey []byte
}

// NewWalletService creates a new wallet service
func NewWalletService(db *gorm.DB, masterKey string) *WalletService {
	// Generate encryption key from master key
	hash := sha256.Sum256([]byte(masterKey))
	
	return &WalletService{
		db:         db,
		encryptKey: hash[:],
	}
}

// encrypt encrypts data using AES-GCM
func (ws *WalletService) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(ws.encryptKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func (ws *WalletService) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(ws.encryptKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// StorePrivateKey stores an encrypted private key for an address
func (ws *WalletService) StorePrivateKey(ctx context.Context, address string, privateKey *btcec.PrivateKey, transactionID uuid.UUID) error {
	// Serialize private key
	privateKeyBytes := privateKey.Serialize()
	
	// Encrypt private key
	encryptedKey, err := ws.encrypt(privateKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Store in database
	wallet := models.Wallet{
		ID:               uuid.New(),
		Address:          address,
		EncryptedPrivKey: hex.EncodeToString(encryptedKey),
		TransactionID:    &transactionID,
		IsActive:         true,
	}

	if err := ws.db.WithContext(ctx).Create(&wallet).Error; err != nil {
		return fmt.Errorf("failed to store wallet: %w", err)
	}

	logrus.Infof("Stored encrypted private key for address: %s", address)
	return nil
}

// GetPrivateKey retrieves and decrypts a private key for an address
func (ws *WalletService) GetPrivateKey(ctx context.Context, address string) (*btcec.PrivateKey, error) {
	var wallet models.Wallet
	if err := ws.db.WithContext(ctx).Where("address = ? AND is_active = ?", address, true).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("private key not found for address: %s", address)
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Decode encrypted key
	encryptedBytes, err := hex.DecodeString(wallet.EncryptedPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted key: %w", err)
	}

	// Decrypt private key
	privateKeyBytes, err := ws.decrypt(encryptedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// Parse private key
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)
	return privateKey, nil
}

// GetWalletByTransaction gets the wallet associated with a transaction
func (ws *WalletService) GetWalletByTransaction(ctx context.Context, transactionID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := ws.db.WithContext(ctx).Where("transaction_id = ? AND is_active = ?", transactionID, true).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet not found for transaction: %s", transactionID)
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// ListActiveWallets lists all active wallets
func (ws *WalletService) ListActiveWallets(ctx context.Context) ([]models.Wallet, error) {
	var wallets []models.Wallet
	if err := ws.db.WithContext(ctx).Where("is_active = ?", true).Find(&wallets).Error; err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}

	return wallets, nil
}

// DeactivateWallet deactivates a wallet (for security)
func (ws *WalletService) DeactivateWallet(ctx context.Context, address string) error {
	result := ws.db.WithContext(ctx).Model(&models.Wallet{}).
		Where("address = ?", address).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to deactivate wallet: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("wallet not found")
	}

	logrus.Infof("Deactivated wallet: %s", address)
	return nil
}
