export const formatCurrency = (amount: number, decimals = 8): string => {
  return amount.toFixed(decimals).replace(/\.?0+$/, '');
};

export const formatNumber = (num: number, decimals = 2): string => {
  if (num >= 1e9) {
    return (num / 1e9).toFixed(decimals) + 'B';
  }
  if (num >= 1e6) {
    return (num / 1e6).toFixed(decimals) + 'M';
  }
  if (num >= 1e3) {
    return (num / 1e3).toFixed(decimals) + 'K';
  }
  return num.toFixed(decimals);
};

export const formatPrice = (price: number): string => {
  if (price < 0.01) {
    return `$${price.toFixed(6)}`;
  }
  if (price < 1) {
    return `$${price.toFixed(4)}`;
  }
  return `$${price.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
};

export const formatPercentage = (percentage: number): string => {
  return `${percentage.toFixed(2)}%`;
};

export const truncateAddress = (address: string, startChars = 6, endChars = 4): string => {
  if (address.length <= startChars + endChars) {
    return address;
  }
  return `${address.slice(0, startChars)}...${address.slice(-endChars)}`;
};

export const getCurrencyIcon = (symbol: string): string => {
  const icons: { [key: string]: string } = {
    BTC: '₿',
    ETH: 'Ξ',
    USDT: '₮',
    USDC: '$',
    ADA: '₳',
    SOL: '◎',
    MATIC: '⬟',
  };
  return icons[symbol.toUpperCase()] || symbol;
};

export const getCurrencyName = (symbol: string): string => {
  const names: { [key: string]: string } = {
    BTC: 'Bitcoin',
    ETH: 'Ethereum',
    USDT: 'Tether',
    USDC: 'USD Coin',
    ADA: 'Cardano',
    SOL: 'Solana',
    MATIC: 'Polygon',
  };
  return names[symbol.toUpperCase()] || symbol;
};

export const getStatusColor = (status: string): string => {
  const colors: { [key: string]: string } = {
    pending: 'text-yellow-500',
    processing: 'text-blue-500',
    completed: 'text-green-500',
    failed: 'text-red-500',
    expired: 'text-gray-500',
  };
  return colors[status.toLowerCase()] || 'text-gray-500';
};

export const getStatusIcon = (status: string): string => {
  const icons: { [key: string]: string } = {
    pending: '⏳',
    processing: '⚡',
    completed: '✅',
    failed: '❌',
    expired: '⏰',
  };
  return icons[status.toLowerCase()] || '⚪';
};

export const isValidBitcoinAddress = (address: string): boolean => {
  // Basic Bitcoin address validation
  const btcRegex = /^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$|^bc1[a-z0-9]{39,59}$/;
  return btcRegex.test(address);
};

export const isValidEthereumAddress = (address: string): boolean => {
  // Basic Ethereum address validation
  const ethRegex = /^0x[a-fA-F0-9]{40}$/;
  return ethRegex.test(address);
};

export const validateCryptoAddress = (address: string, currency: string): boolean => {
  switch (currency.toUpperCase()) {
    case 'BTC':
      return isValidBitcoinAddress(address);
    case 'ETH':
    case 'USDT':
    case 'USDC':
    case 'MATIC':
      return isValidEthereumAddress(address);
    case 'ADA':
      // Basic Cardano address validation
      return /^addr1[a-z0-9]{53,}$/.test(address);
    case 'SOL':
      // Basic Solana address validation
      return /^[1-9A-HJ-NP-Za-km-z]{32,44}$/.test(address);
    default:
      return false;
  }
};

export const calculateExchangeAmount = (
  btcAmount: number,
  btcPrice: number,
  outputPrice: number,
  fee: number = 0.005 // 0.5% default fee
): number => {
  const usdValue = btcAmount * btcPrice;
  const feeAmount = usdValue * fee;
  const netUsdValue = usdValue - feeAmount;
  return netUsdValue / outputPrice;
};

export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch (error) {
    console.error('Failed to copy to clipboard:', error);
    return false;
  }
};

export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  delay: number
): ((...args: Parameters<T>) => void) => {
  let timeoutId: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func.apply(null, args), delay);
  };
};

export const generateQRCodeData = (address: string, amount?: number, currency?: string): string => {
  let qrData = address;
  
  if (amount && currency) {
    if (currency.toUpperCase() === 'BTC') {
      qrData = `bitcoin:${address}?amount=${amount}`;
    } else if (currency.toUpperCase() === 'ETH') {
      qrData = `ethereum:${address}?value=${amount}`;
    }
  }
  
  return qrData;
};
