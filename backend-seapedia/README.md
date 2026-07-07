# SEAPEDIA Backend

Backend API untuk SEAPEDIA — marketplace multi-role (Admin, Seller, Buyer, Driver) untuk COMPFEST 18 Technical Challenge.
Dibangun dengan **Go + Gin + PostgreSQL (pgx)**, arsitektur layer: `handler -> service -> repository -> database`.

Status implementasi: **Level 1 - Level 7 (Core Challenge, 100 pts)** + sebagian bonus (lihat bagian "Yang Sudah & Belum").

---

## 1. Cara Menjalankan

### Opsi A — Docker Compose (disarankan, paling mudah)

```bash
cp .env.example .env
docker compose up --build
```

API akan jalan di `http://localhost:8080`. Migration & seed data demo otomatis dijalankan saat start.

### Opsi B — Manual (Go + PostgreSQL lokal)

```bash
cp .env.example .env
# sesuaikan .env dengan koneksi Postgres lokal Anda

go mod download
go run ./cmd/app
```

Migration (`db/migrations/*.up.sql`) dan seeding otomatis jalan saat aplikasi start (lihat `internal/migration` & `db/seeders.go`). Tidak perlu tool migration terpisah — cukup jalankan aplikasinya.

### Environment Variables (`.env`)

| Key | Keterangan | Default |
|---|---|---|
| `APP_PORT` | Port HTTP server | 8080 |
| `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` | Koneksi PostgreSQL | - |
| `JWT_SECRET` | Secret untuk sign JWT. **Wajib diganti di production.** | - |

---

## 2. Akun Demo (Seed Data)

Dibuat otomatis saat aplikasi start (idempotent, aman dijalankan ulang). Password minimal 8 karakter sesuai validasi register.

| Role | Email | Password |
|---|---|---|
| Admin | `admin@seapedia.com` | `admin12345` |
| Seller (sudah punya toko + 2 produk demo) | `seller@seapedia.com` | `seller12345` |
| Buyer (saldo wallet Rp 1.000.000) | `buyer@seapedia.com` | `buyer12345` |
| Driver | `driver@seapedia.com` | `driver12345` |

Kode diskon demo yang ikut di-seed: voucher **`DEMO10`** (10%, kuota 100) dan promo **`PROMO5K`** (potongan Rp 5.000).

**Membuat admin baru:** tidak ada endpoint publik untuk register sebagai admin (sengaja, demi keamanan — lihat bagian Security Notes). Tambahkan lewat SQL langsung:
```sql
INSERT INTO user_roles (id, user_id, role) VALUES (gen_random_uuid(), '<user_id>', 'admin');
```

---

## 3. Model Role: Guest, Admin, Seller, Buyer, Driver

Sesuai business rule di dokumen tugas, **role yang disimpan di database hanya 4**: `admin`, `seller`, `buyer`, `driver` (tabel `user_roles`, many-to-many terhadap `users` — satu akun non-admin boleh punya lebih dari satu role sekaligus).

**Guest bukan role di database.** Guest adalah kondisi "belum login" — direpresentasikan sebagai request tanpa token valid, ditangani langsung oleh routing (endpoint publik) tanpa butuh flag di database. Ini sengaja dipisah agar konsisten dengan kalimat spek: *"SEAPEDIA has four account roles: Admin, Seller, Buyer, and Driver"*, dan Guest cukup dijelaskan sebagai *"users without an account"*.

### Alur Login & Role Aktif

1. **Register** (`POST /api/v1/register`) → otomatis mendapat role `buyer` (karena buyer adalah role default paling wajar: wajib punya wallet, cart, address).
2. User bisa menambah role lain untuk akunnya sendiri lewat `POST /api/v1/roles` (contoh: `{"role":"seller"}` supaya bisa buka toko, atau `{"role":"driver"}`).
3. **Login** (`POST /api/v1/login`):
   - Kalau user cuma punya **1 role** → langsung dapat token aktif dengan `active_role` terisi.
   - Kalau user punya **>1 role non-admin** → dapat **temp token** (`need_role_selection: true`, `active_role: ""`). Temp token ini **tidak bisa** dipakai untuk endpoint privat manapun sampai memanggil `POST /api/v1/select-role` dengan body `{"role": "seller"}` (atau buyer/driver) untuk mendapat token final.
   - **Admin selalu langsung aktif** tanpa perlu pilih role (dipisahkan dari logika multi-role non-admin, sesuai business rule).
4. Semua endpoint privat memverifikasi `active_role` **dari JWT di server**, bukan dari yang dikirim frontend — lihat `internal/middleware/auth_middleware.go`.

---

## 4. Business Rules & Cara Perhitungan (penting untuk evaluator)

### Single-store checkout
Satu cart (`carts` + `cart_items`) hanya boleh berisi produk dari **satu toko**. Field `carts.store_id` di-set begitu item pertama ditambahkan; kalau buyer coba tambah produk dari toko lain, request **ditolak** dengan pesan jelas minta clear cart dulu (`internal/service/cart_service.go`, fungsi `AddItem`). Cart otomatis "bebas ganti toko" lagi setelah kosong.

