package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction represents a cryptocurrency exchange transaction
type Transaction struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BTCAmount       float64         `json:"btc_amount" gorm:"type:decimal(18,8);not null"`
	OutputCurrency  string          `json:"output_currency" gorm:"type:varchar(10);not null"`
	OutputAddresses OutputAddresses `json:"output_addresses" gorm:"type:jsonb;not null"`
	PaymentAddress  string          `json:"payment_address" gorm:"type:varchar(100);not null"`
	Status          string          `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Fee             float64         `json:"fee" gorm:"type:decimal(18,8);default:0"`
	EstimatedOutput float64         `json:"estimated_output" gorm:"type:decimal(18,8)"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// OutputAddress represents a destination address with percentage allocation
type OutputAddress struct {
	Address    string  `json:"address"`
	Percentage float64 `json:"percentage"`
}

// OutputAddresses is a slice of OutputAddress that implements sql.Scanner and driver.Valuer
type OutputAddresses []OutputAddress

// Scan implements sql.Scanner interface
func (oa *OutputAddresses) Scan(value interface{}) error {
	if value == nil {
		*oa = OutputAddresses{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, oa)
}

// Value implements driver.Valuer interface
func (oa OutputAddresses) Value() (driver.Value, error) {
	if len(oa) == 0 {
		return nil, nil
	}
	return json.Marshal(oa)
}

// PriceCache represents cached cryptocurrency prices
type PriceCache struct {
	Currency    string    `json:"currency" gorm:"primary_key;type:varchar(10)"`
	PriceUSD    float64   `json:"price_usd" gorm:"type:decimal(18,8);not null"`
	LastUpdated time.Time `json:"last_updated" gorm:"default:now()"`
}

// SupportedCurrency represents supported cryptocurrencies
type SupportedCurrency struct {
	Symbol      string  `json:"symbol" gorm:"primary_key;type:varchar(10)"`
	Name        string  `json:"name" gorm:"type:varchar(50);not null"`
	MinAmount   float64 `json:"min_amount" gorm:"type:decimal(18,8);default:0"`
	MaxAmount   float64 `json:"max_amount" gorm:"type:decimal(18,8);default:0"`
	Fee         float64 `json:"fee" gorm:"type:decimal(5,4);default:0.005"` // 0.5% default fee
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TransactionStatus constants
const (
	StatusPending    = "pending"
	StatusWaiting    = "waiting"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusExpired    = "expired"
)

// BeforeCreate will set a UUID rather than numeric ID.
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
