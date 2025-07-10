import React from 'react';
import { Shield } from 'lucide-react';

const Header: React.FC = () => {
  return (
    <header className="flex justify-between items-center p-6 bg-[#1a1a2e]/80 backdrop-blur-xl border-b border-gray-700/50 sticky top-0 z-50">
      <div className="flex items-center gap-3">
        <Shield className="text-indigo-400 w-8 h-8" />
        <span className="text-2xl font-bold bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
          HelloMix
        </span>
      </div>
      
      <nav className="hidden md:flex gap-8">
        <a 
          href="#how-it-works" 
          className="text-gray-300 hover:text-indigo-400 transition-colors duration-300 font-medium"
        >
          How it Works
        </a>
        <a 
          href="#fees" 
          className="text-gray-300 hover:text-indigo-400 transition-colors duration-300 font-medium"
        >
          Fees
        </a>
        <a 
          href="#support" 
          className="text-gray-300 hover:text-indigo-400 transition-colors duration-300 font-medium"
        >
          Support
        </a>
      </nav>

      {/* Mobile menu button */}
      <button className="md:hidden text-gray-300 hover:text-white">
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </button>
    </header>
  );
};

export default Header;
