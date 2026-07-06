"use client";

import { useEffect, useState } from "react";
import { Receipt, ShoppingBag } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { formatIDR } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card } from "@/components/ui/Card";
import type { SpendingReportResponse } from "@/lib/types";

function BuyerReportsContent() {
  const [report, setReport] = useState<SpendingReportResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await api.get<SpendingReportResponse>("/buyer/reports/spending");
        setReport(data);
      } catch (err) {
        setError(err instanceof ApiError ? err.message : "Gagal memuat laporan belanja.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <div className="max-w-2xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Riwayat transaksi</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Laporan belanja</h1>
      <p className="mt-1 text-sm text-ink/60">
        Ringkasan seluruh pesanan yang pernah kamu buat di SEAPEDIA, termasuk diskon, ongkir, dan PPN 12%
        yang sudah dihitung di setiap pesanan.
      </p>

      {loading && <p className="mt-6 text-sm text-ink/50">Memuat laporan...</p>}
      {error && <p className="mt-6 text-sm text-red-600">{error}</p>}

      {report && (
        <div className="mt-6 grid gap-4 sm:grid-cols-2">
          <Card className="flex items-center gap-4 p-6">
            <div className="rounded-xl bg-teal-100 p-3 text-teal-600">
              <ShoppingBag className="h-6 w-6" strokeWidth={1.5} />
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Total pesanan</p>
              <p className="font-display text-2xl font-semibold text-navy-950">{report.total_orders}</p>
            </div>
          </Card>

          <Card className="flex items-center gap-4 p-6">
            <div className="rounded-xl bg-coral-500/10 p-3 text-coral-600">
              <Receipt className="h-6 w-6" strokeWidth={1.5} />
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Total belanja</p>
              <p className="font-display text-2xl font-semibold text-navy-950">
                {formatIDR(report.total_spending)}
              </p>
            </div>
          </Card>
        </div>
      )}

      {report && report.total_orders === 0 && (
        <Card className="mt-4 p-6 text-center text-sm text-ink/50">
          Belum ada transaksi. Yuk mulai belanja dari katalog produk.
        </Card>
      )}
    </div>
  );
}

export default function BuyerReportsPage() {
  return (
    <RequireRole role="buyer">
      <BuyerReportsContent />
    </RequireRole>
  );
}
