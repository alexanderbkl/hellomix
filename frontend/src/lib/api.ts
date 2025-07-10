import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: `${API_BASE_URL}/api/v1`,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface OutputAddress {
  address: string;
  percentage: number;
}

export interface CreateExchangeRequest {
  btc_amount: number;
  output_currency: string;
  output_addresses: OutputAddress[];
}

export interface Transaction {
  id: string;
  payment_address: string;
  btc_amount: number;
  output_currency: string;
  output_addresses: OutputAddress[];
  estimated_output: number;
  fee: number;
  status: string;
  created_at: string;
  updated_at?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  timestamp?: string;
}

export interface SupportedCurrency {
  symbol: string;
  name: string;
  min_amount: number;
  max_amount: number;
  fee: number;
}

export interface PriceData {
  [symbol: string]: number;
}

export interface ValidateAddressRequest {
  address: string;
  currency: string;
}

export interface ValidateAddressResponse {
  valid: boolean;
  address: string;
  currency: string;
}

export const api = {
  // Health check
  healthCheck: async (): Promise<{ status: string; service: string; timestamp: string }> => {
    const response = await apiClient.get('/health');
    return response.data;
  },

  // Get cryptocurrency prices
  getPrices: async (): Promise<PriceData> => {
    const response = await apiClient.get<ApiResponse<PriceData>>('/prices');
    return response.data.data;
  },

  // Get supported currencies
  getSupportedCurrencies: async (): Promise<SupportedCurrency[]> => {
    const response = await apiClient.get<ApiResponse<SupportedCurrency[]>>('/supported-currencies');
    return response.data.data;
  },

  // Create exchange transaction
  createExchange: async (data: CreateExchangeRequest): Promise<ApiResponse<Transaction>> => {
    const response = await apiClient.post<ApiResponse<Transaction>>('/exchange/initiate', data);
    return response.data;
  },

  // Get transaction status
  getTransactionStatus: async (transactionId: string): Promise<ApiResponse<Transaction>> => {
    const response = await apiClient.get<ApiResponse<Transaction>>(`/exchange/status/${transactionId}`);
    return response.data;
  },

  // Generate Bitcoin address
  generateBitcoinAddress: async (): Promise<{ address: string }> => {
    const response = await apiClient.post<ApiResponse<{ address: string }>>('/addresses/generate');
    return response.data.data;
  },

  // Validate address
  validateAddress: async (data: ValidateAddressRequest): Promise<ValidateAddressResponse> => {
    const response = await apiClient.post<ApiResponse<ValidateAddressResponse>>('/addresses/validate', data);
    return response.data.data;
  },
};

// Add request interceptor for debugging
apiClient.interceptors.request.use(
  (config) => {
    console.log(`[API] ${config.method?.toUpperCase()} ${config.url}`, config.data);
    return config;
  },
  (error) => {
    console.error('[API] Request error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => {
    console.log(`[API] Response:`, response.status, response.data);
    return response;
  },
  (error) => {
    console.error('[API] Response error:', error.response?.data || error.message);
    
    // Handle specific error cases
    if (error.response?.status === 429) {
      throw new Error('Too many requests. Please try again later.');
    }
    
    if (error.response?.status >= 500) {
      throw new Error('Server error. Please try again later.');
    }
    
    if (error.response?.data?.error) {
      throw new Error(error.response.data.error);
    }
    
    throw error;
  }
);

export default apiClient;
