'use client';

import React, { useState } from 'react';
import { 
  Clock, 
  CheckCircle, 
  AlertCircle, 
  Copy, 
  ExternalLink, 
  RefreshCw,
  Bitcoin,
  TrendingUp
} from 'lucide-react';
import { useTransactionStatus } from '@/hooks/useApi';
import { formatCurrency, truncateAddress, copyToClipboard } from '@/lib/utils';

interface TransactionTrackerProps {
  transactionId: string;
  onBack?: () => void;
}

const TransactionTracker: React.FC<TransactionTrackerProps> = ({ transactionId, onBack }) => {
  const [copySuccess, setCopySuccess] = useState<{ [key: string]: boolean }>({});
  
  const { data: transactionData, isLoading, error, refetch } = useTransactionStatus(transactionId);
  const transaction = transactionData?.data;

  const handleCopy = async (text: string, type: string) => {
    const success = await copyToClipboard(text);
    if (success) {
      setCopySuccess({ ...copySuccess, [type]: true });
      setTimeout(() => {
        setCopySuccess({ ...copySuccess, [type]: false });
      }, 2000);
    }
  };

  const getStatusDetails = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'pending':
        return {
          icon: Clock,
          color: 'text-yellow-500',
          bgColor: 'bg-yellow-500/10',
          borderColor: 'border-yellow-500/30',
          message: 'Waiting for payment confirmation'
        };
      case 'processing':
        return {
          icon: RefreshCw,
          color: 'text-blue-500',
          bgColor: 'bg-blue-500/10',
          borderColor: 'border-blue-500/30',
          message: 'Processing your exchange'
        };
      case 'completed':
        return {
          icon: CheckCircle,
          color: 'text-green-500',
          bgColor: 'bg-green-500/10',
          borderColor: 'border-green-500/30',
          message: 'Exchange completed successfully'
        };
      case 'failed':
        return {
          icon: AlertCircle,
          color: 'text-red-500',
          bgColor: 'bg-red-500/10',
          borderColor: 'border-red-500/30',
          message: 'Exchange failed - contact support'
        };
      default:
        return {
          icon: Clock,
          color: 'text-gray-500',
          bgColor: 'bg-gray-500/10',
          borderColor: 'border-gray-500/30',
          message: 'Unknown status'
        };
    }
  };

  if (isLoading) {
    return (
      <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
        <div className="flex items-center justify-center space-x-2">
          <RefreshCw className="w-6 h-6 animate-spin text-indigo-500" />
          <span className="text-gray-300">Loading transaction details...</span>
        </div>
      </div>
    );
  }

  if (error || !transaction) {
    return (
      <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
        <div className="text-center">
          <AlertCircle className="w-12 h-12 text-red-500 mx-auto mb-4" />
          <h3 className="text-xl font-bold text-white mb-2">Transaction Not Found</h3>
          <p className="text-gray-400 mb-4">
            Could not find transaction with ID: {transactionId}
          </p>
          {onBack && (
            <button
              onClick={onBack}
              className="bg-indigo-500 hover:bg-indigo-600 text-white px-6 py-2 rounded-lg transition-colors"
            >
              Go Back
            </button>
          )}
        </div>
      </div>
    );
  }

  const statusDetails = getStatusDetails(transaction.status);
  const StatusIcon = statusDetails.icon;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-white">Transaction Status</h2>
          <button
            onClick={() => refetch()}
            className="flex items-center space-x-2 text-gray-400 hover:text-white transition-colors"
          >
            <RefreshCw className="w-4 h-4" />
            <span className="text-sm">Refresh</span>
          </button>
        </div>

        {/* Status Badge */}
        <div className={`inline-flex items-center space-x-2 px-4 py-2 rounded-full ${statusDetails.bgColor} ${statusDetails.borderColor} border`}>
          <StatusIcon className={`w-5 h-5 ${statusDetails.color} ${transaction.status === 'processing' ? 'animate-spin' : ''}`} />
          <span className={`font-medium ${statusDetails.color}`}>
            {transaction.status.charAt(0).toUpperCase() + transaction.status.slice(1)}
          </span>
        </div>

        <p className="text-gray-400 mt-2">{statusDetails.message}</p>
      </div>

      {/* Transaction Details */}
      <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
        <h3 className="text-lg font-semibold text-white mb-4">Transaction Details</h3>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-3">
            <div>
              <label className="text-sm text-gray-400">Transaction ID</label>
              <div className="flex items-center space-x-2 mt-1">
                <span className="font-mono text-sm text-white">{truncateAddress(transaction.id, 8, 8)}</span>
                <button
                  onClick={() => handleCopy(transaction.id, 'id')}
                  className="text-gray-400 hover:text-white transition-colors"
                >
                  <Copy className="w-4 h-4" />
                </button>
                {copySuccess.id && (
                  <span className="text-green-400 text-sm">Copied!</span>
                )}
              </div>
            </div>

            <div>
              <label className="text-sm text-gray-400">Amount to Send</label>
              <div className="flex items-center space-x-2 mt-1">
                <Bitcoin className="w-4 h-4 text-amber-400" />
                <span className="font-semibold text-white">{formatCurrency(transaction.btc_amount, 8)} BTC</span>
              </div>
            </div>

            <div>
              <label className="text-sm text-gray-400">Output Currency</label>
              <div className="flex items-center space-x-2 mt-1">
                <span className="font-semibold text-white">{transaction.output_currency}</span>
                <TrendingUp className="w-4 h-4 text-green-400" />
              </div>
            </div>
          </div>

          <div className="space-y-3">
            <div>
              <label className="text-sm text-gray-400">Expected Output</label>
              <p className="font-semibold text-white mt-1">
                {formatCurrency(transaction.estimated_output, 6)} {transaction.output_currency}
              </p>
            </div>

            <div>
              <label className="text-sm text-gray-400">Fee</label>
              <p className="font-semibold text-white mt-1">
                ${formatCurrency(transaction.fee, 2)}
              </p>
            </div>

            <div>
              <label className="text-sm text-gray-400">Created</label>
              <p className="text-white mt-1">
                {new Date(transaction.created_at).toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Payment Address */}
      <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
        <h3 className="text-lg font-semibold text-white mb-4">Payment Address</h3>
        
        <div className="bg-gray-800/50 p-4 rounded-xl">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm text-gray-400">Send Bitcoin to this address</span>
            <button
              onClick={() => handleCopy(transaction.payment_address, 'address')}
              className="flex items-center space-x-1 text-indigo-400 hover:text-indigo-300 transition-colors"
            >
              <Copy className="w-4 h-4" />
              <span className="text-sm">Copy</span>
            </button>
          </div>
          
          <div className="font-mono text-sm text-white break-all bg-gray-900/50 p-3 rounded border">
            {transaction.payment_address}
          </div>
          
          {copySuccess.address && (
            <div className="text-green-400 text-sm mt-2">✓ Address copied to clipboard!</div>
          )}
        </div>
      </div>

      {/* Output Addresses */}
      <div className="bg-[#232340]/50 backdrop-blur-sm p-8 rounded-3xl border border-gray-700/30">
        <h3 className="text-lg font-semibold text-white mb-4">Output Addresses</h3>
        
        <div className="space-y-3">
          {transaction.output_addresses.map((addr, index) => (
            <div key={index} className="flex items-center justify-between p-3 bg-gray-800/30 rounded-lg">
              <div className="flex-1">
                <div className="font-mono text-sm text-white">
                  {truncateAddress(addr.address, 12, 12)}
                </div>
                <div className="text-xs text-gray-400 mt-1">
                  {addr.percentage}% • {formatCurrency((transaction.estimated_output * addr.percentage) / 100, 6)} {transaction.output_currency}
                </div>
              </div>
              <button
                onClick={() => handleCopy(addr.address, `output-${index}`)}
                className="text-gray-400 hover:text-white transition-colors ml-2"
              >
                <Copy className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>
      </div>

      {/* Important Notes */}
      <div className="bg-blue-500/10 border border-blue-500/30 p-6 rounded-2xl">
        <div className="flex items-start space-x-3">
          <AlertCircle className="w-5 h-5 text-blue-400 mt-0.5" />
          <div>
            <h4 className="font-semibold text-blue-400 mb-2">Important Notes</h4>
            <ul className="text-gray-300 text-sm space-y-1">
              <li>• Send only Bitcoin (BTC) to the payment address</li>
              <li>• Send exactly {formatCurrency(transaction.btc_amount, 8)} BTC</li>
              <li>• Allow 1-6 confirmations for processing</li>
              <li>• Do not send from an exchange - use a personal wallet</li>
              <li>• Keep this page open to monitor progress</li>
            </ul>
          </div>
        </div>
      </div>

      {/* Actions */}
      <div className="flex justify-center space-x-4">
        {onBack && (
          <button
            onClick={onBack}
            className="px-6 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            New Exchange
          </button>
        )}
        
        <button
          onClick={() => window.open(`https://blockchair.com/bitcoin/address/${transaction.payment_address}`, '_blank')}
          className="flex items-center space-x-2 px-6 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-lg transition-colors"
        >
          <ExternalLink className="w-4 h-4" />
          <span>View on Blockchain</span>
        </button>
      </div>
    </div>
  );
};

export default TransactionTracker;
