#!/bin/bash

# Debug script for testing HelloMix backend API

echo "🔧 HelloMix Backend Debug Script"
echo "================================="

BASE_URL="http://localhost:8080/api/v1"

echo ""
echo "1. Testing Health Check..."
curl -s -w "\nStatus: %{http_code}\nTime: %{time_total}s\n" \
     -H "Content-Type: application/json" \
     "$BASE_URL/health" | jq '.' 2>/dev/null || echo "Response not JSON"

echo ""
echo "2. Testing Supported Currencies..."
curl -s -w "\nStatus: %{http_code}\nTime: %{time_total}s\n" \
     -H "Content-Type: application/json" \
     "$BASE_URL/supported-currencies" | jq '.' 2>/dev/null || echo "Response not JSON"

echo ""
echo "3. Testing Prices..."
curl -s -w "\nStatus: %{http_code}\nTime: %{time_total}s\n" \
     -H "Content-Type: application/json" \
     "$BASE_URL/prices" | jq '.' 2>/dev/null || echo "Response not JSON"

echo ""
echo "4. Testing Exchange Estimate..."
curl -s -w "\nStatus: %{http_code}\nTime: %{time_total}s\n" \
     -H "Content-Type: application/json" \
     -d '{"btc_amount": 0.001, "output_currency": "ETH"}' \
     "$BASE_URL/exchange/estimate" | jq '.' 2>/dev/null || echo "Response not JSON"

echo ""
echo "5. Testing Invalid Exchange (should fail)..."
curl -s -w "\nStatus: %{http_code}\nTime: %{time_total}s\n" \
     -H "Content-Type: application/json" \
     -d '{"btc_amount": 0, "output_currency": "INVALID"}' \
     "$BASE_URL/exchange/estimate" | jq '.' 2>/dev/null || echo "Response not JSON"

echo ""
echo "✅ Debug tests completed!"
