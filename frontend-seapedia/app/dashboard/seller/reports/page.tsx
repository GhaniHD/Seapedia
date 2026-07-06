"use client";

import { useEffect, useState } from "react";
import { PackageCheck, TrendingUp, Undo2 } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { formatIDR } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card } from "@/components/ui/Card";
import type { IncomeReportResponse } from "@/lib/types";

function SellerReportsContent() {
  const [report, setReport] = useState<IncomeReportResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const data = await api.get<IncomeReportResponse>("/seller/reports/income");
        setReport(data);
      } catch (err) {
        setError(err instanceof ApiError ? err.message : "Gagal memuat laporan pendapatan.");
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <div className="max-w-2xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Toko saya</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Laporan pendapatan</h1>
      <p className="mt-1 text-sm text-ink/60">
        Ringkasan pesanan yang masuk ke tokomu. Pesanan yang dikembalikan atau di-refund akan mengurangi
        pendapatan lewat transaksi pembalik (reversal), bukan dihapus dari riwayat.
      </p>

      {loading && <p className="mt-6 text-sm text-ink/50">Memuat laporan...</p>}
      {error && <p className="mt-6 text-sm text-red-600">{error}</p>}

      {report && (
        <div className="mt-6 grid gap-4 sm:grid-cols-3">
          <Card className="flex items-center gap-4 p-6">
            <div className="rounded-xl bg-navy-800/10 p-3 text-navy-800">
              <PackageCheck className="h-6 w-6" strokeWidth={1.5} />
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Total pesanan</p>
              <p className="font-display text-2xl font-semibold text-navy-950">{report.total_orders}</p>
            </div>
          </Card>

          <Card className="flex items-center gap-4 p-6">
            <div className="rounded-xl bg-teal-100 p-3 text-teal-600">
              <TrendingUp className="h-6 w-6" strokeWidth={1.5} />
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Total pendapatan</p>
              <p className="font-display text-2xl font-semibold text-navy-950">
                {formatIDR(report.total_income)}
              </p>
            </div>
          </Card>

          <Card className="flex items-center gap-4 p-6">
            <div className="rounded-xl bg-red-100 p-3 text-red-600">
              <Undo2 className="h-6 w-6" strokeWidth={1.5} />
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">
                Dibalik (refund/return)
              </p>
              <p className="font-display text-2xl font-semibold text-navy-950">
                {formatIDR(report.total_reversed)}
              </p>
            </div>
          </Card>
        </div>
      )}

      {report && report.total_orders === 0 && (
        <Card className="mt-4 p-6 text-center text-sm text-ink/50">Belum ada pesanan masuk ke tokomu.</Card>
      )}
    </div>
  );
}

export default function SellerReportsPage() {
  return (
    <RequireRole role="seller">
      <SellerReportsContent />
    </RequireRole>
  );
}
