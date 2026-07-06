"use client";

import { useEffect, useState } from "react";
import { Star } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { formatDate } from "@/lib/format";
import type { CreateReviewRequest, ReviewResponse } from "@/lib/types";
import { Button } from "@/components/ui/Button";
import { Input, Textarea } from "@/components/ui/Input";
import { Card } from "@/components/ui/Card";

function Stars({ rating }: { rating: number }) {
  return (
    <div className="flex gap-0.5 text-teal-500" aria-label={`Rating ${rating} dari 5`}>
      {Array.from({ length: 5 }).map((_, i) => (
        <Star
          key={i}
          className={`h-4 w-4 ${i < rating ? "fill-teal-500 text-teal-500" : "text-sand-300"}`}
        />
      ))}
    </div>
  );
}

export function ReviewSection() {
  const [reviews, setReviews] = useState<ReviewResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const [form, setForm] = useState<CreateReviewRequest>({
    reviewer_name: "",
    rating: 5,
    comment: "",
  });

  async function loadReviews() {
    setLoading(true);
    try {
      const data = await api.get<ReviewResponse[]>("/reviews", false);
      setReviews(data || []);
    } catch {
      setReviews([]);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadReviews();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    setSubmitting(true);
    try {
      await api.post<ReviewResponse>("/reviews", form, false);
      setForm({ reviewer_name: "", rating: 5, comment: "" });
      setSuccess(true);
      await loadReviews();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Gagal mengirim ulasan.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section id="reviews" className="mx-auto max-w-6xl px-5 py-16">
      <div className="mb-8 max-w-xl">
        <p className="text-sm font-semibold uppercase tracking-wide text-teal-500">Suara pengguna</p>
        <h2 className="mt-1 font-display text-3xl font-semibold text-navy-950">
          Apa kata mereka tentang SEAPEDIA?
        </h2>
        <p className="mt-2 text-ink/60">
          Siapa saja boleh berbagi pengalaman memakai aplikasi ini — tanpa perlu checkout.
        </p>
      </div>

      <div className="grid gap-8 lg:grid-cols-[1fr_1.3fr]">
        <Card className="h-fit p-6">
          <h3 className="font-display text-lg font-semibold text-navy-950">Tulis ulasan aplikasi</h3>
          <form onSubmit={handleSubmit} className="mt-4 flex flex-col gap-4">
            <Input
              label="Nama kamu"
              required
              value={form.reviewer_name}
              onChange={(e) => setForm((f) => ({ ...f, reviewer_name: e.target.value }))}
              placeholder="cth. Ghani D."
              maxLength={255}
            />
            <div>
              <label className="text-sm font-medium text-navy-900">Rating</label>
              <div className="mt-1.5 flex gap-1">
                {[1, 2, 3, 4, 5].map((n) => (
                  <button
                    type="button"
                    key={n}
                    onClick={() => setForm((f) => ({ ...f, rating: n }))}
                    className="transition-transform hover:scale-110"
                    aria-label={`${n} bintang`}
                  >
                    <Star
                      className={`h-6 w-6 ${n <= form.rating ? "fill-teal-500 text-teal-500" : "text-sand-300"}`}
                    />
                  </button>
                ))}
              </div>
            </div>
            <Textarea
              label="Komentar"
              required
              rows={3}
              maxLength={2000}
              value={form.comment}
              onChange={(e) => setForm((f) => ({ ...f, comment: e.target.value }))}
              placeholder="Ceritakan pengalamanmu memakai SEAPEDIA..."
            />
            {error && <p className="text-sm text-red-600">{error}</p>}
            {success && <p className="text-sm text-teal-600">Terima kasih! Ulasanmu sudah tampil.</p>}
            <Button type="submit" loading={submitting}>
              Kirim ulasan
            </Button>
          </form>
        </Card>

        <div className="flex max-h-[520px] flex-col gap-4 overflow-y-auto pr-1">
          {loading && <p className="text-sm text-ink/50">Memuat ulasan...</p>}
          {!loading && reviews.length === 0 && (
            <p className="text-sm text-ink/50">Belum ada ulasan. Jadilah yang pertama!</p>
          )}
          {reviews.map((r) => (
            <Card key={r.id} className="p-5">
              <div className="flex items-center justify-between">
                <p className="font-semibold text-navy-950">{r.reviewer_name}</p>
                <Stars rating={r.rating} />
              </div>
              {/* Rendered as plain text via React — safe from script execution */}
              <p className="mt-2 whitespace-pre-wrap break-words text-sm text-ink/70">{r.comment}</p>
              <p className="mt-2 text-xs text-ink/40">{formatDate(r.created_at)}</p>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
