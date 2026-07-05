package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Role     string
	Password string
	CreateAt time.Time
	UpdateAt time.Time
}
