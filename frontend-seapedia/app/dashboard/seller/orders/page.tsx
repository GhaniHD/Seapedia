"use client";

import { useEffect, useState } from "react";
import { ChevronDown, ChevronUp, PackageCheck } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { formatIDR, formatDate, DELIVERY_LABEL } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, StatusPill } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import type { OrderResponse } from "@/lib/types";

function SellerOrdersContent() {
  const [orders, setOrders] = useState<OrderResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [processingId, setProcessingId] = useState<string | null>(null);
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);

  async function load() {
    try {
      const data = await api.get<OrderResponse[]>("/seller/orders");
      setOrders(data || []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal memuat pesanan masuk.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleProcess(orderId: string) {
    setActionError(null);
    setProcessingId(orderId);
    try {
      await api.post(`/seller/orders/${orderId}/process`);
      await load();
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : "Gagal memproses pesanan.");
    } finally {
      setProcessingId(null);
    }
  }

  return (
    <div className="max-w-3xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Toko saya</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Pesanan masuk</h1>
      <p className="mt-1 text-sm text-ink/60">
        Proses pesanan yang masih <span className="font-semibold">Sedang Dikemas</span> agar tersedia
        untuk diambil Driver. Pesanan yang sudah diproses akan berpindah ke status{" "}
        <span className="font-semibold">Menunggu Pengirim</span>.
      </p>

      {loading && <p className="mt-6 text-sm text-ink/50">Memuat...</p>}
      {error && <p className="mt-6 text-sm text-red-600">{error}</p>}
      {actionError && <p className="mt-4 text-sm text-red-600">{actionError}</p>}
      {!loading && orders.length === 0 && (
        <Card className="mt-6 p-8 text-center text-sm text-ink/50">Belum ada pesanan masuk.</Card>
      )}

      <div className="mt-6 flex flex-col gap-3">
        {orders.map((o) => {
          const expanded = expandedId === o.id;
          const canProcess = o.status === "Sedang Dikemas";
          return (
            <Card key={o.id} className="p-5">
              <div className="flex items-center justify-between">
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
              </div>

              <div className="mt-4 flex items-center justify-between border-t border-sand-200 pt-3">
                <button
                  className="flex items-center gap-1 text-xs font-semibold text-ink/50 hover:text-navy-800"
                  onClick={() => setExpandedId(expanded ? null : o.id)}
                >
                  {expanded ? <ChevronUp className="h-3.5 w-3.5" /> : <ChevronDown className="h-3.5 w-3.5" />}
                  Riwayat status
                </button>

                {canProcess ? (
                  <Button
                    size="sm"
                    loading={processingId === o.id}
                    onClick={() => handleProcess(o.id)}
                  >
                    <PackageCheck className="h-4 w-4" />
                    Proses pesanan
                  </Button>
                ) : (
                  <span className="text-xs text-ink/40">
                    {o.status === "Menunggu Pengirim"
                      ? "Menunggu driver mengambil pesanan"
                      : "Sudah diproses"}
                  </span>
                )}
              </div>

              {expanded && (
                <ol className="mt-4 flex flex-col gap-3 border-l-2 border-sand-200 pl-4">
                  {o.status_history?.map((h, idx) => (
                    <li key={idx} className="relative">
                      <span className="absolute -left-[21px] top-1 h-2.5 w-2.5 rounded-full bg-teal-500" />
                      <p className="text-sm font-semibold text-navy-950">{h.status}</p>
                      {h.note && <p className="text-xs text-ink/50">{h.note}</p>}
                      <p className="text-xs text-ink/40">{formatDate(h.created_at)}</p>
                    </li>
                  ))}
                  {!o.status_history?.length && (
                    <li className="text-xs text-ink/40">Belum ada riwayat status.</li>
                  )}
                </ol>
              )}
            </Card>
          );
        })}
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
