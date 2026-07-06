"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import { formatIDR, formatDate } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, Badge } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import type { WalletResponse, WalletTransactionResponse } from "@/lib/types";

const QUICK_AMOUNTS = [50000, 100000, 250000, 500000];

function BuyerWalletContent() {
  const [wallet, setWallet] = useState<WalletResponse>({ balance: 0 });
  const [transactions, setTransactions] = useState<WalletTransactionResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [amount, setAmount] = useState<number>(100000);
  const [topupLoading, setTopupLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function load() {
    setLoading(true);
    try {
      const [w, txs] = await Promise.all([
        api.get<WalletResponse>("/buyer/wallet"),
        api.get<WalletTransactionResponse[]>("/buyer/wallet/transactions"),
      ]);
      setWallet(w);
      setTransactions(txs || []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal memuat dompet.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleTopup(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setTopupLoading(true);
    try {
      await api.post("/buyer/wallet/topup", { amount });
      await load();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal top up.");
    } finally {
      setTopupLoading(false);
    }
  }

  return (
    <div className="max-w-3xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Dompet SEAPEDIA</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Dompet & saldo</h1>

      <Card className="mt-6 bg-navy-950 p-6 text-sand-50">
        <p className="text-xs uppercase tracking-wide text-sand-100/60">Saldo saat ini</p>
        <p className="mt-2 font-display text-4xl font-semibold">
          {loading ? "..." : formatIDR(wallet.balance)}
        </p>
      </Card>

      <Card className="mt-5 p-6">
        <h2 className="font-display text-lg font-semibold text-navy-950">Top up (simulasi)</h2>
        <form onSubmit={handleTopup} className="mt-4">
          <div className="mb-3 flex flex-wrap gap-2">
            {QUICK_AMOUNTS.map((a) => (
              <button
                type="button"
                key={a}
                onClick={() => setAmount(a)}
                className={`rounded-full border px-3 py-1.5 text-xs font-semibold ${
                  amount === a ? "border-teal-500 bg-teal-100 text-teal-600" : "border-sand-300 text-ink/60"
                }`}
              >
                {formatIDR(a)}
              </button>
            ))}
          </div>
          <div className="flex flex-wrap items-end gap-3">
            <div className="w-48">
              <Input
                label="Jumlah top up (Rp)"
                type="number"
                min={1}
                value={amount}
                onChange={(e) => setAmount(Number(e.target.value))}
              />
            </div>
            <Button type="submit" loading={topupLoading}>
              Top up sekarang
            </Button>
          </div>
        </form>
        {error && <p className="mt-3 text-sm text-red-600">{error}</p>}
      </Card>

      <Card className="mt-5 p-6">
        <h2 className="font-display text-lg font-semibold text-navy-950">Riwayat transaksi</h2>
        <div className="mt-4 flex flex-col divide-y divide-sand-200">
          {loading && <p className="py-4 text-sm text-ink/50">Memuat...</p>}
          {!loading && transactions.length === 0 && (
            <p className="py-4 text-sm text-ink/50">Belum ada transaksi dompet.</p>
          )}
          {transactions.map((t) => (
            <div key={t.id} className="flex items-center justify-between py-3">
              <div>
                <p className="text-sm font-medium text-navy-950">{t.description}</p>
                <p className="text-xs text-ink/40">{formatDate(t.created_at)}</p>
              </div>
              <div className="flex items-center gap-2">
                <Badge tone={t.type === "credit" || t.type === "topup" ? "teal" : "coral"}>{t.type}</Badge>
                <span
                  className={`font-mono-num text-sm font-semibold ${
                    t.amount >= 0 ? "text-teal-600" : "text-coral-600"
                  }`}
                >
                  {t.amount >= 0 ? "+" : ""}
                  {formatIDR(t.amount)}
                </span>
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}

export default function BuyerWalletPage() {
  return (
    <RequireRole role="buyer">
      <BuyerWalletContent />
    </RequireRole>
  );
}
