'use client';

/**
 * Search form component for dictionary lookups
 */

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';

export default function SearchForm() {
  const [searchTerm, setSearchTerm] = useState('');
  const router = useRouter();
  
  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    
    if (searchTerm.trim()) {
      // Navigate to the word page
      router.push(`/word/${encodeURIComponent(searchTerm.trim())}`);
    }
  };
  
  return (
    <form onSubmit={handleSubmit} className="w-full max-w-md">
      <div className="flex items-center border-b border-gray-300 py-2">
        <input
          type="text"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Enter a word to look up..."
          className="appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none"
        />
        <button
          type="submit"
          className="flex-shrink-0 bg-blue-500 hover:bg-blue-700 border-blue-500 hover:border-blue-700 text-sm border-4 text-white py-1 px-2 rounded"
          disabled={!searchTerm.trim()}
        >
          Search
        </button>
      </div>
    </form>
  );
}