### Perhitungan Checkout
```
subtotal        = Σ (harga produk × qty)   -- dari harga TERKINI di database, bukan cache cart
discount_amount = dihitung dari voucher/promo (lihat di bawah)
delivery_fee    = tetap per metode pengiriman (lihat tabel di bawah)
tax_amount      = (subtotal - discount_amount) × 12%      <-- PPN dihitung dari nilai barang setelah diskon, TIDAK termasuk ongkir
total           = (subtotal - discount_amount) + delivery_fee + tax_amount
```
Lihat `pkg/utils/order.go` dan `internal/service/checkout_service.go`.

### Voucher vs Promo
- **Tidak bisa digabung** — satu checkout hanya menerima **satu** `discount_code`, dicari di tabel voucher ATAU promo (mana yang cocok). Hasilnya dibedakan lewat field `discount_kind` (`"voucher"` / `"promo"`) di response checkout.
- Voucher: punya `usage_limit` & `usage_count`, ditolak kalau kuota habis atau kedaluwarsa.
- Promo: hanya punya `expiry_date`, tanpa batas kuota.
- Tipe diskon: `percent` (dibatasi maksimal 100% dari subtotal) atau `fixed` (dibatasi maksimal sebesar subtotal, tidak bisa membuat subtotal negatif).

### Ongkir & SLA Pengiriman (dipakai untuk deadline & overdue)
| Metode | Ongkir | SLA (dari saat diambil driver) |
|---|---|---|
| `instant` | Rp 25.000 | 3 jam |
| `next_day` | Rp 12.000 | 24 jam |
| `regular` | Rp 8.000 | 72 jam |

### Driver Earning
Driver mendapat **80% dari ongkir** per job selesai; 20% sisanya dianggap biaya platform. Lihat `internal/service/order_service.go` (dibuat saat Seller memproses order) — konstanta ada di `pkg/utils/order.go`.

### Overdue Handling (simulasi hari & auto refund/return)
- `POST /api/v1/admin/simulate-next-day` (Admin only, body opsional `{"days": 1}`) memajukan **virtual clock** sistem (tabel `system_clock`) — bukan waktu asli server, supaya evaluator bisa demo tanpa menunggu.
- Setelah virtual clock dimajukan, sistem otomatis memindai order berstatus **"Sedang Dikirim"** yang `deadline_at`-nya sudah lewat, lalu:
  1. Status order → **"Dikembalikan"** (+ tercatat di `order_status_history`).
  2. Saldo buyer di-**refund penuh** ke wallet (tercatat di `wallet_transactions` tipe `refund`).
  3. Kalau income seller sudah tercatat, ditandai `seller_income_reversed = true` agar tidak dihitung dobel di laporan income seller.
  4. Stok produk **dikembalikan** sesuai qty di `order_items`.
- Setiap langkah dijaga flag (`refunded`, `seller_income_reversed`, `stock_restored`) di tabel `orders` supaya **tidak terjadi double refund / double reversal / double restore stock** kalau endpoint simulate dipanggil berkali-kali. Lihat `internal/service/overdue_service.go`.

### Status Order (harus selalu salah satu dari 5 ini di UI)
`Sedang Dikemas` → `Menunggu Pengirim` (setelah seller proses) → `Sedang Dikirim` (setelah driver ambil job) → `Pesanan Selesai` (driver konfirmasi selesai) **atau** `Dikembalikan` (overdue).

---

## 5. Security Notes (Level 7)

- **SQL Injection**: semua query di layer `internal/repository` menggunakan **parameterized query** pgx (`$1, $2, ...`) — tidak ada string concatenation untuk membentuk SQL dari input user.
- **XSS**: input publik yang berpotensi berisi HTML/script (review aplikasi, nama toko, deskripsi produk, alamat) di-**escape** (`html.EscapeString`, lihat `pkg/utils/sanitize.go`) sebelum disimpan ke database — sehingga aman ditampilkan mentah oleh frontend manapun tanpa risiko script tereksekusi.
- **Validasi input**: semua request body memakai `binding` tag Gin (`required`, `email`, `min`, `max`, `gt`, `oneof`, dll) — request tidak valid ditolak dengan pesan error sebelum sampai ke business logic.
- **Password**: di-hash dengan **bcrypt** (`pkg/crypto`), tidak pernah disimpan/dikembalikan dalam bentuk plain text.
- **Autentikasi**: JWT HS256, expired 24 jam untuk token aktif dan 15 menit untuk temp-token (role selection). Logout bersifat client-side (hapus token) karena stateless — didokumentasikan di endpoint `POST /logout`.
- **Otorisasi berbasis role di server**: middleware `RequireRole` mengecek `active_role` **dari klaim JWT** (bukan dari body/header yang bisa dipalsukan user), lihat `internal/middleware/auth_middleware.go`.
- **Ownership check**: setiap aksi privat memverifikasi kepemilikan resource di service layer sebelum eksekusi — contoh: Seller hanya bisa update/delete produk miliknya sendiri (`product_service.go`), Buyer hanya bisa lihat order miliknya sendiri (`order_service.go`), Driver hanya bisa complete job yang sudah diambilnya sendiri (`delivery_service.go`). Mengubah `id` di URL untuk mengakses resource orang lain akan ditolak dengan error, bukan mengembalikan data.
- **Race condition pada resource kompetitif**: pengurangan stok (`WHERE stock >= qty`), pengambilan job driver (`WHERE status='available'`), dan pengurangan saldo wallet (`WHERE balance >= amount`) semuanya di-guard langsung di level query SQL (atomic, `RowsAffected` dicek), sehingga tidak butuh row-level lock manual dan aman dari race condition (misal 2 driver klik "take job" bersamaan).

