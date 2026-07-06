import { ProductCard } from "@/components/features/ProductCard";
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

export default async function ProductsPage({
  searchParams,
}: {
  searchParams: Promise<{ q?: string }>;
}) {
  const { q } = await searchParams;
  const allProducts = await getProducts();

  const products = q
    ? allProducts.filter(
        (p) =>
          p.name.toLowerCase().includes(q.toLowerCase()) ||
          p.description?.toLowerCase().includes(q.toLowerCase()) ||
          p.store_name?.toLowerCase().includes(q.toLowerCase())
      )
    : allProducts;

  return (
    <div className="mx-auto max-w-6xl px-5 py-8">
      {/* Category strip */}
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

      <div className="flex items-center justify-between">
        <h1 className="font-display text-xl font-bold text-navy-950">
          {q ? `Hasil untuk "${q}"` : "Semua Produk"}
        </h1>
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
  );
}
