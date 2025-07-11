package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	API      APIConfig
	Wallet   WalletConfig
}

type ServerConfig struct {
	Port    string
	Mode    string
	Host    string
	Timeout int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type APIConfig struct {
	CoinGeckoAPIKey string
	RateLimit       int
}

type WalletConfig struct {
	MasterKey string
	Testnet   bool
}

func Load() (*Config, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			Mode:    getEnv("GIN_MODE", "debug"),
			Host:    getEnv("HOST", "localhost"),
			Timeout: getEnvAsInt("SERVER_TIMEOUT", 30),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "hellomix"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "hellomix"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		API: APIConfig{
			CoinGeckoAPIKey: getEnv("COINGECKO_API_KEY", ""),
			RateLimit:       getEnvAsInt("RATE_LIMIT", 100),
		},
		Wallet: WalletConfig{
			MasterKey: getEnv("WALLET_MASTER_KEY", ""),
			Testnet:   getEnvAsBool("WALLET_TESTNET", false),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
