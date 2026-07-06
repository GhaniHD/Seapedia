export function formatIDR(amount: number): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(amount || 0);
}

export function formatDate(iso: string): string {
  if (!iso) return "-";
  const d = new Date(iso);
  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(d);
}

export const DELIVERY_LABEL: Record<string, string> = {
  instant: "Instant",
  next_day: "Next Day",
  regular: "Regular",
};

export const ROLE_LABEL: Record<string, string> = {
  admin: "Admin",
  seller: "Penjual",
  buyer: "Pembeli",
  driver: "Kurir",
};
