"use client";

import { useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { api, ApiError } from "@/lib/api";
import { formatIDR, ROLE_LABEL } from "@/lib/format";
import { Card, Badge } from "@/components/ui/Card";
import { Button } from "@/components/ui/Button";
import type { AddRoleRequest, Role } from "@/lib/types";

const ADDABLE_ROLES: Role[] = ["buyer", "seller", "driver"];

export default function DashboardHomePage() {
  const { profile, refreshProfile, selectRole } = useAuth();
  const [addingRole, setAddingRole] = useState<Role | null>(null);
  const [error, setError] = useState<string | null>(null);

  if (!profile) return null;

  const missingRoles = ADDABLE_ROLES.filter((r) => !profile.roles.includes(r));

  async function handleAddRole(role: Role) {
    setError(null);
    setAddingRole(role);
    try {
      const body: AddRoleRequest = { role };
      await api.post("/roles", body);
      await refreshProfile();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal menambah peran.");
    } finally {
      setAddingRole(null);
    }
  }

  async function handleSwitchRole(role: Role) {
    setError(null);
    try {
      await selectRole(role);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal berganti peran.");
    }
  }

  return (
    <div className="max-w-3xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Ringkasan akun</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">
        Halo, {profile.name.split(" ")[0]}
      </h1>
      <p className="mt-1 text-sm text-ink/60">{profile.email}</p>

      <Card className="mt-6 p-6">
        <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Peran aktif</p>
        <div className="mt-2 flex items-center gap-2">
          <Badge tone="teal">{ROLE_LABEL[profile.active_role] || profile.active_role}</Badge>
        </div>

        <p className="mt-5 text-xs font-semibold uppercase tracking-wide text-ink/40">Peran yang dimiliki</p>
        <div className="mt-2 flex flex-wrap gap-2">
          {profile.roles.map((r) => (
            <button
              key={r}
              onClick={() => (r !== profile.active_role ? handleSwitchRole(r) : undefined)}
              className={`rounded-full border px-3 py-1.5 text-xs font-semibold transition-colors ${
                r === profile.active_role
                  ? "border-teal-500 bg-teal-100 text-teal-600"
                  : "border-sand-300 text-ink/60 hover:border-navy-800 hover:text-navy-800"
              }`}
            >
              {ROLE_LABEL[r] || r}
              {r !== profile.active_role && " · pindah ke sini"}
            </button>
          ))}
        </div>

        {profile.active_role !== "admin" && missingRoles.length > 0 && (
          <>
            <p className="mt-5 text-xs font-semibold uppercase tracking-wide text-ink/40">
              Tambah peran lain
            </p>
            <div className="mt-2 flex flex-wrap gap-2">
              {missingRoles.map((r) => (
                <Button
                  key={r}
                  size="sm"
                  variant="outline"
                  loading={addingRole === r}
                  onClick={() => handleAddRole(r)}
                >
                  + Jadi {ROLE_LABEL[r]}
                </Button>
              ))}
            </div>
          </>
        )}
        {error && <p className="mt-3 text-sm text-red-600">{error}</p>}
      </Card>

      <div className="mt-6 grid gap-4 sm:grid-cols-3">
        <Card className="p-5">
          <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Saldo dompet</p>
          <p className="mt-2 font-display text-2xl font-semibold text-navy-950">
            {profile.wallet_balance != null ? formatIDR(profile.wallet_balance) : "—"}
          </p>
        </Card>
        <Card className="p-5">
          <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Pendapatan toko</p>
          <p className="mt-2 font-display text-2xl font-semibold text-navy-950">
            {profile.store_income != null ? formatIDR(profile.store_income) : "—"}
          </p>
        </Card>
        <Card className="p-5">
          <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Pendapatan kurir</p>
          <p className="mt-2 font-display text-2xl font-semibold text-navy-950">
            {profile.driver_earning != null ? formatIDR(profile.driver_earning) : "—"}
          </p>
        </Card>
      </div>
    </div>
  );
}
