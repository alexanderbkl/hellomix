package crypto

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// BitcoinService handles Bitcoin-related operations
type BitcoinService struct {
	testnet       bool
	walletManager *WalletManager
}

// NewBitcoinService creates a new Bitcoin service
func NewBitcoinService(testnet bool) *BitcoinService {
	return &BitcoinService{
		testnet:       testnet,
		walletManager: NewWalletManager(testnet),
	}
}

// GenerateAddress generates a new Bitcoin address with persistent key storage
func (bs *BitcoinService) GenerateAddress() (string, error) {
	return bs.walletManager.GenerateAddressWithKey()
}

// ValidateAddress validates a Bitcoin address using proper Bitcoin validation
func (bs *BitcoinService) ValidateAddress(address string) bool {
	// Choose the appropriate network parameters
	var netParams *chaincfg.Params
	if bs.testnet {
		netParams = &chaincfg.TestNet3Params
	} else {
		netParams = &chaincfg.MainNetParams
	}

	// Use btcutil to validate the address
	_, err := btcutil.DecodeAddress(address, netParams)
	return err == nil
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

// validateBitcoinAddress validates a Bitcoin address using proper Bitcoin validation
func (av *AddressValidator) validateBitcoinAddress(address string) bool {
	// Try to decode as mainnet first
	_, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err == nil {
		return true
	}

	// Try to decode as testnet
	_, err = btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	return err == nil
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
