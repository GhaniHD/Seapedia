import Link from "next/link";
import { ShoppingBag, Star } from "lucide-react";
import { formatIDR } from "@/lib/format";
import type { ProductResponse } from "@/lib/types";

// Deterministic pseudo-random display values (rating/sold) derived from the
// product id, purely cosmetic to match a marketplace card layout. Real
// rating/sold data isn't part of the Level 1-3 backend yet.
function hashToRange(id: string, min: number, max: number) {
  let hash = 0;
  for (let i = 0; i < id.length; i++) hash = (hash * 31 + id.charCodeAt(i)) >>> 0;
  return min + (hash % (max - min + 1));
}

export function ProductCard({ product }: { product: ProductResponse }) {
  const rating = (hashToRange(product.id, 38, 50) / 10).toFixed(1);
  const sold = hashToRange(product.id, 3, 999);

  return (
    <Link href={`/products/${product.id}`} className="group block">
      <div className="overflow-hidden rounded-lg border border-sand-200 bg-white transition-shadow hover:shadow-lg">
        <div className="relative flex aspect-square items-center justify-center bg-sand-100 text-ink/20">
          <ShoppingBag className="h-10 w-10" strokeWidth={1.5} />
        </div>
        <div className="p-2.5">
          <h3 className="line-clamp-2 min-h-[2.5rem] text-sm text-ink group-hover:text-navy-800">
            {product.name}
          </h3>
          <p className="mt-1.5 font-display text-base font-bold text-coral-500">
            {formatIDR(product.price)}
          </p>
          <div className="mt-1.5 flex items-center gap-1 text-xs text-ink/50">
            <Star className="h-3 w-3 fill-coral-500 text-coral-500" />
            <span>{rating}</span>
            <span>·</span>
            <span>Terjual {sold}</span>
          </div>
          <p className="mt-1 line-clamp-1 text-xs text-ink/40">
            {product.store_name || "Toko SEAPEDIA"}
          </p>
        </div>
      </div>
    </Link>
  );
}
