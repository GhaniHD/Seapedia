package utils

import (
	"fmt"
	"time"
)

// GenerateOrderNo membuat nomor order yang mudah dibaca, contoh: SPD-20260705-AB12CD
func GenerateOrderNo(now time.Time, suffix string) string {
	return fmt.Sprintf("SPD-%s-%s", now.Format("20060102"), suffix)
}

// DeliveryFee mengembalikan ongkir berdasarkan metode pengiriman.
// Didokumentasikan juga di README agar konsisten FE/BE.
func DeliveryFee(method string) float64 {
	switch method {
	case "instant":
		return 25000
	case "next_day":
		return 12000
	case "regular":
		return 8000
	default:
		return 10000
	}
}

// DeliverySLA mengembalikan batas waktu pengiriman untuk tiap metode,
// dipakai untuk menghitung deadline_at & deteksi overdue (Level 6).
func DeliverySLA(method string) time.Duration {
	switch method {
	case "instant":
		return 3 * time.Hour
	case "next_day":
		return 24 * time.Hour
	case "regular":
		return 72 * time.Hour
	default:
		return 72 * time.Hour
	}
}

// DriverEarningRate: porsi ongkir yang menjadi pendapatan driver. Sisanya (20%) menjadi
// biaya platform. Didokumentasikan di README.
const DriverEarningRate = 0.8

const TaxRate = 0.12 // PPN 12%, dihitung dari (subtotal - discount)
