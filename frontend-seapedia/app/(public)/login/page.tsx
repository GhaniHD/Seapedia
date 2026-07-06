"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { ApiError } from "@/lib/api";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card } from "@/components/ui/Card";
import { ShipWheelMark } from "@/components/ui/WaveDivider";

export default function LoginPage() {
  const { login, isAuthenticated, needRoleSelection, loading: authLoading } = useAuth();
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  // Sudah login (atau sedang menunggu pemilihan role) tidak boleh balik lagi
  // ke halaman login. Redirect ke tempat yang seharusnya.
  useEffect(() => {
    if (authLoading) return;
    if (needRoleSelection) {
      router.replace("/select-role");
    } else if (isAuthenticated) {
      router.replace("/dashboard");
    }
  }, [authLoading, isAuthenticated, needRoleSelection, router]);

  if (authLoading || isAuthenticated || needRoleSelection) {
    return (
      <div className="flex min-h-[calc(100vh-64px)] items-center justify-center bg-sand-100/60">
        <p className="text-sm text-ink/50">Memuat...</p>
      </div>
    );
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await login({ email, password });
      if (res.need_role_selection) {
        router.push("/select-role");
      } else {
        router.push("/dashboard");
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal masuk. Coba lagi.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-[calc(100vh-64px)] items-center justify-center bg-sand-100/60 px-5 py-16">
      <Card className="w-full max-w-md p-8">
        <div className="flex items-center gap-2 text-navy-950">
          <ShipWheelMark className="h-7 w-7 text-teal-500" />
          <span className="font-display text-xl">SEAPEDIA</span>
        </div>
        <h1 className="mt-6 font-display text-2xl font-semibold text-navy-950">Selamat datang kembali</h1>
        <p className="mt-1 text-sm text-ink/60">Masuk untuk melanjutkan ke dashboard kamu.</p>

        <form onSubmit={handleSubmit} className="mt-6 flex flex-col gap-4">
          <Input
            label="Email"
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="kamu@email.com"
          />
          <div>
            <Input
              label="Kata sandi"
              type={showPassword ? "text" : "password"}
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
            />
            <label className="mt-2 flex items-center gap-2 text-xs text-ink/60">
              <input
                type="checkbox"
                checked={showPassword}
                onChange={(e) => setShowPassword(e.target.checked)}
                className="h-3.5 w-3.5 rounded border-sand-300 text-teal-500 focus:ring-teal-500"
              />
              Tampilkan sandi
            </label>
          </div>
          {error && <p className="text-sm text-red-600">{error}</p>}
          <Button type="submit" loading={loading} className="mt-2">
            Masuk
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-ink/60">
          Belum punya akun?{" "}
          <Link href="/register" className="font-semibold text-coral-600 hover:underline">
            Daftar sekarang
          </Link>
        </p>
      </Card>
    </div>
  );
}