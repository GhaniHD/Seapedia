"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import { RequireRole } from "@/components/features/RequireRole";
import { Card } from "@/components/ui/Card";
import { Input, Textarea } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import type { StoreResponse, UpsertStoreRequest } from "@/lib/types";

function SellerStoreContent() {
  const [store, setStore] = useState<StoreResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState<UpsertStoreRequest>({ name: "", description: "" });
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [saving, setSaving] = useState(false);

  async function load() {
    setLoading(true);
    try {
      const data = await api.get<StoreResponse>("/seller/store");
      setStore(data);
      setForm({ name: data.name, description: data.description });
    } catch {
      setStore(null);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    setSaving(true);
    try {
      const data = await api.post<StoreResponse>("/seller/store", form);
      setStore(data);
      setSuccess(true);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal menyimpan toko.");
    } finally {
      setSaving(false);
    }
  }

  if (loading) return <p className="text-sm text-ink/50">Memuat data toko...</p>;

  return (
    <div className="max-w-2xl">
      <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Profil toko</p>
      <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">
        {store ? "Kelola toko kamu" : "Buka toko baru"}
      </h1>
      <p className="mt-1 text-sm text-ink/60">
        Nama toko akan tampil ke publik dan harus unik di seluruh SEAPEDIA.
      </p>

      <Card className="mt-6 p-6">
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            label="Nama toko"
            required
            minLength={3}
            maxLength={255}
            value={form.name}
            onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            placeholder="cth. Berkah Jaya Store"
          />
          <Textarea
            label="Deskripsi toko"
            rows={4}
            maxLength={2000}
            value={form.description}
            onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            placeholder="Ceritakan tentang tokomu..."
          />
          {error && <p className="text-sm text-red-600">{error}</p>}
          {success && <p className="text-sm text-teal-600">Toko berhasil disimpan.</p>}
          <Button type="submit" loading={saving} className="w-fit">
            {store ? "Simpan perubahan" : "Buka toko"}
          </Button>
        </form>
      </Card>
    </div>
  );
}

export default function SellerStorePage() {
  return (
    <RequireRole role="seller">
      <SellerStoreContent />
    </RequireRole>
  );
}
