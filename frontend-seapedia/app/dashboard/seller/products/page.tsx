"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import { formatIDR } from "@/lib/format";
import { RequireRole } from "@/components/features/RequireRole";
import { Card } from "@/components/ui/Card";
import { Input, Textarea } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import type { ProductResponse, UpsertProductRequest } from "@/lib/types";

const emptyForm: UpsertProductRequest = { name: "", description: "", price: 0, stock: 0 };

function SellerProductsContent() {
  const [products, setProducts] = useState<ProductResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState<UpsertProductRequest>(emptyForm);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  async function load() {
    setLoading(true);
    try {
      const data = await api.get<ProductResponse[]>("/seller/products");
      setProducts(data || []);
    } catch {
      setProducts([]);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  function startCreate() {
    setEditingId(null);
    setForm(emptyForm);
    setError(null);
    setShowForm(true);
  }

  function startEdit(p: ProductResponse) {
    setEditingId(p.id);
    setForm({ name: p.name, description: p.description, price: p.price, stock: p.stock });
    setError(null);
    setShowForm(true);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setSaving(true);
    try {
      if (editingId) {
        await api.put(`/seller/products/${editingId}`, form);
      } else {
        await api.post("/seller/products", form);
      }
      setShowForm(false);
      await load();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal menyimpan produk.");
    } finally {
      setSaving(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm("Hapus produk ini?")) return;
    try {
      await api.del(`/seller/products/${id}`);
      await load();
    } catch (err) {
      alert(err instanceof ApiError ? err.message : "Gagal menghapus produk.");
    }
  }

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Produk saya</p>
          <h1 className="mt-1 font-display text-3xl font-semibold text-navy-950">Kelola produk</h1>
        </div>
        <Button onClick={startCreate}>+ Tambah produk</Button>
      </div>

      {showForm && (
        <Card className="mt-6 max-w-xl p-6">
          <h2 className="font-display text-lg font-semibold text-navy-950">
            {editingId ? "Edit produk" : "Produk baru"}
          </h2>
          <form onSubmit={handleSubmit} className="mt-4 flex flex-col gap-4">
            <Input
              label="Nama produk"
              required
              minLength={2}
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
            />
            <Textarea
              label="Deskripsi"
              rows={3}
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            />
            <div className="grid grid-cols-2 gap-4">
              <Input
                label="Harga (Rp)"
                type="number"
                required
                min={1}
                value={form.price || ""}
                onChange={(e) => setForm((f) => ({ ...f, price: Number(e.target.value) }))}
              />
              <Input
                label="Stok"
                type="number"
                required
                min={0}
                value={form.stock}
                onChange={(e) => setForm((f) => ({ ...f, stock: Number(e.target.value) }))}
              />
            </div>
            {error && <p className="text-sm text-red-600">{error}</p>}
            <div className="flex gap-3">
              <Button type="submit" loading={saving}>
                Simpan
              </Button>
              <Button type="button" variant="ghost" onClick={() => setShowForm(false)}>
                Batal
              </Button>
            </div>
          </form>
        </Card>
      )}

      <div className="mt-6 overflow-hidden rounded-2xl border border-sand-200 bg-white">
        {loading ? (
          <p className="p-6 text-sm text-ink/50">Memuat produk...</p>
        ) : products.length === 0 ? (
          <p className="p-6 text-sm text-ink/50">Belum ada produk. Tambahkan produk pertamamu.</p>
        ) : (
          <table className="w-full text-left text-sm">
            <thead className="bg-sand-100/60 text-xs uppercase tracking-wide text-ink/50">
              <tr>
                <th className="px-5 py-3">Produk</th>
                <th className="px-5 py-3">Harga</th>
                <th className="px-5 py-3">Stok</th>
                <th className="px-5 py-3"></th>
              </tr>
            </thead>
            <tbody>
              {products.map((p) => (
                <tr key={p.id} className="border-t border-sand-200">
                  <td className="px-5 py-3 font-medium text-navy-950">{p.name}</td>
                  <td className="px-5 py-3">{formatIDR(p.price)}</td>
                  <td className="px-5 py-3">{p.stock}</td>
                  <td className="px-5 py-3">
                    <div className="flex justify-end gap-2">
                      <Button size="sm" variant="outline" onClick={() => startEdit(p)}>
                        Edit
                      </Button>
                      <Button size="sm" variant="danger" onClick={() => handleDelete(p.id)}>
                        Hapus
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

export default function SellerProductsPage() {
  return (
    <RequireRole role="seller">
      <SellerProductsContent />
    </RequireRole>
  );
}
