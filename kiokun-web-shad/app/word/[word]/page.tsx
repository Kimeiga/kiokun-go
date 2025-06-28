/**
 * Dictionary lookup page
 */

import { notFound } from 'next/navigation';
import { Suspense } from 'react';
import { Layout } from '@/components/Layout';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { DictionaryStreamingResults } from '@/components/dictionary/DictionaryStreamingResults';
import { extractShardType } from '@/lib/dictionary-utils';

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
  const shardType = extractShardType(decodedWord);

  const getShardDescription = (shard: number) => {
    switch (shard) {
      case 0: return "Non-Han (no Han characters)";
      case 1: return "Han 1 Character";
      case 2: return "Han 2 Characters";
      case 3: return "Han 3+ Characters";
      default: return `Shard ${shard}`;
    }
  };

  return (
    <Layout>
      <div className="space-y-6">
        {/* Search result header */}
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-3xl font-bold">{decodedWord}</CardTitle>
              <Badge variant="secondary">
                {getShardDescription(shardType)}
              </Badge>
            </div>
            <p className="text-muted-foreground">
              Search results for &ldquo;{decodedWord}&rdquo;
            </p>
          </CardHeader>
        </Card>

        {/* Streaming results */}
        <Suspense fallback={
          <div className="flex justify-center items-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
          </div>
        }>
          <DictionaryStreamingResults word={decodedWord} />
        </Suspense>
      </div>
    </Layout>
  );
}
