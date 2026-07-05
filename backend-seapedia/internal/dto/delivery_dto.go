package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryJobResponse struct {
	ID          uuid.UUID  `json:"id"`
	OrderID     uuid.UUID  `json:"order_id"`
	OrderNo     string     `json:"order_no"`
	StoreName   string     `json:"store_name"`
	Address     string     `json:"address"`
	Fee         float64    `json:"fee"`
	Status      string     `json:"status"`
	TakenAt     *time.Time `json:"taken_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type DriverEarningResponse struct {
	CompletedJobs int     `json:"completed_jobs"`
	TotalEarning  float64 `json:"total_earning"`
}
