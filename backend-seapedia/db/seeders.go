package db

import (
	"context"
	"log"
	"time"

	"backend-seapedia/internal/model"
	"backend-seapedia/pkg/crypto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SeedDemoData membuat akun demo (Admin, Seller, Buyer, Driver) + voucher/promo contoh,
// supaya evaluator tidak perlu menebak cara membuat akun (Level 7 - Prepare Final Documentation and Demo Data).
// Aman dipanggil berkali-kali (idempotent, pakai ON CONFLICT DO NOTHING).
func SeedDemoData(pool *pgxpool.Pool) {
	ctx := context.Background()

	seedUser := func(name, email, password, role string) uuid.UUID {
		var id uuid.UUID
		err := pool.QueryRow(ctx, `SELECT id FROM users WHERE email=$1`, email).Scan(&id)
		if err != nil {
			hashed, _ := crypto.HashPassword(password)
			id = uuid.New()
			_, err = pool.Exec(ctx, `INSERT INTO users (id, name, email, password) VALUES ($1,$2,$3,$4)`, id, name, email, hashed)
			if err != nil {
				log.Println("seed user gagal:", err)
				return id
			}
		}
		_, _ = pool.Exec(ctx, `INSERT INTO user_roles (id, user_id, role) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`, uuid.New(), id, role)
		_, _ = pool.Exec(ctx, `INSERT INTO wallets (id, user_id, balance) VALUES ($1,$2,0) ON CONFLICT (user_id) DO NOTHING`, uuid.New(), id)
		_, _ = pool.Exec(ctx, `INSERT INTO carts (id, user_id) VALUES ($1,$2) ON CONFLICT (user_id) DO NOTHING`, uuid.New(), id)
		return id
	}

	adminID := seedUser("Admin SEAPEDIA", "admin@seapedia.com", "admin12345", model.RoleAdmin)
	sellerID := seedUser("Toko Demo", "seller@seapedia.com", "seller12345", model.RoleSeller)
	buyerID := seedUser("Buyer Demo", "buyer@seapedia.com", "buyer12345", model.RoleBuyer)
	driverID := seedUser("Driver Demo", "driver@seapedia.com", "driver12345", model.RoleDriver)
	_ = adminID
	_ = driverID

	// beri buyer demo saldo awal supaya bisa langsung checkout
	_, _ = pool.Exec(ctx, `UPDATE wallets SET balance = 1000000 WHERE user_id = $1`, buyerID)

	// store + produk demo untuk seller
	var storeID uuid.UUID
	err := pool.QueryRow(ctx, `SELECT id FROM stores WHERE user_id=$1`, sellerID).Scan(&storeID)
	if err != nil {
		storeID = uuid.New()
		_, err = pool.Exec(ctx, `INSERT INTO stores (id, user_id, name, description) VALUES ($1,$2,$3,$4)`,
			storeID, sellerID, "Toko Demo SEAPEDIA", "Toko contoh untuk testing & demo evaluator")
		if err != nil {
			log.Println("seed store gagal:", err)
		} else {
			_, _ = pool.Exec(ctx, `INSERT INTO products (id, store_id, name, description, price, stock) VALUES
				($1,$2,'Kaos Polos Demo','Kaos katun combed 24s',75000,50),
				($3,$2,'Topi Demo','Topi trucker demo',45000,30)`,
				uuid.New(), storeID, uuid.New())
		}
	}

	// voucher & promo demo
	var vCount int
	_ = pool.QueryRow(ctx, `SELECT COUNT(*) FROM vouchers`).Scan(&vCount)
	if vCount == 0 {
		_, _ = pool.Exec(ctx, `INSERT INTO vouchers (id, code, discount_type, discount_value, expiry_date, usage_limit) VALUES
			($1,'DEMO10','percent',10,$2,100)`, uuid.New(), time.Now().AddDate(1, 0, 0))
	}
	var pCount int
	_ = pool.QueryRow(ctx, `SELECT COUNT(*) FROM promos`).Scan(&pCount)
	if pCount == 0 {
		_, _ = pool.Exec(ctx, `INSERT INTO promos (id, code, discount_type, discount_value, expiry_date) VALUES
			($1,'PROMO5K','fixed',5000,$2)`, uuid.New(), time.Now().AddDate(1, 0, 0))
	}

	log.Println("Seed demo data selesai. Lihat README untuk daftar akun demo.")
}
