"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import { formatIDR, formatDate, DELIVERY_LABEL } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, StatusPill } from "@/components/ui/Card";
import type { OrderResponse } from "@/lib/types";

function BuyerOrderDetailContent() {
  const params = useParams<{ id: string }>();
  const [order, setOrder] = useState<OrderResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await api.get<OrderResponse>(`/buyer/orders/${params.id}`);
        setOrder(data);
      } catch (err) {
        setError(err instanceof ApiError ? err.message : "Gagal memuat detail pesanan.");
      } finally {
        setLoading(false);
      }
    })();
  }, [params.id]);

  if (loading) return <p className="text-sm text-ink/50">Memuat detail pesanan...</p>;
  if (error || !order) return <p className="text-sm text-red-600">{error || "Pesanan tidak ditemukan."}</p>;

  return (
    <div className="max-w-2xl">
      <p className="font-mono-num text-xs text-ink/40">{order.order_no}</p>
      <div className="mt-1 flex items-center justify-between">
        <h1 className="font-display text-2xl font-semibold text-navy-950">{order.store_name}</h1>
        <StatusPill status={order.status} />
      </div>
      <p className="mt-1 text-sm text-ink/50">
        {DELIVERY_LABEL[order.delivery_method]} · dibuat {formatDate(order.created_at)}
      </p>

      <Card className="mt-6 p-5">
        <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Item pesanan</p>
        <ul className="mt-3 flex flex-col gap-2 text-sm">
          {order.items?.map((it, idx) => (
            <li key={idx} className="flex justify-between">
              <span>{it.product_name} × {it.quantity}</span>
              <span>{formatIDR(it.price * it.quantity)}</span>
            </li>
          ))}
        </ul>
      </Card>

      <Card className="mt-4 p-5">
        <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Rincian pembayaran</p>
        <div className="mt-3 flex flex-col gap-2 text-sm text-ink/70">
          <div className="flex justify-between"><span>Subtotal</span><span>{formatIDR(order.subtotal)}</span></div>
          {order.discount_amount > 0 && (
            <div className="flex justify-between"><span>Diskon</span><span>− {formatIDR(order.discount_amount)}</span></div>
          )}
          <div className="flex justify-between"><span>Ongkir</span><span>{formatIDR(order.delivery_fee)}</span></div>
          <div className="flex justify-between"><span>PPN 12%</span><span>{formatIDR(order.tax_amount)}</span></div>
          <div className="my-1 h-px bg-sand-300" />
          <div className="flex justify-between font-display text-base font-semibold text-navy-950">
            <span>Total</span><span>{formatIDR(order.total)}</span>
          </div>
        </div>
      </Card>

      <Card className="mt-4 p-5">
        <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Riwayat status</p>
        <ol className="mt-3 flex flex-col gap-4 border-l-2 border-sand-200 pl-4">
          {order.status_history?.map((h, idx) => (
            <li key={idx} className="relative">
              <span className="absolute -left-[21px] top-1 h-2.5 w-2.5 rounded-full bg-teal-500" />
              <p className="text-sm font-semibold text-navy-950">{h.status}</p>
              {h.note && <p className="text-xs text-ink/50">{h.note}</p>}
              <p className="text-xs text-ink/40">{formatDate(h.created_at)}</p>
            </li>
          ))}
        </ol>
      </Card>
    </div>
  );
}

export default function BuyerOrderDetailPage() {
  return (
    <RequireRole role="buyer">
      <BuyerOrderDetailContent />
    </RequireRole>
  );
}