**Suggested test case (SQL Injection)**: coba login dengan email `' OR '1'='1` — akan gagal karena parameterized query memperlakukannya sebagai string literal, bukan bagian dari SQL.
**Suggested test case (XSS)**: submit review aplikasi dengan comment `<script>alert(1)</script>` — akan tersimpan sebagai teks ter-escape (`&lt;script&gt;...`) dan tidak pernah tereksekusi saat ditampilkan.

---

## 6. Daftar Endpoint (API Reference)

Base URL: `/api/v1`. Endpoint privat butuh header `Authorization: Bearer <token>`.

### Public (Guest)
| Method | Endpoint | Keterangan |
|---|---|---|
| POST | `/register` | Daftar akun baru (otomatis role buyer) |
| POST | `/login` | Login |
| GET | `/products` | List semua produk (katalog publik) |
| GET | `/products/:id` | Detail produk |
| GET | `/stores` | List semua toko |
| GET | `/stores/:id` | Detail toko |
| POST | `/reviews` | Submit review aplikasi (guest/login, tanpa checkout) |
| GET | `/reviews` | List review aplikasi |

### Auth Umum (butuh login)
| Method | Endpoint | Keterangan |
|---|---|---|
| POST | `/select-role` | Pilih role aktif (pakai temp-token) |
| POST | `/logout` | Logout (client discard token) |
| GET | `/profile` | Profil + roles + active_role + ringkasan saldo lintas role |
| POST | `/roles` | Tambah role baru untuk akun sendiri |
| GET | `/vouchers`, `/promos` | Lihat kode diskon yang tersedia |

### Buyer (`active_role=buyer`)
| Method | Endpoint | Keterangan |
|---|---|---|
| POST | `/buyer/wallet/topup` | Dummy top up |
| GET | `/buyer/wallet` | Cek saldo |
| GET | `/buyer/wallet/transactions` | Riwayat transaksi wallet |
| POST/GET | `/buyer/addresses` | Tambah / list alamat |
| GET / POST / PUT / DELETE | `/buyer/cart...` | Kelola cart (single-store rule) |
| POST | `/buyer/checkout` | Checkout |
| GET | `/buyer/orders`, `/buyer/orders/:id` | Riwayat & detail order |
| GET | `/buyer/reports/spending` | Laporan pengeluaran |

### Seller (`active_role=seller`)
| Method | Endpoint | Keterangan |
|---|---|---|
| POST/PUT/GET | `/seller/store` | Buat/update/lihat toko sendiri |
| POST/PUT/DELETE/GET | `/seller/products...` | CRUD produk sendiri |
| GET | `/seller/orders` | Order masuk |
| POST | `/seller/orders/:id/process` | Proses order → buka job untuk driver |
| GET | `/seller/reports/income` | Laporan income |

### Driver (`active_role=driver`)
| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/driver/jobs`, `/driver/jobs/:id` | Cari & lihat job tersedia |
| POST | `/driver/jobs/:id/take` | Ambil job |
| POST | `/driver/jobs/:id/complete` | Konfirmasi selesai |
| GET | `/driver/my-jobs`, `/driver/earnings` | Riwayat & pendapatan |

### Admin (`active_role=admin`)
| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/admin/dashboard` | Monitoring users/stores/products/orders/vouchers/promos/deliveries/overdue |
| POST | `/admin/simulate-next-day` | Simulasi hari berikutnya + trigger overdue handling |
| POST/GET | `/admin/vouchers`, `/admin/promos` | Kelola diskon |

Contoh request lengkap tersedia di folder `tests/rest-client` (format `.http`, bisa langsung dijalankan dengan ekstensi REST Client di VS Code) — silakan lengkapi dengan skenario di atas untuk keperluan demo/testing manual.

---

## 8. Catatan Teknis

- Dependency di `go.mod` disesuaikan ke versi yang kompatibel dengan toolchain Go yang tersedia di lingkungan build (beberapa `replace` directive mengarahkan modul `golang.org/x/*` ke mirror GitHub-nya) — ini murni penyesuaian versi, tidak mengubah perilaku library.
- Migration dijalankan lewat runner custom ringan (`internal/migration`, baca file `.sql` di `db/migrations` secara berurutan dan mencatat yang sudah dijalankan di tabel `schema_migrations`) — tanpa dependency migration tool eksternal, supaya "works on any machine" tanpa perlu tool tambahan di luar Go itu sendiri.
