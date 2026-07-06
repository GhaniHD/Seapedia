"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { ShoppingBag, Store, Truck, Compass, type LucideIcon } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { ApiError } from "@/lib/api";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { ShipWheelMark } from "@/components/ui/WaveDivider";
import { ROLE_LABEL } from "@/lib/format";
import type { Role } from "@/lib/types";

const ROLE_ICON: Record<string, LucideIcon> = {
  buyer: ShoppingBag,
  seller: Store,
  driver: Truck,
};

export default function SelectRolePage() {
  const { availableRoles, selectRole } = useAuth();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loadingRole, setLoadingRole] = useState<Role | null>(null);

  async function handleSelect(role: Role) {
    setError(null);
    setLoadingRole(role);
    try {
      await selectRole(role);
      router.push("/dashboard");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal memilih peran.");
      setLoadingRole(null);
    }
  }

  if (availableRoles.length === 0) {
    return (
      <div className="flex min-h-[calc(100vh-64px)] items-center justify-center px-5">
        <Card className="max-w-md p-8 text-center">
          <p className="text-ink/60">
            Tidak ada peran untuk dipilih. Silakan{" "}
            <a href="/login" className="text-coral-600 underline">
              masuk kembali
            </a>
            .
          </p>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex min-h-[calc(100vh-64px)] items-center justify-center bg-sand-100/60 px-5 py-16">
      <Card className="w-full max-w-lg p-8">
        <div className="flex items-center gap-2 text-navy-950">
          <ShipWheelMark className="h-7 w-7 text-teal-500" />
          <span className="font-display text-xl">SEAPEDIA</span>
        </div>
        <h1 className="mt-6 font-display text-2xl font-semibold text-navy-950">
          Pilih peran aktif untuk sesi ini
        </h1>
        <p className="mt-1 text-sm text-ink/60">
          Akunmu memiliki lebih dari satu peran. Semua akses dan halaman privat mengikuti
          peran yang kamu pilih di sini.
        </p>

        <div className="mt-6 grid gap-3">
          {availableRoles.map((role) => {
            const Icon = ROLE_ICON[role] || Compass;
            return (
              <button
                key={role}
                onClick={() => handleSelect(role)}
                disabled={loadingRole !== null}
                className="flex items-center gap-4 rounded-xl border border-sand-200 bg-white p-4 text-left transition-colors hover:border-teal-500 hover:bg-teal-100/30 disabled:opacity-50"
              >
                <Icon className="h-6 w-6 text-navy-800" strokeWidth={1.5} />
                <div className="flex-1">
                  <p className="font-display font-semibold text-navy-950">{ROLE_LABEL[role] || role}</p>
                  <p className="text-xs text-ink/50">Masuk sebagai {ROLE_LABEL[role] || role}</p>
                </div>
                {loadingRole === role && (
                  <span className="h-4 w-4 animate-spin rounded-full border-2 border-teal-500 border-t-transparent" />
                )}
              </button>
            );
          })}
        </div>
        {error && <p className="mt-4 text-sm text-red-600">{error}</p>}
        <Button variant="ghost" className="mt-4 w-full" onClick={() => router.push("/login")}>
          Batal, kembali ke login
        </Button>
      </Card>
    </div>
  );
}
