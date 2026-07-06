import { HTMLAttributes } from "react";

export function Card({ className = "", children, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={`rounded-2xl border border-sand-200 bg-white shadow-[0_1px_2px_rgba(7,26,51,0.06)] ${className}`}
      {...props}
    >
      {children}
    </div>
  );
}

export function Badge({
  children,
  tone = "navy",
}: {
  children: React.ReactNode;
  tone?: "navy" | "teal" | "coral" | "sand";
}) {
  const tones: Record<string, string> = {
    navy: "bg-navy-800/10 text-navy-800",
    teal: "bg-teal-100 text-teal-500",
    coral: "bg-coral-500/10 text-coral-600",
    sand: "bg-sand-200 text-ink/70",
  };
  return (
    <span className={`inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold ${tones[tone]}`}>
      {children}
    </span>
  );
}

const STATUS_TONE: Record<string, string> = {
  "Sedang Dikemas": "bg-amber-100 text-amber-700",
  "Menunggu Pengirim": "bg-blue-100 text-blue-700",
  "Sedang Dikirim": "bg-teal-100 text-teal-600",
  "Pesanan Selesai": "bg-emerald-100 text-emerald-700",
  Dikembalikan: "bg-red-100 text-red-700",
};

export function StatusPill({ status }: { status: string }) {
  const tone = STATUS_TONE[status] || "bg-sand-200 text-ink/70";
  return (
    <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold ${tone}`}>
      {status}
    </span>
  );
}
