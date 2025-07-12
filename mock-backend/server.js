const express = require('express');
const cors = require('cors');
const { v4: uuidv4 } = require('uuid');

const app = express();
const PORT = 8080;

// Middleware
app.use(cors());
app.use(express.json());

// Mock data
let transactions = {};
let priceData = {
  BTC: 45000.00,
  ETH: 3200.00,
  USDT: 1.00,
  USDC: 1.00,
  ADA: 0.45,
  SOL: 110.00,
  MATIC: 0.95
};

const supportedCurrencies = [
  { symbol: 'ETH', name: 'Ethereum', min_amount: 0.01, max_amount: 10, fee: 0.005 },
  { symbol: 'USDT', name: 'Tether', min_amount: 50, max_amount: 50000, fee: 0.005 },
  { symbol: 'USDC', name: 'USD Coin', min_amount: 50, max_amount: 50000, fee: 0.005 },
  { symbol: 'ADA', name: 'Cardano', min_amount: 100, max_amount: 100000, fee: 0.005 },
  { symbol: 'SOL', name: 'Solana', min_amount: 0.5, max_amount: 500, fee: 0.005 },
  { symbol: 'MATIC', name: 'Polygon', min_amount: 50, max_amount: 50000, fee: 0.005 }
];

// Helper function to generate Bitcoin address
const generateBitcoinAddress = () => {
  const chars = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz';
  let address = '1';
  for (let i = 0; i < 33; i++) {
    address += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return address;
};

// Routes
app.get('/api/v1/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: 'HelloMix Mock Backend',
    timestamp: new Date().toISOString()
  });
});

app.get('/api/v1/prices', (req, res) => {
  // Add some random variation to prices
  const variationFactor = 0.02; // 2% variation
  const updatedPrices = {};
  
  Object.keys(priceData).forEach(currency => {
    const variation = (Math.random() - 0.5) * 2 * variationFactor;
    updatedPrices[currency] = priceData[currency] * (1 + variation);
  });
  
  res.json({
    success: true,
    data: updatedPrices,
    timestamp: new Date().toISOString()
  });
});

app.get('/api/v1/supported-currencies', (req, res) => {
  res.json({
    success: true,
    data: supportedCurrencies,
    timestamp: new Date().toISOString()
  });
});

app.post('/api/v1/exchange/initiate', (req, res) => {
  const { btc_amount, output_currency, output_addresses } = req.body;
  
  // Validate request
  if (!btc_amount || !output_currency || !output_addresses || !Array.isArray(output_addresses)) {
    return res.status(400).json({
      success: false,
      error: 'Missing required fields'
    });
  }

  // Calculate estimated output
  const btcPrice = priceData.BTC;
  const outputPrice = priceData[output_currency.toUpperCase()];
  const currency = supportedCurrencies.find(c => c.symbol === output_currency.toUpperCase());
  const fee = currency?.fee || 0.005;
  
  const usdValue = btc_amount * btcPrice;
  const feeAmount = usdValue * fee;
  const netUsdValue = usdValue - feeAmount;
  const estimatedOutput = netUsdValue / outputPrice;

  const transaction = {
    id: uuidv4(),
    payment_address: generateBitcoinAddress(),
    btc_amount: btc_amount,
    output_currency: output_currency.toUpperCase(),
    output_addresses: output_addresses,
    estimated_output: estimatedOutput,
    fee: feeAmount,
    status: 'pending',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  };

  transactions[transaction.id] = transaction;

  res.json({
    success: true,
    data: transaction,
    timestamp: new Date().toISOString()
  });
});

app.get('/api/v1/exchange/status/:id', (req, res) => {
  const transaction = transactions[req.params.id];
  
  if (!transaction) {
    return res.status(404).json({
      success: false,
      error: 'Transaction not found'
    });
  }

  res.json({
    success: true,
    data: transaction,
    timestamp: new Date().toISOString()
  });
});

app.post('/api/v1/addresses/generate', (req, res) => {
  res.json({
    success: true,
    data: { address: generateBitcoinAddress() },
    timestamp: new Date().toISOString()
  });
});

app.post('/api/v1/addresses/validate', (req, res) => {
  const { address, currency } = req.body;
  
  if (!address || !currency) {
    return res.status(400).json({
      success: false,
      error: 'Address and currency are required'
    });
  }

  // Simple validation - just check if address has minimum length
  const isValid = address.length >= 26 && address.length <= 62;
  
  res.json({
    success: true,
    data: {
      valid: isValid,
      address: address,
      currency: currency
    },
    timestamp: new Date().toISOString()
  });
});

// Start server
app.listen(PORT, () => {
  console.log(`Mock backend server running on http://localhost:${PORT}`);
  console.log('Available endpoints:');
  console.log('  GET  /api/v1/health');
  console.log('  GET  /api/v1/prices');
  console.log('  GET  /api/v1/supported-currencies');
  console.log('  POST /api/v1/exchange/initiate');
  console.log('  GET  /api/v1/exchange/status/:id');
  console.log('  POST /api/v1/addresses/generate');
  console.log('  POST /api/v1/addresses/validate');
});
