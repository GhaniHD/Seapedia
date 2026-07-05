package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderItemResponse struct {
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type StatusHistoryResponse struct {
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderResponse struct {
	ID             uuid.UUID                `json:"id"`
	OrderNo        string                   `json:"order_no"`
	StoreName      string                   `json:"store_name,omitempty"`
	BuyerName      string                   `json:"buyer_name,omitempty"`
	DeliveryMethod string                   `json:"delivery_method"`
	Subtotal       float64                  `json:"subtotal"`
	DiscountAmount float64                  `json:"discount_amount"`
	DeliveryFee    float64                  `json:"delivery_fee"`
	TaxAmount      float64                  `json:"tax_amount"`
	Total          float64                  `json:"total"`
	Status         string                   `json:"status"`
	DeadlineAt     *time.Time               `json:"deadline_at,omitempty"`
	Items          []OrderItemResponse      `json:"items,omitempty"`
	StatusHistory  []StatusHistoryResponse  `json:"status_history,omitempty"`
	CreatedAt      time.Time                `json:"created_at"`
}

type SpendingReportResponse struct {
	TotalOrders   int     `json:"total_orders"`
	TotalSpending float64 `json:"total_spending"`
}

type IncomeReportResponse struct {
	TotalOrders    int     `json:"total_orders"`
	TotalIncome    float64 `json:"total_income"`
	TotalReversed  float64 `json:"total_reversed"`
}
