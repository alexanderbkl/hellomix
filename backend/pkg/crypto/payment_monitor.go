package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/sirupsen/logrus"
)

// WalletManager handles Bitcoin wallet operations with persistent storage
type WalletManager struct {
	testnet     bool
	httpClient  *http.Client
	addresses   map[string]*btcec.PrivateKey // address -> private key mapping
	netParams   *chaincfg.Params
}

// NewWalletManager creates a new wallet manager
func NewWalletManager(testnet bool) *WalletManager {
	var netParams *chaincfg.Params
	if testnet {
		netParams = &chaincfg.TestNet3Params
	} else {
		netParams = &chaincfg.MainNetParams
	}

	return &WalletManager{
		testnet:    testnet,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		addresses:  make(map[string]*btcec.PrivateKey),
		netParams:  netParams,
	}
}

// GenerateAddressWithKey generates a new Bitcoin address and stores the private key
func (wm *WalletManager) GenerateAddressWithKey() (string, error) {
	// Generate a random private key
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create a pay-to-pubkey-hash address
	pubKey := privateKey.PubKey()
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, wm.netParams)
	if err != nil {
		return "", fmt.Errorf("failed to create address: %w", err)
	}

	addressStr := address.EncodeAddress()
	
	// Store the private key for this address
	wm.addresses[addressStr] = privateKey
	
	logrus.Infof("Generated new Bitcoin address: %s", addressStr)
	return addressStr, nil
}

// GetPrivateKey retrieves the private key for an address
func (wm *WalletManager) GetPrivateKey(address string) (*btcec.PrivateKey, error) {
	privateKey, exists := wm.addresses[address]
	if !exists {
		return nil, fmt.Errorf("private key not found for address: %s", address)
	}
	return privateKey, nil
}

// BlockchainExplorer handles blockchain API interactions
type BlockchainExplorer struct {
	testnet    bool
	httpClient *http.Client
	apiURL     string
}

// NewBlockchainExplorer creates a new blockchain explorer
func NewBlockchainExplorer(testnet bool) *BlockchainExplorer {
	apiURL := "https://blockstream.info/api"
	if testnet {
		apiURL = "https://blockstream.info/testnet/api"
	}

	return &BlockchainExplorer{
		testnet:    testnet,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiURL:     apiURL,
	}
}

// AddressInfo represents address information from blockchain API
type AddressInfo struct {
	Address               string `json:"address"`
	ChainStats            Stats  `json:"chain_stats"`
	MempoolStats          Stats  `json:"mempool_stats"`
	TotalReceived         int64  `json:"-"` // Will be calculated
	ConfirmedBalance      int64  `json:"-"` // Will be calculated
	UnconfirmedBalance    int64  `json:"-"` // Will be calculated
}

// Stats represents transaction statistics
type Stats struct {
	FundedTxoCount int64 `json:"funded_txo_count"`
	FundedTxoSum   int64 `json:"funded_txo_sum"`
	SpentTxoCount  int64 `json:"spent_txo_count"`
	SpentTxoSum    int64 `json:"spent_txo_sum"`
	TxCount        int64 `json:"tx_count"`
}

// Transaction represents a Bitcoin transaction
type Transaction struct {
	TXID     string `json:"txid"`
	Version  int    `json:"version"`
	Locktime int64  `json:"locktime"`
	Vin      []Vin  `json:"vin"`
	Vout     []Vout `json:"vout"`
	Status   Status `json:"status"`
	Fee      int64  `json:"fee"`
}

// Vin represents transaction input
type Vin struct {
	TXID    string `json:"txid"`
	Vout    int    `json:"vout"`
	Prevout Vout   `json:"prevout"`
}

// Vout represents transaction output
type Vout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int64  `json:"value"`
}

// Status represents transaction status
type Status struct {
	Confirmed   bool  `json:"confirmed"`
	BlockHeight int64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64 `json:"block_time"`
}

