package dto

import "github.com/google/uuid"

type AddCartItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type CartItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Subtotal  float64   `json:"subtotal"`
}

type CartResponse struct {
	StoreID   *uuid.UUID          `json:"store_id"`
	StoreName string              `json:"store_name,omitempty"`
	Items     []CartItemResponse  `json:"items"`
	Subtotal  float64             `json:"subtotal"`
}
