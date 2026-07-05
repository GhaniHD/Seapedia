package dto

import (
	"time"

	"github.com/google/uuid"
)

type TopupRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WalletResponse struct {
	Balance float64 `json:"balance"`
}

type WalletTransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpsertAddressRequest struct {
	Label     string `json:"label" binding:"required,min=2,max=100"`
	Detail    string `json:"detail" binding:"required,min=5,max=1000"`
	IsDefault bool   `json:"is_default"`
}

type AddressResponse struct {
	ID        uuid.UUID `json:"id"`
	Label     string    `json:"label"`
	Detail    string    `json:"detail"`
	IsDefault bool      `json:"is_default"`
}
