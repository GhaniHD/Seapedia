package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	DeliveryJobAvailable = "available"
	DeliveryJobTaken     = "taken"
	DeliveryJobCompleted = "completed"
)

type Delivery struct {
	ID            uuid.UUID
	OrderID       uuid.UUID
	DriverID      *uuid.UUID
	Status        string
	Fee           float64
	DriverEarning float64
	TakenAt       *time.Time
	CompletedAt   *time.Time
	CreatedAt     time.Time
}
