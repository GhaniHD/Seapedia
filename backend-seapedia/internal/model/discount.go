package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	DiscountPercent = "percent"
	DiscountFixed   = "fixed"
)

type Voucher struct {
	ID            uuid.UUID
	Code          string
	DiscountType  string
	DiscountValue float64
	ExpiryDate    time.Time
	UsageLimit    int
	UsageCount    int
	CreatedAt     time.Time
}

type Promo struct {
	ID            uuid.UUID
	Code          string
	DiscountType  string
	DiscountValue float64
	ExpiryDate    time.Time
	CreatedAt     time.Time
}
