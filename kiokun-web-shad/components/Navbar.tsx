"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
} from "@/components/ui/navigation-menu";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Menu, Search } from "lucide-react";
import { ModeToggle } from "@/components/mode-toggle";

export const Navbar = () => {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const router = useRouter();

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      router.push(`/word/${encodeURIComponent(searchQuery.trim())}`);
      setSearchQuery("");
    }
  };

  return (
    <header className="sticky border-b top-0 z-40 w-full bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <NavigationMenu className="mx-auto">
        <NavigationMenuList className="container h-16 px-4 w-screen flex justify-between items-center">
          {/* Logo */}
          <NavigationMenuItem className="font-bold flex">
            <Link
              href="/"
              className="ml-2 font-bold text-xl flex items-center text-foreground hover:text-primary transition-colors"
            >
              📚 Kiokun Dictionary
            </Link>
          </NavigationMenuItem>

          {/* Desktop Search */}
          <div className="hidden md:flex flex-1 max-w-md mx-8">
            <form onSubmit={handleSearch} className="flex w-full gap-2">
              <Input
                type="text"
                placeholder="Search for a word in Japanese or Chinese..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="flex-1"
              />
              <Button type="submit" size="sm" className="px-3">
                <Search className="h-4 w-4" />
              </Button>
            </form>
          </div>

          {/* Mobile Menu */}
          <div className="flex md:hidden">
            <Sheet open={isOpen} onOpenChange={setIsOpen}>
              <SheetTrigger asChild>
                <Button variant="ghost" size="sm" className="px-2">
                  <Menu className="h-5 w-5" />
                  <span className="sr-only">Menu</span>
                </Button>
              </SheetTrigger>

              <SheetContent side="left" className="w-80">
                <SheetHeader>
                  <SheetTitle className="font-bold text-xl text-left">
                    📚 Kiokun Dictionary
                  </SheetTitle>
                </SheetHeader>

                {/* Mobile Search */}
                <div className="mt-6">
                  <form onSubmit={handleSearch} className="flex gap-2">
                    <Input
                      type="text"
                      placeholder="Search for a word..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="flex-1"
                    />
                    <Button type="submit" size="sm" className="px-3">
                      <Search className="h-4 w-4" />
                    </Button>
                  </form>
                </div>

                <div className="mt-6 space-y-4">
                  <div className="text-sm text-muted-foreground">
                    <p>Examples: 水 (water), 日本 (Japan), ありがとう (thank you), 学生 (student)</p>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">Theme</span>
                    <ModeToggle />
                  </div>
                </div>
              </SheetContent>
            </Sheet>
          </div>

          {/* Desktop Examples and Theme Toggle */}
          <div className="hidden lg:flex items-center gap-4">
            <div className="text-xs text-muted-foreground">
              Examples: 水, 日本, ありがとう, 学生
            </div>
            <ModeToggle />
          </div>
        </NavigationMenuList>
      </NavigationMenu>
    </header>
  );
};
