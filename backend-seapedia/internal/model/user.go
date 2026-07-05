package model

import (
	"time"

	"github.com/google/uuid"
)

// User adalah data akun. Role TIDAK disimpan di sini karena satu username
// non-admin bisa punya lebih dari satu role sekaligus (lihat UserRole).
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

// UserRole merepresentasikan role yang dimiliki seorang user.
// Satu user bisa punya banyak baris (contoh: buyer + seller).
type UserRole struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Role   string // admin | seller | buyer | driver
}

const (
	RoleAdmin  = "admin"
	RoleSeller = "seller"
	RoleBuyer  = "buyer"
	RoleDriver = "driver"
	// RoleGuest bukan role yang disimpan di DB. Guest = user tanpa token valid.
	// Dipakai hanya sebagai flag runtime di middleware, tidak pernah ditulis ke database.
	RoleGuest = "guest"
)
