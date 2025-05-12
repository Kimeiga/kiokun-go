import SearchForm from "@/components/SearchForm.client";

export default function Home() {
  return (
    <div className="flex flex-col items-center min-h-screen p-4 md:p-8 bg-gradient-to-b from-white to-gray-50">
      <header className="w-full max-w-5xl py-6 flex justify-center">
        <h1 className="text-4xl font-bold text-gray-800">
          <span className="text-blue-600">Kiokun</span> Dictionary
        </h1>
      </header>

      <main className="flex flex-col items-center w-full max-w-5xl">
        {/* Hero section with search */}
        <div className="w-full text-center mb-12">
          <p className="text-xl text-gray-600 mb-8">
            Look up words in Japanese and Chinese dictionaries
          </p>

          <div className="w-full max-w-2xl mx-auto mb-12">
            <SearchForm />
          </div>

          {/* Example searches */}
          <div className="mb-12">
            <h3 className="text-lg font-semibold mb-4">Popular Searches:</h3>
            <div className="flex flex-wrap justify-center gap-3">
              <a href="/word/日本" className="bg-blue-50 hover:bg-blue-100 text-blue-800 px-4 py-2 rounded-lg shadow-sm transition-colors">
                日本 <span className="text-gray-500 text-sm">(Japan)</span>
              </a>
              <a href="/word/水" className="bg-blue-50 hover:bg-blue-100 text-blue-800 px-4 py-2 rounded-lg shadow-sm transition-colors">
                水 <span className="text-gray-500 text-sm">(water)</span>
              </a>
              <a href="/word/ありがとう" className="bg-blue-50 hover:bg-blue-100 text-blue-800 px-4 py-2 rounded-lg shadow-sm transition-colors">
                ありがとう <span className="text-gray-500 text-sm">(thank you)</span>
              </a>
              <a href="/word/学生" className="bg-blue-50 hover:bg-blue-100 text-blue-800 px-4 py-2 rounded-lg shadow-sm transition-colors">
                学生 <span className="text-gray-500 text-sm">(student)</span>
              </a>
              <a href="/word/図書館" className="bg-blue-50 hover:bg-blue-100 text-blue-800 px-4 py-2 rounded-lg shadow-sm transition-colors">
                図書館 <span className="text-gray-500 text-sm">(library)</span>
              </a>
              <a href="/word/中国" className="bg-blue-50 hover:bg-blue-100 text-blue-800 px-4 py-2 rounded-lg shadow-sm transition-colors">
                中国 <span className="text-gray-500 text-sm">(China)</span>
              </a>
            </div>
          </div>

          {/* Dictionary info */}
          <div className="w-full grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
              <h2 className="text-xl font-semibold mb-4">Dictionary Sources</h2>
              <ul className="space-y-2">
                <li className="flex items-start">
                  <span className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">JMdict</span>
                  <span>Japanese words with English definitions</span>
                </li>
                <li className="flex items-start">
                  <span className="bg-purple-100 text-purple-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">JMnedict</span>
                  <span>Japanese proper names</span>
                </li>
                <li className="flex items-start">
                  <span className="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">Kanjidic</span>
                  <span>Japanese kanji characters</span>
                </li>
                <li className="flex items-start">
                  <span className="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">Chinese Chars</span>
                  <span>Chinese characters (Hanzi)</span>
                </li>
                <li className="flex items-start">
                  <span className="bg-yellow-100 text-yellow-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">Chinese Words</span>
                  <span>Chinese vocabulary</span>
                </li>
              </ul>
            </div>

            <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
              <h2 className="text-xl font-semibold mb-4">Dictionary Structure</h2>
              <p className="mb-4 text-gray-700">
                The dictionary is sharded based on the number of Han characters in each word:
              </p>
              <ul className="space-y-2 text-gray-700">
                <li className="flex items-start">
                  <span className="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">Non-Han</span>
                  <span>Words with no Han characters (e.g., ありがとう)</span>
                </li>
                <li className="flex items-start">
                  <span className="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">Han-1char</span>
                  <span>Words with exactly 1 Han character (e.g., 水)</span>
                </li>
                <li className="flex items-start">
                  <span className="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">Han-2char</span>
                  <span>Words with exactly 2 Han characters (e.g., 日本)</span>
                </li>
                <li className="flex items-start">
                  <span className="bg-gray-100 text-gray-800 text-xs font-medium px-2.5 py-0.5 rounded mr-2 mt-1">Han-3plus</span>
                  <span>Words with 3 or more Han characters (e.g., 図書館)</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </main>

      <footer className="w-full max-w-5xl mt-12 py-6 text-center text-gray-500 text-sm border-t border-gray-200">
        <p>Powered by Next.js and jsDelivr</p>
      </footer>
    </div>
  );
}
