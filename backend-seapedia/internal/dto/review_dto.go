package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateReviewRequest struct {
	ReviewerName string `json:"reviewer_name" binding:"required,min=2,max=255"`
	Rating       int    `json:"rating" binding:"required,min=1,max=5"`
	Comment      string `json:"comment" binding:"required,min=1,max=2000"`
}

type ReviewResponse struct {
	ID           uuid.UUID `json:"id"`
	ReviewerName string    `json:"reviewer_name"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}
