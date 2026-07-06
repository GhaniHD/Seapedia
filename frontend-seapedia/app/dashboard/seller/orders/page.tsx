"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import { formatIDR, formatDate, DELIVERY_LABEL } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, StatusPill } from "@/components/ui/Card";
import type { OrderResponse } from "@/lib/types";

function SellerOrdersContent() {
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await api.get<OrderResponse[]>("/seller/orders");
        setOrders(data || []);
      } catch (err) {
        setError(err instanceof ApiError ? err.message : "Gagal memuat pesanan masuk.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <div className="max-w-3xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Toko saya</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Pesanan masuk</h1>
      <p className="mt-1 text-sm text-ink/60">
        Memproses pesanan (Sedang Dikemas → Menunggu Pengirim) tersedia mulai Level 4.
      </p>

      {loading && <p className="mt-6 text-sm text-ink/50">Memuat...</p>}
      {error && <p className="mt-6 text-sm text-red-600">{error}</p>}
      {!loading && orders.length === 0 && (
        <Card className="mt-6 p-8 text-center text-sm text-ink/50">Belum ada pesanan masuk.</Card>
      )}

      <div className="mt-6 flex flex-col gap-3">
        {orders.map((o) => (
          <Card key={o.id} className="flex items-center justify-between p-5">
            <div>
              <p className="font-mono-num text-xs text-ink/40">{o.order_no}</p>
              <p className="font-semibold text-navy-950">{o.buyer_name}</p>
              <p className="text-xs text-ink/50">
                {DELIVERY_LABEL[o.delivery_method]} · {formatDate(o.created_at)}
              </p>
            </div>
            <div className="text-right">
              <StatusPill status={o.status} />
              <p className="mt-2 font-semibold text-navy-800">{formatIDR(o.total)}</p>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}

export default function SellerOrdersPage() {
  return (
    <RequireRole role="seller">
      <SellerOrdersContent />
    </RequireRole>
  );
}
