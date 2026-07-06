"use client";

import { useEffect } from "react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  ArrowLeft,
  ClipboardList,
  Wallet,
  MapPin,
  ShoppingCart,
  Package,
  Store,
  Boxes,
  Inbox,
  Truck,
  LayoutDashboard,
  type LucideIcon,
} from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { ShipWheelMark } from "@/components/ui/WaveDivider";
import { ROLE_LABEL } from "@/lib/format";
import { Button } from "@/components/ui/Button";

interface NavItem {
  href: string;
  label: string;
  icon: LucideIcon;
}

const NAV_BY_ROLE: Record<string, NavItem[]> = {
  buyer: [
    { href: "/dashboard/buyer/wallet", label: "Dompet", icon: Wallet },
    { href: "/dashboard/buyer/addresses", label: "Alamat", icon: MapPin },
    { href: "/dashboard/buyer/cart", label: "Keranjang", icon: ShoppingCart },
    { href: "/dashboard/buyer/orders", label: "Pesanan saya", icon: Package },
  ],
  seller: [
    { href: "/dashboard/seller/store", label: "Toko saya", icon: Store },
    { href: "/dashboard/seller/products", label: "Produk", icon: Boxes },
    { href: "/dashboard/seller/orders", label: "Pesanan masuk", icon: Inbox },
  ],
  driver: [
    { href: "/dashboard/driver", label: "Job pengiriman", icon: Truck },
  ],
  admin: [
    { href: "/dashboard/admin", label: "Monitoring", icon: LayoutDashboard },
  ],
};

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { profile, loading, isAuthenticated, needRoleSelection, logout } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (loading) return;
    if (needRoleSelection) {
      router.replace("/select-role");
      return;
    }
    if (!isAuthenticated) {
      router.replace("/login");
    }
  }, [loading, isAuthenticated, needRoleSelection, router]);

  if (loading || !isAuthenticated || !profile) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-sand-50">
        <p className="text-sm text-ink/50">Memuat dashboard...</p>
      </div>
    );
  }

  const navItems = NAV_BY_ROLE[profile.active_role] || [];

  return (
    <div className="flex min-h-screen bg-sand-100/50">
      <aside className="hidden w-64 shrink-0 flex-col border-r border-sand-200 bg-navy-950 text-sand-100 md:flex">
        <Link href="/" className="flex items-center gap-2 border-b border-white/10 px-5 py-5 text-sand-50">
          <ShipWheelMark className="h-6 w-6 text-teal-400" />
          <span className="font-display text-lg">SEAPEDIA</span>
        </Link>

        <Link
          href="/products"
          className="flex items-center gap-2 border-b border-white/10 px-5 py-3 text-sm text-sand-100/80 transition-colors hover:bg-white/10 hover:text-white"
        >
          <ArrowLeft className="h-4 w-4" />
          Kembali ke Katalog Produk
        </Link>

        <div className="border-b border-white/10 px-5 py-4">
          <p className="text-sm font-semibold text-sand-50">{profile.name}</p>
          <p className="mt-1 inline-flex rounded-full bg-teal-400/15 px-2 py-0.5 text-xs text-teal-400">
            {ROLE_LABEL[profile.active_role] || profile.active_role}
          </p>
        </div>

        <nav className="flex-1 px-3 py-4">
          <Link
            href="/dashboard"
            className={`mb-1 flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm transition-colors hover:bg-white/10 ${
              pathname === "/dashboard" ? "bg-white/10 text-teal-400" : "text-sand-100/80"
            }`}
          >
            <ClipboardList className="h-4 w-4" /> Ringkasan profil
          </Link>
          {navItems.map((item) => {
            const Icon = item.icon;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`mb-1 flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm transition-colors hover:bg-white/10 ${
                  pathname.startsWith(item.href) ? "bg-white/10 text-teal-400" : "text-sand-100/80"
                }`}
              >
                <Icon className="h-4 w-4" /> {item.label}
              </Link>
            );
          })}
        </nav>

        <div className="border-t border-white/10 p-3">
          <Button variant="ghost" className="w-full text-sand-100 hover:bg-white/10" onClick={() => logout()}>
            Keluar
          </Button>
        </div>
      </aside>

      <div className="flex min-h-screen flex-1 flex-col">
        <header className="flex items-center justify-between gap-3 border-b border-sand-200 bg-white px-5 py-3 md:hidden">
          <Link href="/products" className="flex items-center gap-1.5 text-sm font-medium text-navy-800">
            <ArrowLeft className="h-4 w-4" />
            Katalog
          </Link>
          <span className="font-display text-lg text-navy-950">SEAPEDIA</span>
          <Button size="sm" variant="ghost" onClick={() => logout()}>
            Keluar
          </Button>
        </header>
        <main className="flex-1 px-5 py-8 md:px-10">{children}</main>
      </div>
    </div>
  );
}
