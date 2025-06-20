/**
 * Dictionary lookup page
 */

import { notFound } from 'next/navigation';
import Link from 'next/link';
import { Suspense } from 'react';
import JishoDictionaryStreamingResults from '@/components/JishoDictionaryStreamingResults';
import SearchForm from '@/components/SearchForm.client';

interface WordPageProps {
  params: Promise<{ word: string }>;
}

export default async function WordPage({ params }: WordPageProps) {
  const { word } = await params;

  if (!word) {
    notFound();
  }

  // Decode the word parameter (it comes URL-encoded from the route)
  const decodedWord = decodeURIComponent(word);

  return (
    <div className="min-h-screen bg-[#121212]">
      {/* Header with search */}
      <header className="bg-[#1e1e1e] border-b border-gray-800 sticky top-0 z-10 shadow-sm">
        <div className="container mx-auto px-4 py-4">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
            <Link href="/" className="text-2xl font-bold text-blue-500 hover:text-blue-400 transition-colors">
              Kiokun Dictionary
            </Link>

            <div className="w-full md:w-1/2 lg:w-2/3">
              <SearchForm />
            </div>
          </div>
        </div>
      </header>

      {/* Main content */}
      <main className="container mx-auto px-4 py-8">
        <div className="bg-[#1e1e1e] rounded-lg shadow-sm border border-gray-800 p-6 mb-6">
          <h1 className="text-3xl font-bold mb-2 text-white">{decodedWord}</h1>
          <div className="text-gray-400">
            Search results for &ldquo;{decodedWord}&rdquo;
          </div>
        </div>

        <Suspense fallback={
          <div className="flex justify-center items-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-400"></div>
          </div>
        }>
          <JishoDictionaryStreamingResults word={decodedWord} />
        </Suspense>
      </main>

      {/* Footer */}
      <footer className="mt-12 py-6 border-t border-gray-800">
        <div className="container mx-auto px-4 text-center text-gray-400 text-sm">
          <p>Powered by Next.js and jsDelivr</p>
        </div>
      </footer>
    </div>
  );
}
