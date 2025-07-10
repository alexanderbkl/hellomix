🚀 HelloMix Cryptocurrency Exchange - Full-Stack Implementation
I need you to create a complete, production-ready HelloMix cryptocurrency exchange application from scratch. This will be a professional-grade platform with real cryptocurrency integration and modern security practices.

Project Overview
Create a full-stack anonymous cryptocurrency exchange platform that allows users to send Bitcoin and receive various cryptocurrencies across multiple wallet addresses. The application should be deployment-ready with comprehensive documentation.

Application Architecture
Frontend (React/Next.js Enhancement)
Convert the provided HTML/CSS/JavaScript to a modern React/Next.js application
Implement real-time cryptocurrency price feeds
Add responsive design optimizations
Include proper error handling and loading states
Integrate with backend API for live data
Add form validation and user feedback
Implement progressive web app (PWA) features
Backend (Golang/Gin)
Framework: Gin web framework for high-performance API
Database: PostgreSQL with GORM for transaction history and status tracking
Cache: Redis for price caching and rate limiting
Security: JWT tokens, rate limiting, CORS, input validation
Logging: Structured logging with logrus
Configuration: Environment-based configuration management
Core Features & Requirements
1. Real Cryptocurrency Integration
Price Feeds: CoinGecko API integration for live prices (BTC, ETH, USDT, USDC, ADA, SOL, MATIC)
Bitcoin Address Generation: Create valid Bitcoin addresses using btcutil
Address Validation: Validate wallet addresses for all supported cryptocurrencies
Exchange Rate Calculations: Real-time rate calculations with transparent fee structure
2. Transaction Processing System
Multi-step Processing: Realistic transaction flow with status updates
Transaction Tracking: Unique transaction IDs with status monitoring
Multiple Output Addresses: Support for 1-7 destination addresses with percentage allocation
Processing Simulation: Realistic timing and status updates for user experience
3. Security & Compliance
Rate Limiting: API rate limiting with Redis
Input Validation: Comprehensive input sanitization
CORS Protection: Proper CORS configuration
Error Handling: Graceful error handling with user-friendly messages
Logging: Comprehensive audit logging for all transactions
4. API Endpoints
GET    /api/v1/health              - Health check
GET    /api/v1/prices              - Live cryptocurrency prices
POST   /api/v1/exchange/initiate   - Initialize exchange transaction
GET    /api/v1/exchange/status/:id - Get transaction status
POST   /api/v1/addresses/generate  - Generate Bitcoin payment address
POST   /api/v1/addresses/validate  - Validate wallet addresses
GET    /api/v1/supported-currencies - Get supported currencies
5. Database Schema
sql
-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    btc_amount DECIMAL(18,8) NOT NULL,
    output_currency VARCHAR(10) NOT NULL,
    output_addresses JSONB NOT NULL,
    payment_address VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Price cache table
CREATE TABLE price_cache (
    currency VARCHAR(10) PRIMARY KEY,
    price_usd DECIMAL(18,8) NOT NULL,
    last_updated TIMESTAMP DEFAULT NOW()
);
Technical Stack
Backend Dependencies
go
// Core framework and middleware
github.com/gin-gonic/gin
github.com/gin-contrib/cors
github.com/gin-contrib/ratelimit

// Database and ORM
gorm.io/gorm
gorm.io/driver/postgres
github.com/go-redis/redis/v8

// Cryptocurrency utilities
github.com/btcsuite/btcutil
github.com/btcsuite/btcd/chaincfg

// Utilities
github.com/sirupsen/logrus
github.com/joho/godotenv
github.com/google/uuid
Frontend Dependencies
json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0",
    "tailwindcss": "^3.0.0",
    "axios": "^1.0.0",
    "react-query": "^3.0.0",
    "react-hook-form": "^7.0.0",
    "qrcode": "^1.5.0",
    "lucide-react": "^0.263.1"
  }
}
Project Structure
hellomix/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes/
│   │   ├── config/
│   │   ├── database/
│   │   ├── models/
│   │   ├── services/
│   │   └── utils/
│   ├── pkg/
│   │   ├── crypto/
│   │   └── validation/
│   ├── docker/
│   │   └── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── hooks/
│   │   ├── services/
│   │   └── utils/
│   ├── public/
│   ├── package.json
│   └── next.config.js
├── docker-compose.yml
├── README.md
└── .env.example
Deployment Configuration
Docker Support
Multi-stage Dockerfile for Go backend
Docker Compose for local development
Production-ready container configuration
Environment Configuration
env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=hellomix
DB_PASSWORD=your_password
DB_NAME=hellomix

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# API Keys
COINGECKO_API_KEY=your_api_key

# Server
PORT=8080
GIN_MODE=release
Testing Requirements
Unit tests for all service functions
Integration tests for API endpoints
End-to-end tests for critical user flows
Performance testing for high-load scenarios
Documentation Requirements
Comprehensive README with setup instructions
API documentation with OpenAPI/Swagger
Deployment guides for different platforms
Security best practices documentation
Success Criteria
Functionality: All features work as specified with real cryptocurrency data
Performance: API responses under 200ms for 95% of requests
Security: Passes basic security audits and best practices
Deployment: Successfully deployable on cloud platforms
Documentation: Complete documentation for setup and usage
Provided Assets
The original HTML, CSS, and JavaScript files are provided as reference for the UI/UX design. These should be modernized and converted to React/Next.js while maintaining the visual design and user experience.

Please create a complete, production-ready application that demonstrates modern web development practices with actual cryptocurrency integration. The application should be immediately deployable and functional for real-world use.

This refactored prompt provides a comprehensive specification that maintains the original vision while adding professional structure, technical depth, and clear deliverables for a production-ready application.