"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, ApiError } from "@/lib/api";
import { formatIDR, formatDate, DELIVERY_LABEL } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, StatusPill } from "@/components/ui/Card";
import type { OrderResponse } from "@/lib/types";

function BuyerOrdersContent() {
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await api.get<OrderResponse[]>("/buyer/orders");
        setOrders(data || []);
      } catch (err) {
        setError(err instanceof ApiError ? err.message : "Gagal memuat pesanan.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <div className="max-w-3xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Riwayat transaksi</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Pesanan saya</h1>

      {loading && <p className="mt-6 text-sm text-ink/50">Memuat pesanan...</p>}
      {error && <p className="mt-6 text-sm text-red-600">{error}</p>}

      {!loading && orders.length === 0 && (
        <Card className="mt-6 p-8 text-center text-sm text-ink/50">Belum ada pesanan.</Card>
      )}

      <div className="mt-6 flex flex-col gap-3">
        {orders.map((o) => (
          <Link key={o.id} href={`/dashboard/buyer/orders/${o.id}`}>
            <Card className="flex items-center justify-between p-5 transition-shadow hover:shadow-md">
              <div>
                <p className="font-mono-num text-xs text-ink/40">{o.order_no}</p>
                <p className="font-semibold text-navy-950">{o.store_name}</p>
                <p className="text-xs text-ink/50">
                  {DELIVERY_LABEL[o.delivery_method]} · {formatDate(o.created_at)}
                </p>
              </div>
              <div className="text-right">
                <StatusPill status={o.status} />
                <p className="mt-2 font-semibold text-navy-800">{formatIDR(o.total)}</p>
              </div>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  );
}

export default function BuyerOrdersPage() {
  return (
    <RequireRole role="buyer">
      <BuyerOrdersContent />
    </RequireRole>
  );
}
