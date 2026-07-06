"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { CheckCircle2 } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { formatIDR, DELIVERY_LABEL } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, Badge } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import type {
  AddressResponse,
  CartResponse,
  CheckoutRequest,
  CheckoutSummaryResponse,
  DeliveryMethod,
  PromoResponse,
  VoucherResponse,
} from "@/lib/types";

// Mirrors backend pkg/utils/order.go — used only for the pre-confirmation
// estimate shown in the UI. The authoritative numbers always come from the
// /buyer/checkout response after confirmation.
const DELIVERY_FEE: Record<DeliveryMethod, number> = {
  instant: 25000,
  next_day: 12000,
  regular: 8000,
};
const TAX_RATE = 0.12;

function BuyerCheckoutContent() {
  const router = useRouter();
  const [cart, setCart] = useState<CartResponse | null>(null);
  const [addresses, setAddresses] = useState<AddressResponse[]>([]);
  const [loading, setLoading] = useState(true);

  const [addressId, setAddressId] = useState("");
  const [deliveryMethod, setDeliveryMethod] = useState<DeliveryMethod>("regular");
  const [discountCode, setDiscountCode] = useState("");
  const [vouchers, setVouchers] = useState<VoucherResponse[]>([]);
  const [promos, setPromos] = useState<PromoResponse[]>([]);

  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<CheckoutSummaryResponse | null>(null);

  async function load() {
    setLoading(true);
    try {
      const [cartData, addrData, voucherData, promoData] = await Promise.all([
        api.get<CartResponse>("/buyer/cart"),
        api.get<AddressResponse[]>("/buyer/addresses"),
        api.get<VoucherResponse[]>("/vouchers").catch(() => []),
        api.get<PromoResponse[]>("/promos").catch(() => []),
      ]);
      setCart(cartData);
      setAddresses(addrData || []);
      const def = addrData?.find((a) => a.is_default) || addrData?.[0];
      if (def) setAddressId(def.id);

      const now = Date.now();
      setVouchers(
        (voucherData || []).filter(
          (v) => new Date(v.expiry_date).getTime() > now && v.usage_count < v.usage_limit
        )
      );
      setPromos((promoData || []).filter((p) => new Date(p.expiry_date).getTime() > now));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal memuat data checkout.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  const estimate = useMemo(() => {
    const subtotal = cart?.subtotal || 0;
    const deliveryFee = DELIVERY_FEE[deliveryMethod];
    const taxAmount = subtotal * TAX_RATE; // discount unknown until server validates the code
    return { subtotal, deliveryFee, taxAmount, total: subtotal + deliveryFee + taxAmount };
  }, [cart, deliveryMethod]);

  async function handleCheckout() {
    if (!addressId) {
      setError("Pilih alamat pengiriman terlebih dahulu.");
      return;
    }
    setError(null);
    setSubmitting(true);
    try {
      const body: CheckoutRequest = {
        address_id: addressId,
        delivery_method: deliveryMethod,
        discount_code: discountCode.trim() || undefined,
      };
      const res = await api.post<CheckoutSummaryResponse>("/buyer/checkout", body);
      setResult(res);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal melakukan checkout.");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) return <p className="text-sm text-ink/50">Memuat checkout...</p>;

  if (result) {
    return (
      <Card className="max-w-lg p-8 text-center">
        <CheckCircle2 className="mx-auto h-10 w-10 text-teal-500" strokeWidth={1.5} />
        <h1 className="mt-3 font-display text-2xl font-semibold text-navy-950">Pesanan berhasil dibuat!</h1>
        <p className="mt-1 font-mono-num text-sm text-ink/50">{result.order_no}</p>

        <div className="mt-6 flex flex-col gap-2 rounded-xl bg-sand-100/60 p-4 text-left text-sm">
          <Row label="Subtotal" value={formatIDR(result.subtotal)} />
          {result.discount_amount > 0 && (
            <Row
              label={`Diskon${result.discount_kind ? ` (${result.discount_kind})` : ""}`}
              value={`− ${formatIDR(result.discount_amount)}`}
            />
          )}
          <Row label="Ongkir" value={formatIDR(result.delivery_fee)} />
          <Row label={`PPN ${(result.tax_rate * 100).toFixed(0)}%`} value={formatIDR(result.tax_amount)} />
          <div className="my-1 h-px bg-sand-300" />
          <Row label="Total dibayar" value={formatIDR(result.total)} bold />
        </div>

        <div className="mt-6 flex justify-center gap-3">
          <Button onClick={() => router.push(`/dashboard/buyer/orders/${result.order_id}`)}>
            Lihat detail pesanan
          </Button>
          <Button variant="ghost" onClick={() => router.push("/products")}>
            Belanja lagi
          </Button>
        </div>
      </Card>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <Card className="max-w-md p-8 text-center">
        <p className="text-ink/60">Keranjangmu kosong. Tambahkan produk dulu sebelum checkout.</p>
        <Button className="mt-4" onClick={() => router.push("/products")}>
          Ke katalog
        </Button>
      </Card>
    );
  }

  return (
    <div className="max-w-3xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Checkout</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Selesaikan pesananmu</h1>

      <div className="mt-6 grid gap-6 lg:grid-cols-[1.3fr_1fr]">
        <div className="flex flex-col gap-5">
          <Card className="p-5">
            <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Belanja dari toko</p>
            <p className="font-display text-lg font-semibold text-navy-950">{cart.store_name}</p>
            <ul className="mt-3 flex flex-col gap-1.5 text-sm text-ink/70">
              {cart.items.map((i) => (
                <li key={i.id} className="flex justify-between">
                  <span>{i.name} × {i.quantity}</span>
                  <span>{formatIDR(i.subtotal)}</span>
                </li>
              ))}
            </ul>
          </Card>

          <Card className="p-5">
            <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-ink/40">Alamat pengiriman</p>
            {addresses.length === 0 ? (
              <p className="text-sm text-red-600">
                Belum ada alamat. Tambahkan alamat di menu Alamat terlebih dahulu.
              </p>
            ) : (
              <div className="flex flex-col gap-2">
                {addresses.map((a) => (
                  <label
                    key={a.id}
                    className={`flex cursor-pointer items-start gap-3 rounded-xl border p-3 text-sm ${
                      addressId === a.id ? "border-teal-500 bg-teal-100/30" : "border-sand-200"
                    }`}
                  >
                    <input
                      type="radio"
                      className="mt-1"
                      checked={addressId === a.id}
                      onChange={() => setAddressId(a.id)}
                    />
                    <span>
                      <span className="font-semibold text-navy-950">{a.label}</span>
                      {a.is_default && <Badge tone="teal">Utama</Badge>}
                      <br />
                      <span className="text-ink/60">{a.detail}</span>
                    </span>
                  </label>
                ))}
              </div>
            )}
          </Card>

          <Card className="p-5">
            <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-ink/40">Metode pengiriman</p>
            <div className="grid grid-cols-3 gap-2">
              {(["instant", "next_day", "regular"] as DeliveryMethod[]).map((m) => (
                <button
                  key={m}
                  onClick={() => setDeliveryMethod(m)}
                  className={`rounded-xl border p-3 text-left text-sm ${
                    deliveryMethod === m ? "border-teal-500 bg-teal-100/30" : "border-sand-200"
                  }`}
                >
                  <p className="font-semibold text-navy-950">{DELIVERY_LABEL[m]}</p>
                  <p className="text-xs text-ink/50">{formatIDR(DELIVERY_FEE[m])}</p>
                </button>
              ))}
            </div>
          </Card>

          <Card className="p-5">
            <Input
              label="Kode voucher / promo (opsional)"
              value={discountCode}
              onChange={(e) => setDiscountCode(e.target.value.toUpperCase())}
              placeholder="cth. SEA10"
            />
            {(vouchers.length > 0 || promos.length > 0) && (
              <div className="mt-3 flex flex-wrap gap-2">
                {vouchers.map((v) => (
                  <button
                    key={v.id}
                    type="button"
                    onClick={() => setDiscountCode(v.code)}
                    className="flex items-center gap-1.5 rounded-full border border-teal-300 bg-teal-100/40 px-3 py-1 text-xs font-semibold text-teal-600 hover:bg-teal-100"
                  >
                    <Badge tone="teal">Voucher</Badge>
                    {v.code}
                  </button>
                ))}
                {promos.map((p) => (
                  <button
                    key={p.id}
                    type="button"
                    onClick={() => setDiscountCode(p.code)}
                    className="flex items-center gap-1.5 rounded-full border border-coral-500/30 bg-coral-500/10 px-3 py-1 text-xs font-semibold text-coral-600 hover:bg-coral-500/20"
                  >
                    <Badge tone="coral">Promo</Badge>
                    {p.code}
                  </button>
                ))}
              </div>
            )}
            <p className="mt-2 text-xs text-ink/40">Klik salah satu kode untuk menerapkannya.</p>
          </Card>
        </div>

        <Card className="h-fit p-5">
          <p className="text-xs font-semibold uppercase tracking-wide text-ink/40">Ringkasan pembayaran</p>
          <div className="mt-3 flex flex-col gap-2 text-sm">
            <Row label="Subtotal" value={formatIDR(estimate.subtotal)} />
            <Row label="Ongkir" value={formatIDR(estimate.deliveryFee)} />
            <Row label="PPN 12%" value={formatIDR(estimate.taxAmount)} />
            <p className="text-xs text-ink/40">
              *Diskon (jika kode valid) dihitung final saat konfirmasi.
            </p>
            <div className="my-1 h-px bg-sand-300" />
            <Row label="Estimasi total" value={formatIDR(estimate.total)} bold />
          </div>
          {error && <p className="mt-3 text-sm text-red-600">{error}</p>}
          <Button className="mt-4 w-full" loading={submitting} onClick={handleCheckout}>
            Konfirmasi & bayar dengan dompet
          </Button>
        </Card>
      </div>
    </div>
  );
}

function Row({ label, value, bold }: { label: string; value: string; bold?: boolean }) {
  return (
    <div className={`flex justify-between ${bold ? "font-display text-base font-semibold text-navy-950" : "text-ink/70"}`}>
      <span>{label}</span>
      <span>{value}</span>
    </div>
  );
}

export default function BuyerCheckoutPage() {
  return (
    <RequireRole role="buyer">
      <BuyerCheckoutContent />
    </RequireRole>
  );
}
