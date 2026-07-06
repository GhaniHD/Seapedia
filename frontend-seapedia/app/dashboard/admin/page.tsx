"use client";

import { LayoutDashboard } from "lucide-react";
import { RequireRole } from "@/components/features/RequireRole";
import { Card } from "@/components/ui/Card";

function AdminShellContent() {
  return (
    <Card className="max-w-lg p-8">
      <LayoutDashboard className="h-8 w-8 text-navy-800" strokeWidth={1.5} />
      <h1 className="mt-3 font-display text-2xl font-semibold text-navy-950">
        Dashboard Admin
      </h1>
      <p className="mt-2 text-sm text-ink/60">
        Monitoring marketplace penuh, manajemen voucher/promo, dan overdue
        handling akan tersedia di Level 6 (Admin Monitoring and Overdue
        Handling). Halaman ini sudah disiapkan sebagai placeholder navigasi.
      </p>
    </Card>
  );
}

export default function AdminDashboardPage() {
  return (
    <RequireRole role="admin">
      <AdminShellContent />
    </RequireRole>
  );
}
