"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search } from "lucide-react";

export default function SearchForm() {
  const [word, setWord] = useState("");
  const router = useRouter();

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (word.trim()) {
      router.push(`/word/${encodeURIComponent(word.trim())}`);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full space-y-3">
      <div className="flex gap-2">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            type="text"
            value={word}
            onChange={(e) => setWord(e.target.value)}
            className="pl-10 text-base h-12"
            placeholder="Search for a word in Japanese or Chinese..."
            required
          />
        </div>
        <Button type="submit" size="lg" className="h-12 px-6">
          Search
        </Button>
      </div>

      <div className="text-sm text-muted-foreground text-center">
        Examples: 水 (water), 日本 (Japan), ありがとう (thank you), 学生 (student)
      </div>
    </form>
  );
}
