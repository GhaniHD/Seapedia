"use client";

import { Truck } from "lucide-react";
import { RequireRole } from "@/components/features/RequireRole";
import { Card } from "@/components/ui/Card";

function DriverShellContent() {
  return (
    <Card className="max-w-lg p-8">
      <Truck className="h-8 w-8 text-navy-800" strokeWidth={1.5} />
      <h1 className="mt-3 font-display text-2xl font-semibold text-navy-950">
        Dashboard Kurir
      </h1>
      <p className="mt-2 text-sm text-ink/60">
        Fitur temukan job, ambil job, dan riwayat pengiriman akan tersedia di
        Level 5 (Delivery and Driver Workflow). Halaman ini sudah disiapkan
        sebagai placeholder navigasi.
      </p>
    </Card>
  );
}

export default function DriverDashboardPage() {
  return (
    <RequireRole role="driver">
      <DriverShellContent />
    </RequireRole>
  );
}
