"use client";

import { useState } from "react";
import Link from "next/link";
import { CheckCircle2 } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { api, ApiError } from "@/lib/api";
import { Button } from "@/components/ui/Button";
import type { CartResponse } from "@/lib/types";

export function AddToCartButton({ productId }: { productId: string }) {
  const { profile, isAuthenticated } = useAuth();
  const [qty, setQty] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [clearing, setClearing] = useState(false);

  if (!isAuthenticated) {
    return (
      <Link href="/login" className="inline-block">
        <Button variant="primary">Masuk untuk membeli</Button>
      </Link>
    );
  }

  if (profile?.active_role !== "buyer") {
    return (
      <p className="text-sm text-ink/50">
        Beralih ke peran Pembeli dari dashboard untuk menambahkan produk ke keranjang.
      </p>
    );
  }

  const isConflict = error?.includes("single-store") || error?.includes("toko lain");

  async function handleAdd() {
    setError(null);
    setSuccess(false);
    setLoading(true);
    try {
      await api.post<CartResponse>("/buyer/cart/items", { product_id: productId, quantity: qty });
      setSuccess(true);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal menambahkan ke keranjang.");
    } finally {
      setLoading(false);
    }
  }

  async function handleClearAndAdd() {
    setClearing(true);
    try {
      await api.del("/buyer/cart");
      await handleAdd();
    } finally {
      setClearing(false);
    }
  }

  return (
    <div>
      <div className="flex items-center gap-3">
        <input
          type="number"
          min={1}
          value={qty}
          onChange={(e) => setQty(Math.max(1, Number(e.target.value)))}
          className="w-20 rounded-lg border border-sand-300 px-3 py-2.5 text-sm"
        />
        <Button onClick={handleAdd} loading={loading}>
          Tambah ke keranjang
        </Button>
      </div>
      {success && (
        <p className="mt-2 flex items-center gap-1.5 text-sm text-teal-600">
          <CheckCircle2 className="h-4 w-4" /> Ditambahkan ke keranjang
        </p>
      )}
      {error && (
        <div className="mt-2 text-sm text-red-600">
          <p>{error}</p>
          {isConflict && (
            <Button size="sm" variant="outline" className="mt-2" loading={clearing} onClick={handleClearAndAdd}>
              Kosongkan keranjang & tambah produk ini
            </Button>
          )}
        </div>
      )}
    </div>
  );
}
