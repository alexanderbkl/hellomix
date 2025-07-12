#!/usr/bin/env bash

# HelloMix Backend Local Debug Script
# This script starts the backend locally with Delve debugger

set -e

cd backend

echo "🔧 Starting HelloMix Backend in Debug Mode..."
echo "Debug server will be available on localhost:2345"
echo "Backend API will be available on localhost:8080"
echo ""
echo "Make sure your PostgreSQL and Redis are running:"
echo "  PostgreSQL: localhost:5432 (database: hellomix, user: hellomix)"
echo "  Redis: localhost:6379"
echo ""

# Set environment variables for local development
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=hellomix
export DB_USER=hellomix
export DB_PASSWORD=hellomix_password
export DB_SSLMODE=disable
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=
export REDIS_DB=0
export SERVER_MODE=debug
export SERVER_PORT=8080
export COINGECKO_API_KEY=${COINGECKO_API_KEY:-}

# Start the backend with Delve debugger
echo "Starting backend with Delve debugger..."
echo "Connect your IDE to localhost:2345 for debugging"
echo ""

dlv debug ./cmd/server --headless --listen=:2345 --api-version=2 --accept-multiclient --continue
