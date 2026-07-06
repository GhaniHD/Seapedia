"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import { RequireRole } from "@/components/features/RequireRole";
import { Card, Badge } from "@/components/ui/Card";
import { Input, Textarea } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import type { AddressResponse, UpsertAddressRequest } from "@/lib/types";

const emptyForm: UpsertAddressRequest = { label: "", detail: "", is_default: false };

function BuyerAddressesContent() {
  const [addresses, setAddresses] = useState<AddressResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<UpsertAddressRequest>(emptyForm);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  async function load() {
    setLoading(true);
    try {
      const data = await api.get<AddressResponse[]>("/buyer/addresses");
      setAddresses(data || []);
    } catch {
      setAddresses([]);
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
    setSaving(true);
    try {
      await api.post("/buyer/addresses", form);
      setForm(emptyForm);
      setShowForm(false);
      await load();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal menyimpan alamat.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="max-w-2xl">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Alamat pengiriman</p>
          <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Alamat saya</h1>
        </div>
        <Button onClick={() => setShowForm((v) => !v)}>{showForm ? "Tutup" : "+ Tambah alamat"}</Button>
      </div>

      {showForm && (
        <Card className="mt-6 p-6">
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <Input
              label="Label alamat"
              required
              minLength={2}
              placeholder="cth. Rumah, Kantor"
              value={form.label}
              onChange={(e) => setForm((f) => ({ ...f, label: e.target.value }))}
            />
            <Textarea
              label="Detail alamat lengkap"
              required
              minLength={5}
              rows={3}
              value={form.detail}
              onChange={(e) => setForm((f) => ({ ...f, detail: e.target.value }))}
            />
            <label className="flex items-center gap-2 text-sm text-ink/70">
              <input
                type="checkbox"
                checked={form.is_default}
                onChange={(e) => setForm((f) => ({ ...f, is_default: e.target.checked }))}
              />
              Jadikan alamat utama
            </label>
            {error && <p className="text-sm text-red-600">{error}</p>}
            <Button type="submit" loading={saving} className="w-fit">
              Simpan alamat
            </Button>
          </form>
        </Card>
      )}

      <div className="mt-6 flex flex-col gap-3">
        {loading && <p className="text-sm text-ink/50">Memuat alamat...</p>}
        {!loading && addresses.length === 0 && (
          <p className="text-sm text-ink/50">Belum ada alamat tersimpan.</p>
        )}
        {addresses.map((a) => (
          <Card key={a.id} className="flex items-start justify-between p-5">
            <div>
              <div className="flex items-center gap-2">
                <p className="font-semibold text-navy-950">{a.label}</p>
                {a.is_default && <Badge tone="teal">Utama</Badge>}
              </div>
              <p className="mt-1 text-sm text-ink/60">{a.detail}</p>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}

export default function BuyerAddressesPage() {
  return (
    <RequireRole role="buyer">
      <BuyerAddressesContent />
    </RequireRole>
  );
}
