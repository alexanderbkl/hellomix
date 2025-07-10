import React from 'react';

const HeroSection: React.FC = () => {
  return (
    <section className="py-16 px-6 text-center">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-5xl md:text-7xl font-bold mb-6 leading-tight">
          Anonymous Crypto Exchange{' '}
          <span className="bg-gradient-to-r from-indigo-400 via-purple-400 to-amber-400 bg-clip-text text-transparent">
            Made Simple
          </span>
        </h1>
        
        <p className="text-xl md:text-2xl text-gray-300 mb-12 max-w-3xl mx-auto leading-relaxed">
          Send Bitcoin, receive any cryptocurrency anonymously across multiple wallets. 
          No registration, no KYC, no traces.
        </p>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-2xl mx-auto">
          <div className="bg-[#232340]/50 backdrop-blur-sm p-6 rounded-2xl border border-gray-700/30">
            <div className="text-3xl font-bold text-indigo-400 mb-2">$2.5M+</div>
            <div className="text-gray-400 text-lg">Mixed</div>
          </div>
          
          <div className="bg-[#232340]/50 backdrop-blur-sm p-6 rounded-2xl border border-gray-700/30">
            <div className="text-3xl font-bold text-purple-400 mb-2">1,247</div>
            <div className="text-gray-400 text-lg">Transactions</div>
          </div>
          
          <div className="bg-[#232340]/50 backdrop-blur-sm p-6 rounded-2xl border border-gray-700/30">
            <div className="text-3xl font-bold text-amber-400 mb-2">99.9%</div>
            <div className="text-gray-400 text-lg">Uptime</div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default HeroSection;
