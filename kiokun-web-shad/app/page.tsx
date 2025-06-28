import Link from "next/link";
import { Layout } from "@/components/Layout";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import SearchForm from "@/components/SearchForm.client";

export default function Home() {
  return (
    <Layout>
      <div className="flex flex-col items-center space-y-12">
        {/* Hero section */}
        <div className="text-center space-y-6 max-w-4xl">
          <h1 className="text-4xl md:text-6xl font-bold">
            <span className="bg-gradient-to-r from-primary to-blue-600 text-transparent bg-clip-text">
              Kiokun
            </span>{" "}
            Dictionary
          </h1>

          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Look up words in Japanese and Chinese dictionaries with real-time streaming results
          </p>

          <div className="w-full max-w-2xl mx-auto">
            <SearchForm />
          </div>

          {/* Example searches */}
          <div className="space-y-4">
            <h3 className="text-lg font-semibold">Popular Searches:</h3>
            <div className="flex flex-wrap justify-center gap-3">
              <Button variant="outline" asChild>
                <Link href="/word/日本">
                  日本 <span className="text-muted-foreground ml-1">(Japan)</span>
                </Link>
              </Button>
              <Button variant="outline" asChild>
                <Link href="/word/水">
                  水 <span className="text-muted-foreground ml-1">(water)</span>
                </Link>
              </Button>
              <Button variant="outline" asChild>
                <Link href="/word/ありがとう">
                  ありがとう <span className="text-muted-foreground ml-1">(thank you)</span>
                </Link>
              </Button>
              <Button variant="outline" asChild>
                <Link href="/word/学生">
                  学生 <span className="text-muted-foreground ml-1">(student)</span>
                </Link>
              </Button>
              <Button variant="outline" asChild>
                <Link href="/word/図書館">
                  図書館 <span className="text-muted-foreground ml-1">(library)</span>
                </Link>
              </Button>
              <Button variant="outline" asChild>
                <Link href="/word/中国">
                  中国 <span className="text-muted-foreground ml-1">(China)</span>
                </Link>
              </Button>
            </div>
          </div>
        </div>

        {/* Dictionary info */}
        <div className="w-full grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl">
          <Card>
            <CardHeader>
              <CardTitle>Dictionary Sources</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-3">
                <li className="flex items-start gap-3">
                  <Badge variant="secondary" className="bg-blue-100 text-blue-800">JMdict</Badge>
                  <span className="text-sm">Japanese words with English definitions</span>
                </li>
                <li className="flex items-start gap-3">
                  <Badge variant="secondary" className="bg-purple-100 text-purple-800">JMnedict</Badge>
                  <span className="text-sm">Japanese proper names</span>
                </li>
                <li className="flex items-start gap-3">
                  <Badge variant="secondary" className="bg-green-100 text-green-800">Kanjidic</Badge>
                  <span className="text-sm">Japanese kanji characters</span>
                </li>
                <li className="flex items-start gap-3">
                  <Badge variant="secondary" className="bg-red-100 text-red-800">Chinese Chars</Badge>
                  <span className="text-sm">Chinese characters (Hanzi)</span>
                </li>
                <li className="flex items-start gap-3">
                  <Badge variant="secondary" className="bg-yellow-100 text-yellow-800">Chinese Words</Badge>
                  <span className="text-sm">Chinese vocabulary</span>
                </li>
              </ul>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Dictionary Structure</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="mb-4 text-sm text-muted-foreground">
                The dictionary is sharded based on the number of Han characters in each word:
              </p>
              <ul className="space-y-3">
                <li className="flex items-start gap-3">
                  <Badge variant="outline">Non-Han</Badge>
                  <span className="text-sm">Words with no Han characters (e.g., ありがとう)</span>
                </li>
                <li className="flex items-start gap-3">
                  <Badge variant="outline">Han-1char</Badge>
                  <span className="text-sm">Words with exactly 1 Han character (e.g., 水)</span>
                </li>
                <li className="flex items-start gap-3">
                  <Badge variant="outline">Han-2char</Badge>
                  <span className="text-sm">Words with exactly 2 Han characters (e.g., 日本)</span>
                </li>
                <li className="flex items-start gap-3">
                  <Badge variant="outline">Han-3plus</Badge>
                  <span className="text-sm">Words with 3 or more Han characters (e.g., 図書館)</span>
                </li>
              </ul>
            </CardContent>
          </Card>
        </div>
      </div>
    </Layout>
  );
}
