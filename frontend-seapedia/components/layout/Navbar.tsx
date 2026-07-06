"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Search, ShoppingCart, Menu, X } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { ShipWheelMark } from "@/components/ui/WaveDivider";
import { ROLE_LABEL } from "@/lib/format";

const categoryLinks = [
  { href: "/products", label: "Semua Produk" },
  { href: "/#reviews", label: "Ulasan" },
];

export function Navbar() {
  const { profile, isAuthenticated, logout } = useAuth();
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    router.push(query.trim() ? `/products?q=${encodeURIComponent(query.trim())}` : "/products");
  }

  return (
    <header className="sticky top-0 z-40 bg-navy-950">
      {/* Top thin bar */}
      <div className="hidden border-b border-white/10 md:block">
        <div className="mx-auto flex max-w-6xl items-center justify-end gap-4 px-5 py-1.5 text-xs text-sand-100/70">
          {isAuthenticated ? (
            <>
              <span>
                {profile?.name}
                {profile?.active_role && (
                  <span className="ml-1.5 rounded-full bg-white/10 px-2 py-0.5 text-white">
                    {ROLE_LABEL[profile.active_role] || profile.active_role}
                  </span>
                )}
              </span>
              <Link href="/dashboard" className="hover:text-white">Dashboard</Link>
              <button onClick={() => logout()} className="hover:text-white">Keluar</button>
            </>
          ) : (
            <>
              <Link href="/register" className="hover:text-white">Daftar</Link>
              <Link href="/login" className="hover:text-white">Masuk</Link>
            </>
          )}
        </div>
      </div>

      {/* Main bar: logo + search + cart */}
      <div className="mx-auto flex max-w-6xl items-center gap-4 px-5 py-3">
        <Link href="/" className="flex shrink-0 items-center gap-2 text-white">
          <ShipWheelMark className="h-8 w-8 text-white" />
          <span className="font-display text-2xl font-extrabold tracking-tight">SEAPEDIA</span>
        </Link>

        <form onSubmit={handleSearch} className="hidden flex-1 md:flex">
          <div className="flex w-full overflow-hidden rounded-md bg-white">
            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Cari produk, toko, atau brand..."
              className="w-full px-4 py-2.5 text-sm text-ink outline-none"
            />
            <button
              type="submit"
              className="flex items-center justify-center bg-coral-500 px-5 text-white hover:bg-coral-600"
              aria-label="Cari"
            >
              <Search className="h-4 w-4" />
            </button>
          </div>
        </form>

        <div className="hidden shrink-0 items-center gap-4 md:flex">
          <Link href="/dashboard/buyer/cart" className="flex items-center gap-1.5 text-white hover:text-teal-100">
            <ShoppingCart className="h-5 w-5" />
            <span className="text-sm font-medium">Keranjang</span>
          </Link>
        </div>

        <button
          className="ml-auto flex h-9 w-9 items-center justify-center rounded-md text-white md:hidden"
          onClick={() => setOpen((v) => !v)}
          aria-label="Buka menu"
        >
          {open ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
        </button>
      </div>

      {/* Category row */}
      <div className="hidden border-t border-white/10 md:block">
        <div className="mx-auto flex max-w-6xl gap-6 px-5 py-2 text-sm">
          {categoryLinks.map((link) => (
            <Link key={link.href} href={link.href} className="text-sand-100/80 hover:text-white">
              {link.label}
            </Link>
          ))}
        </div>
      </div>

      {open && (
        <div className="border-t border-white/10 bg-navy-950 px-5 py-4 md:hidden">
          <form onSubmit={handleSearch} className="mb-3 flex overflow-hidden rounded-md bg-white">
            <input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Cari produk..."
              className="w-full px-4 py-2 text-sm text-ink outline-none"
            />
            <button type="submit" className="bg-coral-500 px-4 text-white">
              <Search className="h-4 w-4" />
            </button>
          </form>
          <div className="flex flex-col gap-3">
            {categoryLinks.map((link) => (
              <Link key={link.href} href={link.href} className="text-sand-100" onClick={() => setOpen(false)}>
                {link.label}
              </Link>
            ))}
            <Link
              href="/dashboard/buyer/cart"
              className="flex items-center gap-1.5 text-sand-100"
              onClick={() => setOpen(false)}
            >
              <ShoppingCart className="h-4 w-4" /> Keranjang
            </Link>
            <div className="my-1 h-px bg-white/10" />
            {isAuthenticated ? (
              <>
                <Link href="/dashboard" className="text-teal-400" onClick={() => setOpen(false)}>
                  Dashboard ({profile?.name})
                </Link>
                <button className="text-left text-sand-100" onClick={() => { setOpen(false); logout(); }}>
                  Keluar
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="text-sand-100" onClick={() => setOpen(false)}>
                  Masuk
                </Link>
                <Link href="/register" className="text-teal-400" onClick={() => setOpen(false)}>
                  Daftar
                </Link>
              </>
            )}
          </div>
        </div>
      )}
    </header>
  );
}
