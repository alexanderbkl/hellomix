'use client';

import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { 
  Bitcoin, 
  ArrowRight, 
  Plus,
  Trash2,
  Copy,
  QrCode,
  CheckCircle,
  Clock,
  AlertCircle,
  Loader2
} from 'lucide-react';
import QRCode from 'qrcode';
import { 
  usePrices, 
  useSupportedCurrencies, 
  useCreateExchange, 
  useTransactionStatus,
  usePriceCalculation
} from '@/hooks/useApi';
import { 
  formatCurrency, 
  formatPrice, 
  truncateAddress, 
  getCurrencyIcon,
  getCurrencyName,
  validateCryptoAddress,
  copyToClipboard
} from '@/lib/utils';

interface OutputAddress {
  address: string;
  percentage: number;
}

interface ExchangeFormData {
  btcAmount: number;
  outputCurrency: string;
  outputAddresses: OutputAddress[];
}

interface Transaction {
  id: string;
  payment_address: string;
  btc_amount: number;
  output_currency: string;
  output_addresses: OutputAddress[];
  estimated_output: number;
  fee: number;
  status: string;
  created_at: string;
}

const ExchangeForm: React.FC = () => {
  const [currentStep, setCurrentStep] = useState(1);
  const [transaction, setTransaction] = useState<Transaction | null>(null);
  const [qrCodeUrl, setQrCodeUrl] = useState<string>('');
  const [addresses, setAddresses] = useState<OutputAddress[]>([
    { address: '', percentage: 100 }
  ]);

  const { register, handleSubmit, watch, setValue, formState: { errors } } = useForm<ExchangeFormData>({
    defaultValues: {
      btcAmount: 0.001,
      outputCurrency: '',
      outputAddresses: addresses
    }
  });

  const watchedAmount = watch('btcAmount');
  const watchedCurrency = watch('outputCurrency');

  // API hooks
  const { data: prices, isLoading: pricesLoading } = usePrices();
  const { data: supportedCurrencies } = useSupportedCurrencies();
  const createExchange = useCreateExchange();

  const addAddress = () => {
    if (addresses.length < 7) {
      const newPercentage = Math.floor(100 / (addresses.length + 1));
      const updatedAddresses = addresses.map(addr => ({ ...addr, percentage: newPercentage }));
      updatedAddresses.push({ address: '', percentage: newPercentage });
      setAddresses(updatedAddresses);
      setValue('outputAddresses', updatedAddresses);
    }
  };

  const removeAddress = (index: number) => {
    if (addresses.length > 1) {
      const updatedAddresses = addresses.filter((_, i) => i !== index);
      // Redistribute percentages
      const equalPercentage = Math.floor(100 / updatedAddresses.length);
      const redistributed = updatedAddresses.map(addr => ({ ...addr, percentage: equalPercentage }));
      setAddresses(redistributed);
      setValue('outputAddresses', redistributed);
    }
  };

  const updateAddress = (index: number, field: keyof OutputAddress, value: string | number) => {
    const updatedAddresses = [...addresses];
    updatedAddresses[index] = { ...updatedAddresses[index], [field]: value };
    setAddresses(updatedAddresses);
    setValue('outputAddresses', updatedAddresses);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const calculateEstimatedOutput = () => {
    if (!prices || !watchedCurrency || !watchedAmount) return 0;
    const btcPrice = prices.BTC || 0;
    const outputPrice = prices[watchedCurrency] || 0;
    if (btcPrice === 0 || outputPrice === 0) return 0;
    return (watchedAmount * btcPrice) / outputPrice;
  };

  const calculateUSDValue = () => {
    if (!prices || !watchedAmount) return 0;
    return watchedAmount * (prices.BTC || 0);
  };

  const onSubmit = (data: ExchangeFormData) => {
    if (currentStep === 3) {
      createExchange.mutate({
        btc_amount: data.btcAmount,
        output_currency: data.outputCurrency,
        output_addresses: data.outputAddresses
      }, {
        onSuccess: async (response) => {
          setTransaction(response.data);
          // Generate QR code for payment address
          try {
            const qrUrl = await QRCode.toDataURL(response.data.payment_address);
            setQrCodeUrl(qrUrl);
          } catch (error) {
            console.error('Error generating QR code:', error);
          }
          setCurrentStep(4);
        },
        onError: (error) => {
          console.error('Error creating exchange:', error);
        }
      });
    } else {
      setCurrentStep(currentStep + 1);
    }
  };

  const currencies = [
    { id: 'ETH', name: 'Ethereum', icon: '⟠' },
    { id: 'USDT', name: 'Tether', icon: '₮' },
    { id: 'USDC', name: 'USD Coin', icon: '$' },
    { id: 'ADA', name: 'Cardano', icon: '₳' },
    { id: 'SOL', name: 'Solana', icon: '◎' },
    { id: 'MATIC', name: 'Polygon', icon: '⬟' }
  ];

  return (
    <div className="max-w-4xl mx-auto">
      {/* Progress Steps */}
      <div className="flex justify-center mb-8">
        <div className="flex items-center space-x-4">
          {[1, 2, 3, 4].map((step) => (
            <React.Fragment key={step}>
              <div className={`flex items-center justify-center w-10 h-10 rounded-full font-medium ${
                step <= currentStep 
                  ? 'bg-indigo-500 text-white' 
                  : 'bg-gray-700 text-gray-400'
              }`}>
                {step}
              </div>
              {step < 4 && <ArrowRight className={`w-5 h-5 ${
                step < currentStep ? 'text-indigo-500' : 'text-gray-600'
              }`} />}
            </React.Fragment>
          ))}
        </div>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
        {/* Step 1: Amount Input */}
        {currentStep === 1 && (
          <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
            <h2 className="text-2xl font-bold mb-6 text-center">Enter Bitcoin Amount</h2>
            
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-2">
                  Bitcoin Amount
                </label>
                <div className="relative">
                  <input
                    type="number"
                    step="0.00000001"
                    min="0.001"
                    max="10"
                    {...register('btcAmount', { 
                      required: 'Bitcoin amount is required',
                      min: { value: 0.001, message: 'Minimum amount is 0.001 BTC' },
                      max: { value: 10, message: 'Maximum amount is 10 BTC' }
                    })}
                    className="w-full px-4 py-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 pr-12"
                    placeholder="0.001"
                  />
                  <Bitcoin className="absolute right-3 top-3 w-6 h-6 text-amber-400" />
                </div>
                {errors.btcAmount && (
                  <p className="text-red-400 text-sm mt-1">{errors.btcAmount.message}</p>
                )}
              </div>

              <div className="text-center">
                <p className="text-gray-400">
                  ≈ ${calculateUSDValue().toLocaleString()} USD
                </p>
                {pricesLoading && (
                  <p className="text-gray-500 text-sm">Loading prices...</p>
                )}
              </div>

              <button
                type="submit"
                disabled={!watchedAmount || watchedAmount < 0.001}
                className="w-full bg-gradient-to-r from-indigo-500 to-purple-500 text-white py-4 rounded-xl font-medium hover:from-indigo-600 hover:to-purple-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Continue
              </button>
            </div>
          </div>
        )}

        {/* Step 2: Currency Selection */}
        {currentStep === 2 && (
          <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
            <h2 className="text-2xl font-bold mb-6 text-center">Select Output Currency</h2>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
              {currencies.map((currency) => (
                <button
                  key={currency.id}
                  type="button"
                  onClick={() => setValue('outputCurrency', currency.id)}
                  className={`p-4 rounded-xl border-2 transition-all duration-300 ${
                    watchedCurrency === currency.id
                      ? 'border-indigo-500 bg-indigo-500/20'
                      : 'border-gray-600 bg-gray-800/30 hover:border-gray-500'
                  }`}
                >
                  <div className="flex items-center space-x-3">
                    <span className="text-2xl">{currency.icon}</span>
                    <div className="text-left">
                      <div className="font-medium">{currency.name}</div>
                      <div className="text-sm text-gray-400">{currency.id}</div>
                      {prices && (
                        <div className="text-sm text-green-400">
                          ${prices[currency.id]?.toLocaleString() || 'N/A'}
                        </div>
                      )}
                    </div>
                  </div>
                </button>
              ))}
            </div>

            {watchedCurrency && (
              <div className="text-center mb-6">
                <p className="text-gray-400">
                  You will receive approximately{' '}
                  <span className="text-white font-medium">
                    {calculateEstimatedOutput().toFixed(6)} {watchedCurrency}
                  </span>
                </p>
              </div>
            )}

            <button
              type="submit"
              disabled={!watchedCurrency}
              className="w-full bg-gradient-to-r from-indigo-500 to-purple-500 text-white py-4 rounded-xl font-medium hover:from-indigo-600 hover:to-purple-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Continue
            </button>
          </div>
        )}

        {/* Step 3: Output Addresses */}
        {currentStep === 3 && (
          <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
            <h2 className="text-2xl font-bold mb-6 text-center">
              Output Addresses for {watchedCurrency}
            </h2>
            
            <div className="space-y-4 mb-6">
              {addresses.map((addr, index) => (
                <div key={index} className="flex gap-2">
                  <div className="flex-1">
                    <input
                      type="text"
                      placeholder={`Enter ${watchedCurrency} address`}
                      value={addr.address}
                      onChange={(e) => updateAddress(index, 'address', e.target.value)}
                      className="w-full px-4 py-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white placeholder-gray-400 focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
                    />
                  </div>
                  <div className="w-24">
                    <input
                      type="number"
                      min="1"
                      max="100"
                      value={addr.percentage}
                      onChange={(e) => updateAddress(index, 'percentage', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-3 bg-gray-800/50 border border-gray-600 rounded-xl text-white text-center focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500"
                    />
                  </div>
                  {addresses.length > 1 && (
                    <button
                      type="button"
                      onClick={() => removeAddress(index)}
                      className="p-3 text-red-400 hover:text-red-300 hover:bg-red-400/10 rounded-xl transition-all"
                    >
                      <Trash2 className="w-5 h-5" />
                    </button>
                  )}
                </div>
              ))}
            </div>

            {addresses.length < 7 && (
              <button
                type="button"
                onClick={addAddress}
                className="flex items-center gap-2 text-indigo-400 hover:text-indigo-300 mb-6 mx-auto"
              >
                <Plus className="w-4 h-4" />
                Add Another Address
              </button>
            )}

            <div className="text-center mb-6">
              <p className="text-gray-400">
                Total: {addresses.reduce((sum, addr) => sum + addr.percentage, 0)}%
              </p>
            </div>

            <button
              type="submit"
              disabled={
                createExchange.isPending ||
                addresses.some(addr => !addr.address) ||
                addresses.reduce((sum, addr) => sum + addr.percentage, 0) !== 100
              }
              className="w-full bg-gradient-to-r from-indigo-500 to-purple-500 text-white py-4 rounded-xl font-medium hover:from-indigo-600 hover:to-purple-600 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {createExchange.isPending ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin" />
                  Creating Exchange...
                </>
              ) : (
                'Create Exchange'
              )}
            </button>
          </div>
        )}

        {/* Step 4: Payment Instructions */}
        {currentStep === 4 && transaction && (
          <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
            <h2 className="text-2xl font-bold mb-6 text-center">Payment Instructions</h2>
            
            <div className="space-y-6">
              <div className="text-center">
                <p className="text-gray-400 mb-4">Send exactly</p>
                <p className="text-3xl font-bold text-amber-400">
                  {transaction.btc_amount} BTC
                </p>
                <p className="text-gray-400">to the address below</p>
              </div>

              <div className="bg-gray-800/50 p-6 rounded-xl">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-sm font-medium text-gray-300">Payment Address</span>
                  <button
                    onClick={() => copyToClipboard(transaction.payment_address)}
                    className="flex items-center gap-2 text-indigo-400 hover:text-indigo-300 text-sm"
                  >
                    <Copy className="w-4 h-4" />
                    Copy
                  </button>
                </div>
                <p className="font-mono text-sm break-all text-white bg-gray-900/50 p-3 rounded">
                  {transaction.payment_address}
                </p>
              </div>

              {qrCodeUrl && (
                <div className="text-center">
                  <img src={qrCodeUrl} alt="Payment QR Code" className="mx-auto rounded-xl" />
                  <p className="text-sm text-gray-400 mt-2">Scan with your Bitcoin wallet</p>
                </div>
              )}

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div className="bg-gray-800/30 p-4 rounded-xl">
                  <p className="text-gray-400">Transaction ID</p>
                  <p className="font-mono text-white break-all">{transaction.id}</p>
                </div>
                <div className="bg-gray-800/30 p-4 rounded-xl">
                  <p className="text-gray-400">Status</p>
                  <div className="flex items-center gap-2">
                    <Clock className="w-4 h-4 text-yellow-400" />
                    <span className="text-yellow-400 capitalize">{transaction.status}</span>
                  </div>
                </div>
              </div>

              <div className="bg-blue-500/10 border border-blue-500/30 p-4 rounded-xl">
                <div className="flex items-start gap-3">
                  <AlertCircle className="w-5 h-5 text-blue-400 mt-0.5" />
                  <div className="text-sm">
                    <p className="text-blue-400 font-medium mb-1">Important Notes:</p>
                    <ul className="text-gray-300 space-y-1">
                      <li>• Send only Bitcoin (BTC) to this address</li>
                      <li>• Send the exact amount specified</li>
                      <li>• Allow 1-6 confirmations for processing</li>
                      <li>• Keep this page open to track progress</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </form>
    </div>
  );
};

export default ExchangeForm;
