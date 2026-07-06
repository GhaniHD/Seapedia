import Link from "next/link";
import { notFound } from "next/navigation";
import { ShoppingBag } from "lucide-react";
import { formatIDR } from "@/lib/format";
import type { ProductResponse, StoreResponse } from "@/lib/types";
import { AddToCartButton } from "@/components/features/AddToCartButton";

const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function getProduct(id: string): Promise<ProductResponse | null> {
  try {
    const res = await fetch(`${API_BASE}/products/${id}`, { cache: "no-store" });
    if (!res.ok) return null;
    const json = await res.json();
    return json?.data ?? null;
  } catch {
    return null;
  }
}

async function getStore(id: string): Promise<StoreResponse | null> {
  try {
    const res = await fetch(`${API_BASE}/stores/${id}`, { cache: "no-store" });
    if (!res.ok) return null;
    const json = await res.json();
    return json?.data ?? null;
  } catch {
    return null;
  }
}

export default async function ProductDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  const product = await getProduct(id);
  if (!product) notFound();

  const store = await getStore(product.store_id);

  return (
    <div className="mx-auto max-w-5xl px-5 py-12">
      <Link href="/products" className="text-sm text-navy-800 hover:text-coral-600">
        ← Kembali ke katalog
      </Link>

      <div className="mt-6 grid gap-10 lg:grid-cols-2">
        <div className="flex aspect-square items-center justify-center rounded-xl border border-sand-200 bg-sand-100 text-ink/20">
          <ShoppingBag className="h-20 w-20" strokeWidth={1} />
        </div>

        <div>
          <h1 className="mt-1 font-display text-2xl font-bold text-navy-950">{product.name}</h1>
          <p className="mt-2 font-display text-3xl font-extrabold text-coral-500">{formatIDR(product.price)}</p>
          <p className="mt-1 text-sm text-ink/50">Stok tersedia: {product.stock}</p>

          <div className="mt-5">
            <AddToCartButton productId={product.id} />
          </div>

          <p className="mt-6 whitespace-pre-wrap text-ink/70">
            {product.description || "Belum ada deskripsi untuk produk ini."}
          </p>

          <div className="mt-8 rounded-2xl border border-sand-200 bg-white p-5">
            <p className="text-xs font-semibold uppercase tracking-wide text-teal-500">Toko</p>
            <p className="mt-1 font-display text-lg font-semibold text-navy-950">
              {store?.name || product.store_name || "Toko SEAPEDIA"}
            </p>
            {store?.description && (
              <p className="mt-1 text-sm text-ink/60">{store.description}</p>
            )}
            <Link
              href="/login"
              className="mt-3 inline-block text-sm font-semibold text-coral-600 hover:underline"
            >
              Masuk untuk membeli →
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
