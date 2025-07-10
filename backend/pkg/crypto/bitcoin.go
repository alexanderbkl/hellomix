package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// BitcoinService handles Bitcoin-related operations
type BitcoinService struct {
	testnet bool
}

// NewBitcoinService creates a new Bitcoin service
func NewBitcoinService(testnet bool) *BitcoinService {
	return &BitcoinService{
		testnet: testnet,
	}
}

// GenerateAddress generates a new Bitcoin address
func (bs *BitcoinService) GenerateAddress() (string, error) {
	// Generate a random private key
	privateKey := make([]byte, 32)
	_, err := rand.Read(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Create a hash of the private key for the address
	hash := sha256.Sum256(privateKey)
	
	// Create a simple Bitcoin-like address (for demo purposes)
	prefix := "1"
	if bs.testnet {
		prefix = "m"
	}
	
	address := prefix + hex.EncodeToString(hash[:16])
	return address, nil
}

// ValidateAddress validates a Bitcoin address (basic validation)
func (bs *BitcoinService) ValidateAddress(address string) bool {
	if len(address) < 26 || len(address) > 35 {
		return false
	}
	
	// Basic Bitcoin address validation
	return strings.HasPrefix(address, "1") || strings.HasPrefix(address, "3") || 
		   strings.HasPrefix(address, "bc1") || strings.HasPrefix(address, "tb1") ||
		   strings.HasPrefix(address, "m") || strings.HasPrefix(address, "2")
}

// AddressValidator provides validation for various cryptocurrency addresses
type AddressValidator struct{}

// NewAddressValidator creates a new address validator
func NewAddressValidator() *AddressValidator {
	return &AddressValidator{}
}

// ValidateAddress validates an address for the given cryptocurrency
func (av *AddressValidator) ValidateAddress(address, currency string) bool {
	switch currency {
	case "BTC":
		return av.validateBitcoinAddress(address)
	case "ETH", "USDT", "USDC", "MATIC":
		return av.validateEthereumAddress(address)
	case "ADA":
		return av.validateCardanoAddress(address)
	case "SOL":
		return av.validateSolanaAddress(address)
	default:
		return false
	}
}

// validateBitcoinAddress validates a Bitcoin address
func (av *AddressValidator) validateBitcoinAddress(address string) bool {
	if len(address) < 26 || len(address) > 35 {
		return false
	}
	
	// Basic Bitcoin address validation
	return strings.HasPrefix(address, "1") || strings.HasPrefix(address, "3") || 
		   strings.HasPrefix(address, "bc1") || strings.HasPrefix(address, "tb1") ||
		   strings.HasPrefix(address, "m") || strings.HasPrefix(address, "2")
}

// validateEthereumAddress validates an Ethereum-based address
func (av *AddressValidator) validateEthereumAddress(address string) bool {
	// Basic Ethereum address validation
	if len(address) != 42 {
		return false
	}
	
	if address[:2] != "0x" {
		return false
	}
	
	// Check if all characters after 0x are valid hex
	for _, char := range address[2:] {
		if !((char >= '0' && char <= '9') || 
			 (char >= 'a' && char <= 'f') || 
			 (char >= 'A' && char <= 'F')) {
			return false
		}
	}
	
	return true
}

// validateCardanoAddress validates a Cardano address
func (av *AddressValidator) validateCardanoAddress(address string) bool {
	// Basic Cardano address validation
	// Cardano addresses are typically 103 characters long and start with 'addr1'
	if len(address) < 50 || len(address) > 120 {
		return false
	}
	
	// Check for Shelley era addresses
	if len(address) >= 4 && address[:4] == "addr" {
		return true
	}
	
	// Check for Byron era addresses (legacy)
	if len(address) >= 2 && (address[:2] == "Ae" || address[:2] == "Dd") {
		return true
	}
	
	return false
}

// validateSolanaAddress validates a Solana address
func (av *AddressValidator) validateSolanaAddress(address string) bool {
	// Solana addresses are base58 encoded and typically 32-44 characters
	if len(address) < 32 || len(address) > 44 {
		return false
	}
	
	// Basic base58 character check
	validChars := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, char := range address {
		valid := false
		for _, validChar := range validChars {
			if char == validChar {
				valid = true
				break
			}
		}
		if !valid {
			return false
		}
	}
	
	return true
}
