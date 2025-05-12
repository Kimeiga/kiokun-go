/**
 * Dictionary lookup page with query parameter support
 */

import { Suspense } from 'react';
import { redirect } from 'next/navigation';
import DictionaryResults from '@/components/DictionaryResults';
import LoadingResults from '@/components/LoadingResults';
import SearchForm from '@/components/SearchForm';

interface WordPageProps {
  searchParams: { word?: string };
}

export default async function WordPage({ searchParams }: WordPageProps) {
  const { word } = searchParams;
  
  // If no word is provided, show the search form
  if (!word) {
    return (
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold mb-6">Dictionary Lookup</h1>
        <SearchForm />
      </div>
    );
  }
  
  // Redirect to the dynamic route for better SEO
  redirect(`/word/${encodeURIComponent(word)}`);
}
