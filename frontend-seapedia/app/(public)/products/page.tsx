import Link from "next/link";
import {
  ShoppingBasket,
  Shirt,
  Smartphone,
  Home as HomeIcon,
  Sparkles,
  Blocks,
  BookOpen,
  Dumbbell,
  Flame,
  Zap,
  ShoppingBag,
  Store,
  Truck,
  LayoutDashboard,
  type LucideIcon,
} from "lucide-react";
import { ReviewSection } from "@/components/features/ReviewSection";
import { ProductCard } from "@/components/features/ProductCard";
import { Button } from "@/components/ui/Button";
import type { ProductResponse } from "@/lib/types";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function getProducts(): Promise<ProductResponse[]> {
  try {
    const res = await fetch(`${API_BASE}/products`, { cache: "no-store" });
    if (!res.ok) return [];
    const json = await res.json();
    const data = json?.data;
    return Array.isArray(data) ? data : [];
  } catch {
    return [];
  }
}

const categories: { icon: LucideIcon; label: string }[] = [
  { icon: ShoppingBasket, label: "Kelontong" },
  { icon: Shirt, label: "Fashion" },
  { icon: Smartphone, label: "Elektronik" },
  { icon: HomeIcon, label: "Rumah Tangga" },
  { icon: Sparkles, label: "Kecantikan" },
  { icon: Blocks, label: "Mainan" },
  { icon: BookOpen, label: "Buku" },
  { icon: Dumbbell, label: "Olahraga" },
];

const roleCards: { title: string; desc: string; icon: LucideIcon }[] = [
  { title: "Pembeli", desc: "Dompet, alamat, keranjang, dan lacak pesananmu.", icon: ShoppingBag },
  { title: "Penjual", desc: "Buka toko dan kelola produk & pesanan masuk.", icon: Store },
  { title: "Kurir", desc: "Temukan job, ambil, dan selesaikan pengiriman.", icon: Truck },
  { title: "Admin", desc: "Pantau marketplace & kelola voucher/promo.", icon: LayoutDashboard },
];

