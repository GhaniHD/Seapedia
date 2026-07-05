package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Balance   float64
	UpdatedAt time.Time
}

type WalletTransaction struct {
	ID          uuid.UUID
	WalletID    uuid.UUID
	Type        string // topup | checkout | refund
	Amount      float64
	Description string
	OrderID     *uuid.UUID
	CreatedAt   time.Time
}

type Address struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Label     string
	Detail    string
	IsDefault bool
	CreatedAt time.Time
}
