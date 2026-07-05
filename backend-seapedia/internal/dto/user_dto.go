package dto

import "github.com/google/uuid"

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type ProfileResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Roles      []string  `json:"roles"`
	ActiveRole string    `json:"active_role"`
	// Ringkasan saldo lintas role (Level 1: placeholder, diisi penuh di level lanjutan)
	WalletBalance *float64 `json:"wallet_balance,omitempty"`
	StoreIncome   *float64 `json:"store_income,omitempty"`
	DriverEarning *float64 `json:"driver_earning,omitempty"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse: kalau user punya >1 role non-admin, Token yang dikembalikan
// adalah "temp token" (active_role kosong) dan NeedRoleSelection=true.
// Frontend wajib memanggil POST /api/v1/select-role sebelum bisa akses dashboard privat.
type LoginResponse struct {
	Token             string       `json:"token"`
	NeedRoleSelection bool         `json:"need_role_selection"`
	Roles             []string     `json:"roles"`
	ActiveRole        string       `json:"active_role"`
	User              UserResponse `json:"user"`
}

type SelectRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type AddRoleRequest struct {
	Role string `json:"role" binding:"required"` // seller | buyer | driver
}
