package database

import (
	"fmt"

	"hellomix-backend/internal/config"
	"hellomix-backend/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func New(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logrus.Info("Connected to PostgreSQL database")

	database := &Database{DB: db}
	
	// Run migrations
	if err := database.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Seed initial data
	if err := database.Seed(); err != nil {
		logrus.Warn("Failed to seed database: ", err)
	}

	return database, nil
}

func (d *Database) Migrate() error {
	logrus.Info("Running database migrations...")
	
	return d.DB.AutoMigrate(
		&models.Transaction{},
		&models.PriceCache{},
		&models.SupportedCurrency{},
	)
}

func (d *Database) Seed() error {
	logrus.Info("Seeding database with initial data...")

	// Check if supported currencies already exist
	var count int64
	d.DB.Model(&models.SupportedCurrency{}).Count(&count)
	if count > 0 {
		logrus.Info("Database already seeded")
		return nil
	}

	supportedCurrencies := []models.SupportedCurrency{
		{
			Symbol:    "BTC",
			Name:      "Bitcoin",
			MinAmount: 0.001,
			MaxAmount: 10,
			Fee:       0.002, // 0.2%
			IsActive:  true,
		},
		{
			Symbol:    "ETH",
			Name:      "Ethereum",
			MinAmount: 0.01,
			MaxAmount: 100,
			Fee:       0.005, // 0.5%
			IsActive:  true,
		},
		{
			Symbol:    "USDT",
			Name:      "Tether",
			MinAmount: 10,
			MaxAmount: 50000,
			Fee:       0.005, // 0.5%
			IsActive:  true,
		},
		{
			Symbol:    "USDC",
			Name:      "USD Coin",
			MinAmount: 10,
			MaxAmount: 50000,
			Fee:       0.005, // 0.5%
			IsActive:  true,
		},
		{
			Symbol:    "ADA",
			Name:      "Cardano",
			MinAmount: 100,
			MaxAmount: 500000,
			Fee:       0.005, // 0.5%
			IsActive:  true,
		},
		{
			Symbol:    "SOL",
			Name:      "Solana",
			MinAmount: 1,
			MaxAmount: 10000,
			Fee:       0.005, // 0.5%
			IsActive:  true,
		},
		{
			Symbol:    "MATIC",
			Name:      "Polygon",
			MinAmount: 100,
			MaxAmount: 1000000,
			Fee:       0.005, // 0.5%
			IsActive:  true,
		},
	}

	for _, currency := range supportedCurrencies {
		if err := d.DB.Create(&currency).Error; err != nil {
			logrus.Errorf("Failed to seed currency %s: %v", currency.Symbol, err)
		}
	}

	logrus.Info("Database seeded successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
