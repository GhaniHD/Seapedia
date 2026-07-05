package model

import (
	"time"

	"github.com/google/uuid"
)

// AppReview adalah review publik tentang aplikasi/website SEAPEDIA
// (bukan review produk atau transaksi). Bisa diisi guest maupun user login.
type AppReview struct {
	ID           uuid.UUID
	ReviewerName string
	Rating       int
	Comment      string
	CreatedAt    time.Time
}
