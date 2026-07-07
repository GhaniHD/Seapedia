# SEAPEDIA — Frontend

Frontend marketplace SEAPEDIA (Level 1–5), dibangun dengan **Next.js 16 (App Router)
+ React 19 + Tailwind CSS 4**, memakai native `fetch` + React Context untuk state
(tanpa library data-fetching tambahan).

## Menjalankan proyek

```bash
npm install
cp .env.example .env.local   # lalu sesuaikan NEXT_PUBLIC_API_URL
npm run dev
```

Buka [http://localhost:3000](http://localhost:3000).

### Environment variable

| Variabel | Default | Keterangan |
| --- | --- | --- |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080/api/v1` | Base URL backend SEAPEDIA (Go/Echo). |

Backend harus berjalan lebih dulu (lihat `backend-seapedia/README.md`) sebelum
fitur yang butuh data (katalog, checkout, dsb) bisa dicoba.

## Struktur folder

```
app/
  (public)/        halaman publik: landing, katalog, detail produk, login, register, select-role
  dashboard/        area privat (butuh login), sidebar berbeda per peran aktif
lib/
  api.ts            fetch wrapper (auto attach JWT, auto unwrap {data: ...} envelope backend)
  types.ts          tipe TS yang mirror DTO backend
  format.ts         helper format Rupiah/tanggal/label
context/
  AuthContext.tsx   state login, profil, daftar peran, peran aktif
components/
  ui/               komponen dasar (Button, Input, Card, StatusPill, dst)
  layout/           Navbar, Footer
  features/         komponen fitur (ReviewSection, ProductCard, RequireRole, AddToCartButton)
```

Struktur `app/dashboard/` per peran (ringkas):
```
dashboard/
  page.tsx              ringkasan profil & peran aktif
  buyer/
    wallet/              saldo + riwayat top-up
    addresses/           alamat pengiriman
    cart/                keranjang (single-store)
    checkout/            checkout (delivery method, voucher/promo, ringkasan pajak)
    orders/[id]/         riwayat & detail pesanan (status history)
    reports/             laporan pengeluaran buyer
  seller/
    store/               profil toko (nama unik)
    products/            CRUD produk
    orders/               pesanan masuk + aksi "Proses pesanan"
    reports/             laporan pendapatan seller
  driver/                 cari job, ambil job, konfirmasi selesai, riwayat & penghasilan
  admin/                  placeholder navigasi (fitur penuh menyusul di Level 6)
```

## Alur bisnis yang diimplementasikan (Level 1–5)

### 1. Autentikasi & peran ganda
- Registrasi selalu menghasilkan akun dengan peran **Buyer** saja.
- Peran lain (Seller/Driver) ditambahkan lewat halaman Dashboard → tombol
  "+ Jadi Penjual/Kurir" (`POST /roles`).
- Jika akun memiliki >1 peran non-admin, user **wajib memilih peran aktif**
  setelah login (`/select-role`) sebelum masuk ke dashboard manapun.
- Semua endpoint privat & tampilan sidebar mengikuti **peran aktif**, bukan
  daftar seluruh peran yang dimiliki. Ini diberlakukan di frontend lewat
  komponen `<RequireRole role="...">` dan tetap divalidasi ulang oleh backend.

### 2. Single-store checkout
Satu keranjang hanya boleh berisi produk dari satu toko. Jika buyer mencoba
menambahkan produk dari toko lain, backend menolak dengan pesan error yang
mengandung indikasi konflik toko; frontend (`AddToCartButton`) menampilkan
tombol **"Kosongkan keranjang & tambah produk ini"** sebagai jalan keluar,
sesuai aturan di soal ("prevent it or clearly ask the buyer to clear the cart
first").

### 3. Perhitungan checkout
Ditampilkan di halaman `/dashboard/buyer/checkout` dan detail pesanan:

```
tax_base = subtotal - discount_amount
tax_amount (PPN 12%) = tax_base * 0.12
total = tax_base + delivery_fee + tax_amount
```

PPN dihitung dari **(subtotal − diskon)**, bukan dari ongkir — mengikuti
implementasi backend di `pkg/utils/order.go` & `checkout_service.go`.

Ongkir tetap per metode (harus sinkron dengan backend):

| Metode | Ongkir |
| --- | --- |
| Instant | Rp25.000 |
| Next Day | Rp12.000 |
| Regular | Rp8.000 |

> Catatan: Nilai ongkir di halaman checkout adalah **estimasi pra-konfirmasi**
> (ditulis ulang di frontend agar UI bisa menampilkan ringkasan sebelum
> submit). Angka final yang benar-benar dipakai (termasuk potongan
> voucher/promo yang tervalidasi) selalu berasal dari respons
> `POST /buyer/checkout`, bukan dari estimasi ini.

### 4. Diskon: Voucher & Promo
Di halaman checkout (`/dashboard/buyer/checkout`), buyer bisa memasukkan kode
voucher/promo (opsional) sebelum konfirmasi. Voucher & promo yang masih aktif
(belum expired, dan untuk voucher masih ada sisa kuota) ditampilkan sebagai
pilihan cepat di form. Validasi kode (expired/habis kuota/tidak ditemukan)
dilakukan di backend saat `POST /buyer/checkout`; hasil validasi (`discount_amount`,
`discount_kind`: `voucher` atau `promo`) ditampilkan jelas di ringkasan checkout,
terpisah dari ongkir dan PPN.

### 5. Pemrosesan pesanan seller & job driver
- **Seller** (`/dashboard/seller/orders`): melihat daftar pesanan masuk beserta
  riwayat status (`status_history` dengan timestamp), dan bisa menekan
  **"Proses pesanan"** untuk memindahkan status dari `Sedang Dikemas` →
  `Menunggu Pengirim`. Hanya pesanan milik toko sendiri yang bisa diproses.
- **Driver** (`/dashboard/driver`): melihat daftar *job tersedia* (pesanan yang
  sudah berstatus `Menunggu Pengirim`), bisa **ambil job** (status →
  `Sedang Dikirim`), lalu **konfirmasi selesai** (status → `Pesanan Selesai`).
  Race condition ditangani di backend — jika dua driver mencoba ambil job yang
  sama, hanya salah satu yang berhasil dan frontend menampilkan pesan error
  lalu me-refresh daftar. Dashboard driver juga menampilkan ringkasan
  penghasilan (total & jumlah job selesai) dan riwayat job yang sudah selesai.

### 6. Ulasan aplikasi publik
Form di landing page (`#reviews`) bisa diisi oleh siapa saja (guest atau user
login) tanpa perlu checkout, sesuai aturan bisnis. Komentar dirender sebagai
teks React biasa (`{comment}`), **tidak pernah** lewat `dangerouslySetInnerHTML`,
sehingga otomatis aman dari eksekusi script (hardening XSS penuh menyusul di
Level 7).

