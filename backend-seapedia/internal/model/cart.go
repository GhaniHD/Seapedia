package model

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	StoreID *uuid.UUID
}

type CartItem struct {
	ID        uuid.UUID
	CartID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	CreatedAt time.Time
}
