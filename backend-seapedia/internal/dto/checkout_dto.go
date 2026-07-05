package dto

import (
	"time"

	"github.com/google/uuid"
)

type CheckoutRequest struct {
	AddressID      string `json:"address_id" binding:"required,uuid"`
	DeliveryMethod string `json:"delivery_method" binding:"required,oneof=instant next_day regular"`
	DiscountCode   string `json:"discount_code"` // optional: kode voucher ATAU promo
}

// CheckoutSummaryResponse selalu menampilkan breakdown lengkap sesuai business rule:
// subtotal, discount, delivery fee, PPN 12%, final total.
type CheckoutSummaryResponse struct {
	OrderID        uuid.UUID `json:"order_id"`
	OrderNo        string    `json:"order_no"`
	Subtotal       float64   `json:"subtotal"`
	DiscountAmount float64   `json:"discount_amount"`
	DiscountKind   string    `json:"discount_kind,omitempty"` // "voucher" | "promo"
	DeliveryFee    float64   `json:"delivery_fee"`
	TaxAmount      float64   `json:"tax_amount"`
	TaxRate        float64   `json:"tax_rate"`
	Total          float64   `json:"total"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
