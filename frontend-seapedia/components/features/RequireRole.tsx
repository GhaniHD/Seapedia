"use client";

import { useAuth } from "@/context/AuthContext";
import { ROLE_LABEL } from "@/lib/format";
import { Card } from "@/components/ui/Card";
import type { Role } from "@/lib/types";

export function RequireRole({ role, children }: { role: Role; children: React.ReactNode }) {
  const { profile } = useAuth();

  if (!profile) return null;

  if (profile.active_role !== role) {
    return (
      <Card className="max-w-lg p-6">
        <p className="font-display text-lg font-semibold text-navy-950">
          Halaman ini khusus peran {ROLE_LABEL[role] || role}
        </p>
        <p className="mt-2 text-sm text-ink/60">
          Peran aktifmu saat ini adalah {ROLE_LABEL[profile.active_role] || profile.active_role}.
          Ganti peran aktif dari halaman Ringkasan Profil untuk mengakses ini.
        </p>
      </Card>
    );
  }

  return <>{children}</>;
}