### 7. Amplop respons backend
Semua endpoint backend membungkus payload sukses sebagai `{ "data": ... }`
atau `{ "message": ... }`, dan error sebagai `{ "error": ... }`. Ini di-unwrap
sekali di `lib/api.ts` sehingga seluruh kode halaman bekerja langsung dengan
tipe DTO tanpa perlu mengurus envelope berulang-ulang.

## Cakupan level saat ini

Iterasi ini mengimplementasikan **Level 1 sampai Level 5** penuh:
- Level 1: landing, katalog & detail produk publik, auth + role awareness,
  ulasan publik, komponen & routing dasar.
- Level 2: manajemen toko seller (nama unik), CRUD produk seller, katalog
  publik terhubung ke data asli.
- Level 3: dompet & top-up buyer, alamat pengiriman, keranjang single-store,
  checkout (subtotal/ongkir/PPN/total), riwayat & detail pesanan buyer,
  daftar pesanan masuk seller.
- Level 4: voucher & promo (validasi expiry/kuota di backend, efek diskon
  tampil di ringkasan checkout), aksi seller "Proses pesanan" (Sedang Dikemas
  → Menunggu Pengirim) dengan status history bertimestamp, laporan
  pengeluaran buyer dan laporan pendapatan seller.
- Level 5: dashboard driver penuh — cari job tersedia (hanya pesanan
  `Menunggu Pengirim`), ambil job (→ Sedang Dikirim, dengan penanganan race
  condition di backend), konfirmasi selesai (→ Pesanan Selesai), ringkasan
  penghasilan, dan riwayat job.

Dashboard Admin sudah punya entri navigasi sebagai placeholder untuk Level 6
(Admin Monitoring and Overdue Handling), sesuai instruksi soal ("higher levels
assume previous levels are done; placeholders are acceptable until the
relevant level").
