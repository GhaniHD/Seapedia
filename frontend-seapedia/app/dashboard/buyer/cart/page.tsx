"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, ApiError } from "@/lib/api";
import { formatIDR } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, Badge } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import type { CartResponse } from "@/lib/types";

function BuyerCartContent() {
  const [cart, setCart] = useState<CartResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [busyId, setBusyId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const data = await api.get<CartResponse>("/buyer/cart");
      setCart(data);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal memuat keranjang.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function updateQty(productId: string, quantity: number) {
    if (quantity < 1) return;
    setBusyId(productId);
    try {
      const data = await api.put<CartResponse>(`/buyer/cart/items/${productId}`, { quantity });
      setCart(data);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : "Gagal memperbarui jumlah.");
    } finally {
      setBusyId(null);
    }
  }

  async function removeItem(productId: string) {
    setBusyId(productId);
    try {
      const data = await api.del<CartResponse>(`/buyer/cart/items/${productId}`);
      setCart(data);
    } catch (err) {
      alert(err instanceof ApiError ? err.message : "Gagal menghapus item.");
    } finally {
      setBusyId(null);
    }
  }

  async function clearCart() {
    if (!confirm("Kosongkan seluruh keranjang?")) return;
    try {
      await api.del("/buyer/cart");
      await load();
    } catch (err) {
      alert(err instanceof ApiError ? err.message : "Gagal mengosongkan keranjang.");
    }
  }

  if (loading) return <p className="text-sm text-ink/50">Memuat keranjang...</p>;

  const isEmpty = !cart || cart.items.length === 0;

  return (
    <div className="max-w-3xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Keranjang belanja</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Keranjang saya</h1>
      <p className="mt-1 flex items-center gap-1.5 text-sm text-ink/60">
        <Badge tone="sand">Single-store checkout</Badge>
        Satu keranjang hanya bisa berisi produk dari satu toko.
      </p>

      {error && <p className="mt-4 text-sm text-red-600">{error}</p>}

      {isEmpty ? (
        <Card className="mt-6 p-10 text-center">
          <p className="font-display text-lg text-navy-950">Keranjangmu masih kosong</p>
          <Link href="/products">
            <Button className="mt-4">Jelajahi katalog</Button>
          </Link>
        </Card>
      ) : (
        <>
          <Card className="mt-6 p-5">
            <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Belanja dari toko</p>
            <p className="font-display text-lg font-semibold text-navy-950">{cart.store_name}</p>
          </Card>

          <div className="mt-4 flex flex-col gap-3">
            {cart.items.map((item) => (
              <Card key={item.id} className="flex items-center justify-between gap-4 p-4">
                <div className="flex-1">
                  <p className="font-medium text-navy-950">{item.name}</p>
                  <p className="text-sm text-ink/50">{formatIDR(item.price)} / item</p>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => updateQty(item.product_id, item.quantity - 1)}
                    disabled={busyId === item.product_id}
                    className="h-8 w-8 rounded-lg border border-sand-300 text-ink/60 hover:bg-sand-100"
                  >
                    −
                  </button>
                  <span className="w-6 text-center text-sm font-medium">{item.quantity}</span>
                  <button
                    onClick={() => updateQty(item.product_id, item.quantity + 1)}
                    disabled={busyId === item.product_id}
                    className="h-8 w-8 rounded-lg border border-sand-300 text-ink/60 hover:bg-sand-100"
                  >
                    +
                  </button>
                </div>
                <p className="w-28 text-right font-semibold text-navy-800">{formatIDR(item.subtotal)}</p>
                <button
                  onClick={() => removeItem(item.product_id)}
                  disabled={busyId === item.product_id}
                  className="text-sm text-red-600 hover:underline"
                >
                  Hapus
                </button>
              </Card>
            ))}
          </div>

          <Card className="mt-4 flex items-center justify-between p-5">
            <div>
              <p className="text-xs uppercase tracking-wide text-ink/40">Subtotal</p>
              <p className="font-display text-xl font-semibold text-navy-950">{formatIDR(cart.subtotal)}</p>
            </div>
            <div className="flex gap-3">
              <Button variant="ghost" onClick={clearCart}>
                Kosongkan
              </Button>
              <Link href="/dashboard/buyer/checkout">
                <Button>Lanjut ke checkout →</Button>
              </Link>
            </div>
          </Card>
        </>
      )}
    </div>
  );
}

export default function BuyerCartPage() {
  return (
    <RequireRole role="buyer">
      <BuyerCartContent />
    </RequireRole>
  );
}