// GetAddressInfo gets information about a Bitcoin address
func (be *BlockchainExplorer) GetAddressInfo(ctx context.Context, address string) (*AddressInfo, error) {
	url := fmt.Sprintf("%s/address/%s", be.apiURL, address)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := be.httpClient.Do(req)
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

	var addressInfo AddressInfo
	if err := json.Unmarshal(body, &addressInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Calculate balances
	addressInfo.TotalReceived = addressInfo.ChainStats.FundedTxoSum + addressInfo.MempoolStats.FundedTxoSum
	addressInfo.ConfirmedBalance = addressInfo.ChainStats.FundedTxoSum - addressInfo.ChainStats.SpentTxoSum
	addressInfo.UnconfirmedBalance = addressInfo.MempoolStats.FundedTxoSum - addressInfo.MempoolStats.SpentTxoSum

	return &addressInfo, nil
}

// GetAddressTransactions gets transactions for a Bitcoin address
func (be *BlockchainExplorer) GetAddressTransactions(ctx context.Context, address string) ([]Transaction, error) {
	url := fmt.Sprintf("%s/address/%s/txs", be.apiURL, address)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := be.httpClient.Do(req)
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

	var transactions []Transaction
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return transactions, nil
}

// CheckPayment checks if a payment has been received to an address
func (be *BlockchainExplorer) CheckPayment(ctx context.Context, address string, expectedAmount int64) (*PaymentStatus, error) {
	addressInfo, err := be.GetAddressInfo(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get address info: %w", err)
	}

	transactions, err := be.GetAddressTransactions(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	status := &PaymentStatus{
		Address:            address,
		ExpectedAmount:     expectedAmount,
		TotalReceived:      addressInfo.TotalReceived,
		ConfirmedBalance:   addressInfo.ConfirmedBalance,
		UnconfirmedBalance: addressInfo.UnconfirmedBalance,
		Transactions:       transactions,
	}

	// Check if payment is sufficient
	if addressInfo.ConfirmedBalance >= expectedAmount {
		status.Status = "confirmed"
		status.Confirmations = 1 // At least 1 confirmation
		
		// Get exact confirmation count for the payment transaction
		for _, tx := range transactions {
			if tx.Status.Confirmed {
				// Find if this transaction has an output to our address with sufficient amount
				for _, vout := range tx.Vout {
					if vout.ScriptPubKeyAddress == address && vout.Value >= expectedAmount {
						// Calculate confirmations (simplified)
						if tx.Status.BlockHeight > 0 {
							// In a real implementation, you'd get current block height
							status.Confirmations = 1 // Simplified
						}
						status.PaymentTXID = tx.TXID
						break
					}
				}
			}
		}
	} else if addressInfo.UnconfirmedBalance >= expectedAmount {
		status.Status = "unconfirmed"
		status.Confirmations = 0
		
		// Find the unconfirmed transaction
		for _, tx := range transactions {
			if !tx.Status.Confirmed {
				for _, vout := range tx.Vout {
					if vout.ScriptPubKeyAddress == address && vout.Value >= expectedAmount {
						status.PaymentTXID = tx.TXID
						break
					}
				}
			}
		}
	} else {
		status.Status = "pending"
		status.Confirmations = 0
	}

	return status, nil
}

// PaymentStatus represents the status of a payment
type PaymentStatus struct {
	Address            string        `json:"address"`
	ExpectedAmount     int64         `json:"expected_amount"`
	TotalReceived      int64         `json:"total_received"`
	ConfirmedBalance   int64         `json:"confirmed_balance"`
	UnconfirmedBalance int64         `json:"unconfirmed_balance"`
	Status             string        `json:"status"` // pending, unconfirmed, confirmed
	Confirmations      int           `json:"confirmations"`
	PaymentTXID        string        `json:"payment_txid,omitempty"`
	Transactions       []Transaction `json:"transactions,omitempty"`
}

// PaymentMonitor monitors Bitcoin payments
type PaymentMonitor struct {
	explorer *BlockchainExplorer
	wallet   *WalletManager
}

// NewPaymentMonitor creates a new payment monitor
func NewPaymentMonitor(testnet bool) *PaymentMonitor {
	return &PaymentMonitor{
		explorer: NewBlockchainExplorer(testnet),
		wallet:   NewWalletManager(testnet),
	}
}

// MonitorPayment monitors a payment to an address
func (pm *PaymentMonitor) MonitorPayment(ctx context.Context, address string, expectedAmountSats int64) (*PaymentStatus, error) {
	return pm.explorer.CheckPayment(ctx, address, expectedAmountSats)
}

// GeneratePaymentAddress generates a new address for receiving payments
func (pm *PaymentMonitor) GeneratePaymentAddress() (string, error) {
	return pm.wallet.GenerateAddressWithKey()
}

// SatoshisToBTC converts satoshis to BTC
func SatoshisToBTC(satoshis int64) float64 {
	return float64(satoshis) / 100000000.0
}

// BTCToSatoshis converts BTC to satoshis
func BTCToSatoshis(btc float64) int64 {
	return int64(btc * 100000000.0)
}
