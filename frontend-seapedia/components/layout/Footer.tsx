import Link from "next/link";
import { ShipWheelMark } from "@/components/ui/WaveDivider";

export function Footer() {
  return (
    <footer className="mt-auto border-t border-sand-200 bg-navy-950 text-sand-100">
      <div className="mx-auto max-w-6xl px-5 py-10">
        <div className="flex flex-col gap-8 md:flex-row md:justify-between">
          <div className="max-w-sm">
            <div className="flex items-center gap-2 text-sand-50">
              <ShipWheelMark className="h-6 w-6 text-teal-400" />
              <span className="font-display text-lg">SEAPEDIA</span>
            </div>
            <p className="mt-3 text-sm text-sand-100/60">
              Satu marketplace untuk banyak toko. SEAPEDIA menghubungkan penjual,
              pembeli, dan kurir dalam satu pengalaman belanja.
            </p>
          </div>
          <div className="grid grid-cols-2 gap-8 text-sm sm:grid-cols-3">
            <div>
              <p className="mb-2 font-semibold text-sand-50">Marketplace</p>
              <ul className="flex flex-col gap-1.5 text-sand-100/60">
                <li><Link href="/products" className="hover:text-teal-400">Katalog produk</Link></li>
                <li><Link href="/#reviews" className="hover:text-teal-400">Ulasan pengguna</Link></li>
              </ul>
            </div>
            <div>
              <p className="mb-2 font-semibold text-sand-50">Akun</p>
              <ul className="flex flex-col gap-1.5 text-sand-100/60">
                <li><Link href="/login" className="hover:text-teal-400">Masuk</Link></li>
                <li><Link href="/register" className="hover:text-teal-400">Daftar</Link></li>
              </ul>
            </div>
          </div>
        </div>
        <div className="mt-8 border-t border-white/10 pt-5 text-xs text-sand-100/40">
          © {new Date().getFullYear()} SEAPEDIA. Dibuat untuk COMPFEST Software Engineering Academy.
        </div>
      </div>
    </footer>
  );
}
