"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/context/AuthContext";
import { ApiError } from "@/lib/api";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Card } from "@/components/ui/Card";
import { ShipWheelMark } from "@/components/ui/WaveDivider";

export default function RegisterPage() {
  const { register, login } = useAuth();
  const router = useRouter();
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await register(form);
      // auto-login after register; new accounts always start as buyer only
      const res = await login({ email: form.email, password: form.password });
      router.push(res.need_role_selection ? "/select-role" : "/dashboard");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal mendaftar. Coba lagi.");
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
        <h1 className="mt-6 font-display text-2xl font-semibold text-navy-950">Buat akun baru</h1>
        <p className="mt-1 text-sm text-ink/60">
          Akun baru otomatis memiliki peran Pembeli. Kamu bisa menambah peran Penjual atau
          Kurir kapan saja lewat dashboard.
        </p>

        <form onSubmit={handleSubmit} className="mt-6 flex flex-col gap-4">
          <Input
            label="Nama lengkap"
            required
            minLength={2}
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            placeholder="cth. Ghani Husna"
          />
          <Input
            label="Email"
            type="email"
            required
            value={form.email}
            onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
            placeholder="kamu@email.com"
          />
          <Input
            label="Kata sandi"
            type="password"
            required
            minLength={8}
            hint="Minimal 8 karakter."
            value={form.password}
            onChange={(e) => setForm((f) => ({ ...f, password: e.target.value }))}
            placeholder="••••••••"
          />
          {error && <p className="text-sm text-red-600">{error}</p>}
          <Button type="submit" loading={loading} className="mt-2">
            Daftar
          </Button>
        </form>

        <p className="mt-6 text-center text-sm text-ink/60">
          Sudah punya akun?{" "}
          <Link href="/login" className="font-semibold text-coral-600 hover:underline">
            Masuk di sini
          </Link>
        </p>
      </Card>
    </div>
  );
}
