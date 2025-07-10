import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api, CreateExchangeRequest, Transaction, SupportedCurrency, PriceData, ValidateAddressRequest } from '@/lib/api';

// Query Keys
export const queryKeys = {
  prices: ['prices'] as const,
  supportedCurrencies: ['supportedCurrencies'] as const,
  transaction: (id: string) => ['transaction', id] as const,
  health: ['health'] as const,
};

// Prices Hook
export const usePrices = () => {
  return useQuery({
    queryKey: queryKeys.prices,
    queryFn: api.getPrices,
    refetchInterval: 30000, // Refetch every 30 seconds
    staleTime: 25000, // Consider data stale after 25 seconds
    retry: 3,
  });
};

// Supported Currencies Hook
export const useSupportedCurrencies = () => {
  return useQuery({
    queryKey: queryKeys.supportedCurrencies,
    queryFn: api.getSupportedCurrencies,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 3,
  });
};

// Transaction Status Hook
export const useTransactionStatus = (transactionId: string, enabled = true) => {
  return useQuery({
    queryKey: queryKeys.transaction(transactionId),
    queryFn: () => api.getTransactionStatus(transactionId),
    enabled: !!transactionId && enabled,
    refetchInterval: (query) => {
      // Stop refetching if transaction is completed, failed, or expired
      const status = query.state.data?.data?.status;
      if (status === 'completed' || status === 'failed' || status === 'expired') {
        return false;
      }
      return 5000; // Refetch every 5 seconds for pending/processing
    },
    retry: 3,
  });
};

// Health Check Hook
export const useHealthCheck = () => {
  return useQuery({
    queryKey: queryKeys.health,
    queryFn: api.healthCheck,
    staleTime: 60000, // 1 minute
    retry: 1,
  });
};

// Create Exchange Mutation
export const useCreateExchange = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateExchangeRequest) => api.createExchange(data),
    onSuccess: (data) => {
      // Invalidate and refetch transaction data
      queryClient.invalidateQueries({ queryKey: queryKeys.transaction(data.data.id) });
    },
    onError: (error) => {
      console.error('Exchange creation failed:', error);
    },
  });
};

// Validate Address Mutation
export const useValidateAddress = () => {
  return useMutation({
    mutationFn: (data: ValidateAddressRequest) => api.validateAddress(data),
  });
};

// Generate Bitcoin Address Mutation
export const useGenerateBitcoinAddress = () => {
  return useMutation({
    mutationFn: api.generateBitcoinAddress,
  });
};

// Custom hook for real-time price calculations
export const usePriceCalculation = (btcAmount: number, outputCurrency: string) => {
  const { data: prices, isLoading: pricesLoading } = usePrices();
  const { data: currencies, isLoading: currenciesLoading } = useSupportedCurrencies();

  const isLoading = pricesLoading || currenciesLoading;

  const calculation = React.useMemo(() => {
    if (!prices || !currencies || !btcAmount || !outputCurrency) {
      return null;
    }

    const btcPrice = prices.BTC || 0;
    const outputPrice = prices[outputCurrency.toUpperCase()] || 0;
    const currency = currencies.find(c => c.symbol.toUpperCase() === outputCurrency.toUpperCase());
    const fee = currency?.fee || 0.005;

    if (!btcPrice || !outputPrice) {
      return null;
    }

    const usdValue = btcAmount * btcPrice;
    const feeAmount = usdValue * fee;
    const netUsdValue = usdValue - feeAmount;
    const outputAmount = netUsdValue / outputPrice;

    return {
      btcPrice,
      outputPrice,
      usdValue,
      feeAmount,
      feePercentage: fee * 100,
      netUsdValue,
      outputAmount,
      currency,
    };
  }, [prices, currencies, btcAmount, outputCurrency]);

  return {
    calculation,
    isLoading,
    prices,
    currencies,
  };
};

// Custom hook for exchange form state management
export const useExchangeForm = () => {
  const [currentStep, setCurrentStep] = React.useState(1);
  const [transactionId, setTransactionId] = React.useState<string | null>(null);
  
  const createExchange = useCreateExchange();
  const validateAddress = useValidateAddress();

  const goToNextStep = () => setCurrentStep(prev => prev + 1);
  const goToPrevStep = () => setCurrentStep(prev => prev - 1);
  const resetForm = () => {
    setCurrentStep(1);
    setTransactionId(null);
  };

  return {
    currentStep,
    transactionId,
    setTransactionId,
    goToNextStep,
    goToPrevStep,
    resetForm,
    createExchange,
    validateAddress,
  };
};

// React import for useMemo and useState
import React from 'react';
