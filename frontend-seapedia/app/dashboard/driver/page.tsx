"use client";

import { useEffect, useState } from "react";
import { Truck, PackageSearch, MapPin, Wallet, CheckCircle2, History } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { formatIDR, formatDate } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, StatusPill } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";

interface DeliveryJob {
  id: string;
  order_id: string;
  order_no: string;
  store_name: string;
  address: string;
  fee: number;
  status: "available" | "taken" | "completed";
  taken_at?: string | null;
  completed_at?: string | null;
}

interface DriverEarning {
  completed_jobs: number;
  total_earning: number;
}

function DriverDashboardContent() {
  const [availableJobs, setAvailableJobs] = useState<DeliveryJob[]>([]);
  const [myJobs, setMyJobs] = useState<DeliveryJob[]>([]);
  const [earnings, setEarnings] = useState<DriverEarning | null>(null);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionId, setActionId] = useState<string | null>(null);

  async function loadAll() {
    setError(null);
    try {
      const [available, mine, earning] = await Promise.all([
        api.get<DeliveryJob[]>("/driver/jobs"),
        api.get<DeliveryJob[]>("/driver/my-jobs"),
        api.get<DriverEarning>("/driver/earnings"),
      ]);
      setAvailableJobs(available || []);
      setMyJobs(mine || []);
      setEarnings(earning || { completed_jobs: 0, total_earning: 0 });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal memuat data pengiriman.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadAll();
  }, []);

  async function handleTakeJob(jobId: string) {
    setActionError(null);
    setActionId(jobId);
    try {
      await api.post(`/driver/jobs/${jobId}/take`);
      await loadAll();
    } catch (err) {
      setActionError(
        err instanceof ApiError ? err.message : "Gagal mengambil job. Mungkin sudah diambil driver lain."
      );
      await loadAll();
    } finally {
      setActionId(null);
    }
  }

  async function handleCompleteJob(jobId: string) {
    setActionError(null);
    setActionId(jobId);
    try {
      await api.post(`/driver/jobs/${jobId}/complete`);
      await loadAll();
    } catch (err) {
      setActionError(err instanceof ApiError ? err.message : "Gagal mengonfirmasi job selesai.");
    } finally {
      setActionId(null);
    }
  }

  const activeJobs = myJobs.filter((j) => j.status === "taken");
  const completedJobs = myJobs
    .filter((j) => j.status === "completed")
    .sort((a, b) => new Date(b.completed_at || 0).getTime() - new Date(a.completed_at || 0).getTime());

  return (
    <div className="max-w-3xl">
      <div className="flex items-center gap-2 text-teal-500">
        <Truck className="h-5 w-5" />
        <p className="text-sm font-semibold uppercase tracking-wide">Dashboard Kurir</p>
      </div>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Job pengiriman</h1>
      <p className="mt-1 text-sm text-ink/60">
        Ambil pesanan yang statusnya <span className="font-semibold">Menunggu Pengirim</span>, antar, lalu
        konfirmasi selesai untuk mendapatkan penghasilan.
      </p>

      {loading && <p className="mt-6 text-sm text-ink/50">Memuat...</p>}
      {error && <p className="mt-6 text-sm text-red-600">{error}</p>}
      {actionError && <p className="mt-4 text-sm text-red-600">{actionError}</p>}

      {!loading && (
        <>
          {/* Ringkasan penghasilan */}
          <Card className="mt-6 flex items-center justify-between p-5">
            <div className="flex items-center gap-3">
              <div className="rounded-full bg-teal-100 p-2.5 text-teal-600">
                <Wallet className="h-5 w-5" />
              </div>
              <div>
                <p className="text-xs text-ink/50">Total penghasilan</p>
                <p className="font-display text-xl font-semibold text-navy-950">
                  {formatIDR(earnings?.total_earning || 0)}
                </p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-xs text-ink/50">Job selesai</p>
              <p className="font-semibold text-navy-800">{earnings?.completed_jobs || 0}</p>
            </div>
          </Card>

          {/* Job aktif */}
          <section className="mt-8">
            <h2 className="font-display text-lg font-semibold text-navy-950">Job aktif</h2>
            <p className="text-xs text-ink/50">Pesanan yang sedang kamu antar sekarang.</p>

            {activeJobs.length === 0 && (
              <Card className="mt-3 p-6 text-center text-sm text-ink/50">
                Belum ada job aktif. Ambil job dari daftar di bawah.
              </Card>
            )}

            <div className="mt-3 flex flex-col gap-3">
              {activeJobs.map((j) => (
                <Card key={j.id} className="p-5">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="font-mono-num text-xs text-ink/40">{j.order_no}</p>
                      <p className="font-semibold text-navy-950">{j.store_name}</p>
                      <p className="mt-1 flex items-start gap-1 text-xs text-ink/50">
                        <MapPin className="mt-0.5 h-3.5 w-3.5 shrink-0" />
                        {j.address}
                      </p>
                    </div>
                    <div className="text-right">
                      <StatusPill status="Sedang Dikirim" />
                      <p className="mt-2 font-semibold text-navy-800">{formatIDR(j.fee)}</p>
                    </div>
                  </div>
                  <div className="mt-4 flex justify-end border-t border-sand-200 pt-3">
                    <Button
                      size="sm"
                      loading={actionId === j.id}
                      onClick={() => handleCompleteJob(j.id)}
                    >
                      <CheckCircle2 className="h-4 w-4" />
                      Konfirmasi selesai
                    </Button>
                  </div>
                </Card>
              ))}
            </div>
          </section>

          {/* Cari job tersedia */}
          <section className="mt-8">
            <div className="flex items-center gap-2">
              <PackageSearch className="h-4 w-4 text-navy-800" />
              <h2 className="font-display text-lg font-semibold text-navy-950">Job tersedia</h2>
            </div>
            <p className="text-xs text-ink/50">
              Pesanan yang sudah diproses seller dan siap diantar. Ambil sebelum driver lain mengambilnya.
            </p>

            {availableJobs.length === 0 && (
              <Card className="mt-3 p-6 text-center text-sm text-ink/50">
                Belum ada job yang tersedia saat ini.
              </Card>
            )}

            <div className="mt-3 flex flex-col gap-3">
              {availableJobs.map((j) => (
                <Card key={j.id} className="p-5">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="font-mono-num text-xs text-ink/40">{j.order_no}</p>
                      <p className="font-semibold text-navy-950">{j.store_name}</p>
                      <p className="mt-1 flex items-start gap-1 text-xs text-ink/50">
                        <MapPin className="mt-0.5 h-3.5 w-3.5 shrink-0" />
                        {j.address}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-navy-800">{formatIDR(j.fee)}</p>
                    </div>
                  </div>
                  <div className="mt-4 flex justify-end border-t border-sand-200 pt-3">
                    <Button
                      size="sm"
                      variant="secondary"
                      loading={actionId === j.id}
                      disabled={activeJobs.length > 0}
                      onClick={() => handleTakeJob(j.id)}
                    >
                      Ambil job
                    </Button>
                  </div>
                </Card>
              ))}
            </div>
            {activeJobs.length > 0 && availableJobs.length > 0 && (
              <p className="mt-3 text-xs text-ink/40">
                Selesaikan job aktif kamu terlebih dahulu sebelum mengambil job baru.
              </p>
            )}
          </section>

          {/* Riwayat job */}
          <section className="mt-8">
            <div className="flex items-center gap-2">
              <History className="h-4 w-4 text-navy-800" />
              <h2 className="font-display text-lg font-semibold text-navy-950">Riwayat job</h2>
            </div>

            {completedJobs.length === 0 && (
              <Card className="mt-3 p-6 text-center text-sm text-ink/50">Belum ada job yang selesai.</Card>
            )}

            <div className="mt-3 flex flex-col gap-2">
              {completedJobs.map((j) => (
                <Card key={j.id} className="flex items-center justify-between p-4">
                  <div>
                    <p className="font-mono-num text-xs text-ink/40">{j.order_no}</p>
                    <p className="text-sm font-semibold text-navy-950">{j.store_name}</p>
                    <p className="text-xs text-ink/40">{j.completed_at ? formatDate(j.completed_at) : "-"}</p>
                  </div>
                  <p className="font-semibold text-teal-600">+{formatIDR(j.fee)}</p>
                </Card>
              ))}
            </div>
          </section>
        </>
      )}
    </div>
  );
}

export default function DriverDashboardPage() {
  return (
    <RequireRole role="driver">
      <DriverDashboardContent />
    </RequireRole>
  );
}