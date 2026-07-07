# SEAPEDIA

Marketplace multi-role (Admin, Seller, Buyer, Driver) untuk **COMPFEST 18
Technical Challenge**. Satu ekosistem, dua bagian:

```
SEAPEDIA/
├── backend-seapedia/    API (Go + Gin + PostgreSQL)
├── frontend-seapedia/   Web app (Next.js 16 + React 19 + Tailwind CSS 4)
└── README.md            (file ini)
```

- **Backend**: lihat [`backend-seapedia/README.md`](backend-seapedia/README.md) —
  cara jalan, akun demo, business rules & cara perhitungan, security notes,
  daftar endpoint lengkap.
- **Frontend**: lihat [`frontend-seapedia/README.md`](frontend-seapedia/README.md) —
  struktur folder, environment variable, alur bisnis per halaman.

Dokumen ini cuma ringkasan supaya evaluator bisa langsung jalanin dan demo
kedua bagian sekaligus tanpa harus baca dua README secara terpisah dulu.

---

## Status Implementasi

| Bagian | Level yang diklaim | Keterangan |
|---|---|---|
| Backend | **Level 1 – 7** (100 pts core) | Lihat detail di README backend |
| Frontend | **Level 1 – 5** | Admin dashboard masih placeholder navigasi (Level 6), belum ada halaman khusus untuk fitur admin monitoring/overdue trigger/voucher-promo management UI (Level 6) dan belum ada halaman finalisasi keamanan tambahan di sisi UI (Level 7) |

> **Catatan:** backend sudah mengimplementasikan seluruh business logic sampai
> Level 7 (termasuk endpoint admin dashboard, generate voucher/promo,
> simulate-next-day, dan hardening keamanan), sementara frontend saat ini
> fokus menyelesaikan UI sampai Level 5 — dashboard Admin masih berupa
> placeholder navigasi dan akan dilengkapi pada iterasi berikutnya.
> Endpoint-endpoint tersebut tetap bisa didemokan lewat
> `tests/rest-client/*.http` atau Postman/Swagger.

---

## Menjalankan Proyek (Full Stack)

### 1. Jalankan backend lebih dulu

```bash
cd backend-seapedia
cp .env.example .env
docker compose up --build
```

API akan jalan di `http://localhost:8080` (migration & seed data demo jalan
otomatis). Lihat [`backend-seapedia/README.md`](backend-seapedia/README.md)
untuk opsi menjalankan tanpa Docker.

### 2. Jalankan frontend

```bash
cd frontend-seapedia
npm install
cp .env.example .env.local   # pastikan NEXT_PUBLIC_API_URL mengarah ke backend di atas
npm run dev
```

Buka [http://localhost:3000](http://localhost:3000).

---

## Akun Demo

Dibuat otomatis lewat seeder backend saat pertama kali start:

| Role | Email | Password |
|---|---|---|
| Admin | `admin@seapedia.com` | `admin12345` |
| Seller (sudah punya toko + 2 produk demo) | `seller@seapedia.com` | `seller12345` |
| Buyer (saldo wallet Rp 1.000.000) | `buyer@seapedia.com` | `buyer12345` |
| Driver | `driver@seapedia.com` | `driver12345` |

Kode diskon demo: voucher **`DEMO10`** (10%, kuota 100) dan promo **`PROMO5K`**
(potongan Rp 5.000).

---

## Alur Demo Singkat (End-to-End)

1. **Guest**: buka landing page, lihat katalog & detail produk, submit review
   aplikasi tanpa login.
2. **Register/Login** sebagai Buyer → top up wallet → tambah alamat → tambah
   produk ke cart → checkout (pilih metode pengiriman, opsional kode
   `DEMO10`/`PROMO5K`) → cek riwayat pesanan.
3. **Login** sebagai Seller → lihat pesanan masuk → klik "Proses pesanan"
   (status `Sedang Dikemas` → `Menunggu Pengirim`).
4. **Login** sebagai Driver → ambil job yang baru diproses seller → konfirmasi
   selesai → cek ringkasan penghasilan.
5. **Login** sebagai Admin → (UI masih placeholder) gunakan
   `backend-seapedia/tests/rest-client/seller_driver_admin_flow.http` untuk
   demo monitoring dashboard, generate voucher/promo, dan
   `POST /admin/simulate-next-day` untuk memicu auto refund/return pada
   pesanan overdue.

---

## Dokumentasi API

Format `.http` (REST Client, VS Code) tersedia di
`backend-seapedia/tests/rest-client/`. Detail lengkap tiap endpoint ada di
tabel "Daftar Endpoint" pada [`backend-seapedia/README.md`](backend-seapedia/README.md).