package dto

import "time"

type CreateVoucherRequest struct {
	Code          string  `json:"code" binding:"required,min=3,max=50"`
	DiscountType  string  `json:"discount_type" binding:"required,oneof=percent fixed"`
	DiscountValue float64 `json:"discount_value" binding:"required,gt=0"`
	ExpiryDate    string  `json:"expiry_date" binding:"required"` // RFC3339
	UsageLimit    int     `json:"usage_limit" binding:"required,gt=0"`
}

type CreatePromoRequest struct {
	Code          string  `json:"code" binding:"required,min=3,max=50"`
	DiscountType  string  `json:"discount_type" binding:"required,oneof=percent fixed"`
	DiscountValue float64 `json:"discount_value" binding:"required,gt=0"`
	ExpiryDate    string  `json:"expiry_date" binding:"required"`
}

type VoucherResponse struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	ExpiryDate    time.Time `json:"expiry_date"`
	UsageLimit    int       `json:"usage_limit"`
	UsageCount    int       `json:"usage_count"`
}

type PromoResponse struct {
	ID            string    `json:"id"`
	Code          string    `json:"code"`
	DiscountType  string    `json:"discount_type"`
	DiscountValue float64   `json:"discount_value"`
	ExpiryDate    time.Time `json:"expiry_date"`
}