export default async function ProductsPage({
  searchParams,
}: {
  searchParams: Promise<{ q?: string }>;
}) {
  const { q } = await searchParams;
  const allProducts = await getProducts();
  const isSearching = Boolean(q);

  const products = q
    ? allProducts.filter(
        (p) =>
          p.name.toLowerCase().includes(q.toLowerCase()) ||
          p.description?.toLowerCase().includes(q.toLowerCase()) ||
          p.store_name?.toLowerCase().includes(q.toLowerCase())
      )
    : allProducts;

  return (
    <div>
      {!isSearching && (
        <>
          {/* Promo banner */}
          <section className="bg-navy-950">
            <div className="mx-auto max-w-6xl px-5 py-10">
              <div className="grid items-center gap-8 rounded-2xl bg-gradient-to-r from-navy-800 to-navy-600 p-8 text-white lg:grid-cols-[1.2fr_0.8fr]">
                <div>
                  <span className="inline-flex items-center gap-1.5 rounded-full bg-coral-500 px-3 py-1 text-xs font-bold">
                    <Flame className="h-3.5 w-3.5" /> PROMO SPESIAL
                  </span>
                  <h1 className="mt-4 font-display text-3xl font-extrabold leading-tight sm:text-4xl">
                    Belanja apa saja, dari toko mana saja.
                  </h1>
                  <p className="mt-3 max-w-md text-sand-100/80">
                    Ribuan produk dari berbagai toko terpercaya, dengan gratis ongkir dan
                    diskon spesial tiap hari.
                  </p>
                  <div className="mt-6 flex flex-wrap gap-3">
                    <Link href="#katalog">
                      <Button size="lg">Belanja Sekarang</Button>
                    </Link>
                    <Link href="/register">
                      <Button
                        size="lg"
                        variant="outline"
                        className="border-white/40 text-white hover:bg-white/10 hover:text-white"
                      >
                        Buka Toko
                      </Button>
                    </Link>
                  </div>
                </div>
                <div className="hidden justify-center lg:flex">
                  <ShoppingBag className="h-32 w-32 text-white/20" strokeWidth={1} />
                </div>
              </div>
            </div>
          </section>

          {/* Category strip */}
          <section className="mx-auto max-w-6xl px-5 py-8">
            <div className="grid grid-cols-4 gap-3 sm:grid-cols-8">
              {categories.map((c) => (
                <span
                  key={c.label}
                  className="flex flex-col items-center gap-2 rounded-xl border border-sand-200 bg-white p-3 text-center"
                >
                  <c.icon className="h-6 w-6 text-navy-800" strokeWidth={1.5} />
                  <span className="text-xs font-medium text-ink/70">{c.label}</span>
                </span>
              ))}
            </div>
          </section>

          {/* Flash-sale style featured products */}
          <section className="mx-auto max-w-6xl px-5 py-6">
            <div className="flex items-center justify-between rounded-t-xl bg-coral-500 px-5 py-3">
              <div className="flex items-center gap-2 text-white">
                <Zap className="h-5 w-5" />
                <h2 className="font-display text-lg font-bold">Flash Sale</h2>
              </div>
              <Link href="#katalog" className="text-sm font-semibold text-white hover:underline">
                Lihat semua →
              </Link>
            </div>

            {products.length > 0 ? (
              <div className="grid grid-cols-2 gap-3 rounded-b-xl border border-t-0 border-sand-200 bg-white p-4 sm:grid-cols-3 lg:grid-cols-5">
                {products.slice(0, 10).map((p) => (
                  <ProductCard key={p.id} product={p} />
                ))}
              </div>
            ) : (
              <div className="rounded-b-xl border border-t-0 border-sand-200 bg-white p-8 text-center text-sm text-ink/50">
                Belum ada produk terhubung ke backend, atau backend belum berjalan.
              </div>
            )}
          </section>

          {/* Role explainer */}
          <section className="mx-auto max-w-6xl px-5 py-10">
            <h2 className="font-display text-xl font-bold text-navy-950">Satu akun, banyak peran</h2>
            <div className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {roleCards.map((r) => (
                <div key={r.title} className="rounded-xl border border-sand-200 bg-white p-5">
                  <r.icon className="h-6 w-6 text-navy-800" strokeWidth={1.5} />
                  <p className="mt-2 font-display font-bold text-navy-950">{r.title}</p>
                  <p className="mt-1 text-sm text-ink/60">{r.desc}</p>
                </div>
              ))}
            </div>
          </section>
        </>
      )}

      {/* Full catalog */}
      <div id="katalog" className="mx-auto max-w-6xl px-5 py-8">
        {!isSearching && (
          <div className="mb-6 flex gap-3 overflow-x-auto pb-1">
            {["Semua", "Kelontong", "Fashion", "Elektronik", "Rumah Tangga", "Kecantikan"].map((c) => (
              <span
                key={c}
                className="shrink-0 rounded-full border border-sand-200 bg-white px-4 py-1.5 text-sm text-ink/70"
              >
                {c}
              </span>
            ))}
          </div>
        )}

        <div className="flex items-center justify-between">
          <h2 className="font-display text-xl font-bold text-navy-950">
            {isSearching ? `Hasil untuk "${q}"` : "Semua Produk"}
          </h2>
          <span className="text-sm text-ink/50">{products.length} produk</span>
        </div>

        {products.length === 0 ? (
          <div className="mt-10 rounded-xl border border-dashed border-sand-300 bg-white p-10 text-center">
            <p className="font-display text-lg text-navy-950">Produk tidak ditemukan</p>
            <p className="mt-1 text-sm text-ink/50">
              Coba kata kunci lain, atau pastikan backend berjalan di{" "}
              <code className="font-mono-num">{API_BASE}</code>.
            </p>
          </div>
        ) : (
          <div className="mt-5 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
            {products.map((p) => (
              <ProductCard key={p.id} product={p} />
            ))}
          </div>
        )}
      </div>

      {!isSearching && <ReviewSection />}
    </div>
  );
}
